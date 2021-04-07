


package main

import (
	// Make sure dep tools picks up these dependencies
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "github.com/go-openapi/loads"

	"sigs.k8s.io/apiserver-builder-alpha/pkg/cmd/server"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // Enable cloud provider auth

	"generals-space/kube-operator/pkg/apis"
	"generals-space/kube-operator/pkg/openapi"
	_ "generals-space/kube-operator/plugin/admission/install"
)

func main() {
	version := "v0"

	err := server.StartApiServerWithOptions(&server.StartOptions{
		EtcdPath:         "/registry/generals.space",
		Apis:             apis.GetAllApiBuilders(),
		Openapidefs:      openapi.GetOpenAPIDefinitions,
		Title:            "Api",
		Version:          version,

		// TweakConfigFuncs []func(apiServer *apiserver.Config) error
		// FlagConfigFuncs []func(*cobra.Command) error
	})
	if err != nil {
		panic(err)
	}
}
