/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package conversion

import (
	"fmt"
	"reflect"
)

// EnforcePtr 确认 obj 对象为指针类型, 然后返回该指针的引用类型.
// 如 obj 为 &Object{...}, 则会返回 Object{...}, 但要注意, 返回值为 reflect.Value 类型.
//
// EnforcePtr ensures that obj is a pointer of some sort. 
// Returns a reflect.Value of the dereferenced pointer, 
// ensuring that it is settable/addressable.
//
// Returns an error if this is not possible.
func EnforcePtr(obj interface{}) (reflect.Value, error) {
	// v 是一个 reflect.Value{} 对象
	v := reflect.ValueOf(obj) 
	// 如果 obj 是一个指针类型变量, 那么 v.Kind() 必然是 reflect.Ptr.
	if v.Kind() != reflect.Ptr {
		if v.Kind() == reflect.Invalid {
			return reflect.Value{}, fmt.Errorf("expected pointer, but got invalid kind")
		}
		return reflect.Value{}, fmt.Errorf("expected pointer, but got %v type", v.Type())
	}
	// nil 不是指针, 但也没办法取指针.
	if v.IsNil() {
		return reflect.Value{}, fmt.Errorf("expected pointer, but got nil")
	}
	// 运行到这里, 说明 obj 是一个指针对象, Elem() 方法对其取指针, 得到其引用类型
	return v.Elem(), nil
}
