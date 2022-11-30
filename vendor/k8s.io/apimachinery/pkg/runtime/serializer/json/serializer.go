package json

import (
	"encoding/json"
	"io"
	"strconv"

	"sigs.k8s.io/yaml"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/recognizer"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/klog"
)

// identifier computes Identifier of Encoder based on the given options.
func identifier(options SerializerOptions) runtime.Identifier {
	result := map[string]string{
		"name":   "json",
		"yaml":   strconv.FormatBool(options.Yaml),
		"pretty": strconv.FormatBool(options.Pretty),
	}
	identifier, err := json.Marshal(result)
	if err != nil {
		klog.Fatalf("Failed marshaling identifier for json Serializer: %v", err)
	}
	return runtime.Identifier(identifier)
}

// SerializerOptions holds the options which are used to configure a JSON/YAML serializer.
// example:
// (1) To configure a JSON serializer, set `Yaml` to `false`.
// (2) To configure a YAML serializer, set `Yaml` to `true`.
// (3) To configure a strict serializer that can return strictDecodingError, set `Strict` to `true`.
type SerializerOptions struct {
	// Yaml: configures the Serializer to work with JSON(false) or YAML(true).
	// When `Yaml` is enabled, 
	// this serializer only supports the subset of YAML that matches JSON, 
	// and will error if constructs are used that do not serialize to JSON.
	Yaml bool

	// Pretty: configures a JSON enabled Serializer(`Yaml: false`) to 
	// produce human-readable output.
	// This option is silently ignored when `Yaml` is `true`.
	Pretty bool

	// Strict 一般为 false
	//
	// Strict: configures the Serializer to return strictDecodingError's 
	// when duplicate fields are present decoding JSON or YAML.
	// Note that enabling this option is not as performant as the non-strict variant, 
	// and should not be used in fast paths.
	Strict bool
}

type Serializer struct {
	// 一般为 pkg/runtime/serializer/json/meta.go -> DefaultMetaFactory 是一个固定值.
	meta    MetaFactory
	options SerializerOptions
	// creater, typer 都是在 Decode() 成员方法中被真正使用(其实也只是当作参数传入).
	// creater 为 pkg/runtime/scheme.go -> Scheme{} 对象(几乎所有的调用都使用这个)
	creater runtime.ObjectCreater
	// typer 与 creater 相同
	typer   runtime.ObjectTyper

	identifier runtime.Identifier
}

// Serializer implements Serializer
var _ runtime.Serializer = &Serializer{}
var _ recognizer.RecognizingDecoder = &Serializer{}

// NewSerializer 简单地构造一个 Serializer{} 对象.
//
// 	@param meta: 一般为 pkg/runtime/serializer/json/meta.go -> DefaultMetaFactory
//
// NewSerializer creates a JSON serializer that handles encoding versioned objects
// into the proper JSON form.
// If typer is not nil, the object has the group, version, and kind fields set.
//
// 	Deprecated: use NewSerializerWithOptions instead.
func NewSerializer(
	meta MetaFactory, 
	creater runtime.ObjectCreater, 
	typer runtime.ObjectTyper,
	pretty bool,
) *Serializer {
	return NewSerializerWithOptions(
		meta, creater, typer, SerializerOptions{false, pretty, false},
	)
}

// NewYAMLSerializer creates a YAML serializer that handles encoding versioned objects
// into the proper YAML form.
// If typer is not nil, the object has the group, version, and kind fields set.
// This serializer supports only the subset of YAML that matches JSON,
// and will error if constructs are used that do not serialize to JSON.
// 	Deprecated: use NewSerializerWithOptions instead.
func NewYAMLSerializer(
	meta MetaFactory, creater runtime.ObjectCreater, typer runtime.ObjectTyper,
) *Serializer {
	return NewSerializerWithOptions(meta, creater, typer, SerializerOptions{true, false, false})
}

// NewSerializerWithOptions 简单地构造一个 Serializer{} 对象.
//
// 	@param meta: 一般为 pkg/runtime/serializer/json/meta.go -> DefaultMetaFactory
// 	@param creater: pkg/runtime/scheme.go -> NewScheme()的返回值(几乎所有的调用都使用这个)
// 	@param typer: 同 creater 参数.
//
// caller: pkg/runtime/serializer/codec_factory.go -> newSerializersForScheme()
//
// NewSerializerWithOptions creates a JSON/YAML serializer
// that handles encoding versioned objects into the proper JSON/YAML form.
// If typer is not nil, the object has the group, version, and kind fields set.
// Options are copied into the Serializer and are immutable.
func NewSerializerWithOptions(
	meta MetaFactory, 
	creater runtime.ObjectCreater, 
	typer runtime.ObjectTyper,
	options SerializerOptions,
) *Serializer {
	return &Serializer{
		meta:       meta,
		creater:    creater,
		typer:      typer,
		options:    options,
		identifier: identifier(options),
	}
}

// 	@param gvk: 可能为 nil, 需要根据 originalData 的 apiVersion 与 kind 字段自行确认资源类型.
//
// caller: pkg/runtime/serializer/recognizer/recognizer.go -> decoder.Decode()
//			其实最终是 client-go/rest/request.go -> Request.Into() 方法, 
//			将 apiserver 的响应数据转换成指定类型的资源对象
//
// Decode attempts to convert the provided data into YAML or JSON,
// extract the stored schema kind, apply the provided default gvk, and then
// load that data into an object matching the desired schema kind 
// or the provided into.
//
// If into is *runtime.Unknown,
// the raw data will be extracted and no decoding will be performed.
//
// If into is not registered with the typer,
// then the object will be straight decoded using normal JSON/YAML unmarshalling.
//
// If into is provided and the original data is not fully qualified with kind/version/group,
// the type of the into will be used to alter the returned gvk.
//
// If into is nil or data's gvk different from into's gvk,
// it will generate a new Object with ObjectCreater.New(gvk)
//
// On success or most errors, the method will return the calculated schema kind.
// The gvk calculate priority will be originalData > default gvk > into
func (s *Serializer) Decode(
	originalData []byte, gvk *schema.GroupVersionKind, into runtime.Object,
) (runtime.Object, *schema.GroupVersionKind, error) {
	data := originalData
	if s.options.Yaml {
		altered, err := yaml.YAMLToJSON(data)
		if err != nil {
			return nil, nil, err
		}
		data = altered
	}

	// actual 是一个 GVK 对象, 基本确认了 data 中的内容是哪种资源.
	actual, err := s.meta.Interpret(data)
	if err != nil {
		return nil, nil, err
	}

	// client-go 的 rest-client Into()方法中传入的 gvk 参数为 nil.
	if gvk != nil {
		*actual = gvkWithDefaults(*actual, *gvk)
	}

	// 如果目标资源类型不存在.
	if unk, ok := into.(*runtime.Unknown); ok && unk != nil {
		unk.Raw = originalData
		unk.ContentType = runtime.ContentTypeJSON
		unk.GetObjectKind().SetGroupVersionKind(*actual)
		return unk, actual, nil
	}

	if into != nil {
		_, isUnstructured := into.(runtime.Unstructured)
		// types: into 对象可能的资源类型列表, 如[corev1.Pod, appsv1.Deployment...]
		// 一般来说, 只会有一个成员吧? default 块里直接就使用了 types[0] 成员.
		types, _, err := s.typer.ObjectKinds(into)
		switch {
		case runtime.IsNotRegisteredError(err), isUnstructured:
			if err := caseSensitiveJsonIterator.Unmarshal(data, into); err != nil {
				return nil, actual, err
			}
			return into, actual, nil
		case err != nil:
			return nil, actual, err
		default:
			// 这里
			*actual = gvkWithDefaults(*actual, types[0])
		}
	}

	// 判断 actual 对象的合法性
	if len(actual.Kind) == 0 {
		return nil, actual, runtime.NewMissingKindErr(string(originalData))
	}
	if len(actual.Version) == 0 {
		return nil, actual, runtime.NewMissingVersionErr(string(originalData))
	}

	// 根据 actual 类型, 创建一个合适的 object 资源对象(如果 actual 为 nil 的话)
	// use the target if necessary
	obj, err := runtime.UseOrCreateObject(s.typer, s.creater, *actual, into)
	if err != nil {
		return nil, actual, err
	}

	// 将 data 数据解析为 obj 对象
	if err := caseSensitiveJsonIterator.Unmarshal(data, obj); err != nil {
		return nil, actual, err
	}

	// 如果传入的参数中, 开启了 strict 模式, 则还要进行如下步骤

	// If the deserializer is non-strict, return successfully here.
	if !s.options.Strict {
		return obj, actual, nil
	}

	// In strict mode pass the data trough the YAMLToJSONStrict converter.
	// This is done to catch duplicate fields regardless of encoding (JSON or YAML). 
	// For JSON data, the output would equal the input, 
	// unless there is a parsing error such as duplicate fields.
	// As we know this was successful in the non-strict case, 
	// the only error that may be returned here
	// is because of the newly-added strictness. 
	// hence we know we can return the typed strictDecoderError
	// the actual error is that the object contains duplicate fields.
	altered, err := yaml.YAMLToJSONStrict(originalData)
	if err != nil {
		return nil, actual, runtime.NewStrictDecodingError(err.Error(), string(originalData))
	}
	// As performance is not an issue for now for the strict deserializer 
	// (one has regardless to do the unmarshal twice), 
	// we take the sanitized, altered data that is guaranteed to have no duplicated
	// fields, and unmarshal this into a copy of the already-populated obj. 
	// Any error that occurs here is
	// due to that a matching field doesn't exist in the object. 
	// hence we can return a typed strictDecoderError,
	// the actual error is that the object contains unknown field.
	strictObj := obj.DeepCopyObject()
	if err := strictCaseSensitiveJsonIterator.Unmarshal(altered, strictObj); err != nil {
		return nil, actual, runtime.NewStrictDecodingError(err.Error(), string(originalData))
	}
	// Always return the same object as the non-strict serializer to avoid any deviations.
	return obj, actual, nil
}

// Encode serializes the provided object to the given writer.
func (s *Serializer) Encode(obj runtime.Object, w io.Writer) error {
	if co, ok := obj.(runtime.CacheableObject); ok {
		return co.CacheEncode(s.Identifier(), s.doEncode, w)
	}
	return s.doEncode(obj, w)
}

func (s *Serializer) doEncode(obj runtime.Object, w io.Writer) error {
	if s.options.Yaml {
		json, err := caseSensitiveJsonIterator.Marshal(obj)
		if err != nil {
			return err
		}
		data, err := yaml.JSONToYAML(json)
		if err != nil {
			return err
		}
		_, err = w.Write(data)
		return err
	}

	if s.options.Pretty {
		data, err := caseSensitiveJsonIterator.MarshalIndent(obj, "", "  ")
		if err != nil {
			return err
		}
		_, err = w.Write(data)
		return err
	}
	encoder := json.NewEncoder(w)
	return encoder.Encode(obj)
}

// Identifier implements runtime.Encoder interface.
func (s *Serializer) Identifier() runtime.Identifier {
	return s.identifier
}

// RecognizesData implements the RecognizingDecoder interface.
func (s *Serializer) RecognizesData(peek io.Reader) (ok, unknown bool, err error) {
	if s.options.Yaml {
		// we could potentially look for '---'
		return false, true, nil
	}
	_, _, ok = utilyaml.GuessJSONStream(peek, 2048)
	return ok, false, nil
}
