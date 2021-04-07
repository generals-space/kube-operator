


package v1_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/apiserver-builder-alpha/pkg/test"
	"k8s.io/client-go/rest"

	"generals-space/kube-operator/pkg/apis"
	"generals-space/kube-operator/pkg/client/clientset_generated/clientset"
	"generals-space/kube-operator/pkg/openapi"
)

var testenv *test.TestEnvironment
var config *rest.Config
var cs *clientset.Clientset

func TestV1(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "v1 Suite", []Reporter{test.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	testenv = test.NewTestEnvironment(apis.GetAllApiBuilders(), openapi.GetOpenAPIDefinitions)
	config = testenv.Start()
	cs = clientset.NewForConfigOrDie(config)
})

var _ = AfterSuite(func() {
	testenv.Stop()
})
