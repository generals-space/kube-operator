// Code generated by apiregister-gen. DO NOT EDIT.

package install

import (
	aggregatedclientset "generals-space/kube-operator/pkg/client/clientset_generated/clientset"
	aggregatedinformerfactory "generals-space/kube-operator/pkg/client/informers_generated/externalversions"
	initializer "generals-space/kube-operator/plugin/admission"
	. "generals-space/kube-operator/plugin/admission/podcluster"

	"k8s.io/apiserver/pkg/admission"
	genericserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/apiserver-builder-alpha/pkg/cmd/server"
)

func init() {
	server.AggregatedAdmissionInitializerGetter = GetAggregatedResourceAdmissionControllerInitializer
	server.AggregatedAdmissionPlugins["PodCluster"] = NewPodClusterPlugin()

}

func GetAggregatedResourceAdmissionControllerInitializer(config *rest.Config) (admission.PluginInitializer, genericserver.PostStartHookFunc) {
	// init aggregated resource clients
	aggregatedResourceClient := aggregatedclientset.NewForConfigOrDie(config)
	aggregatedInformerFactory := aggregatedinformerfactory.NewSharedInformerFactory(aggregatedResourceClient, 0)
	aggregatedResourceInitializer := initializer.New(aggregatedResourceClient, aggregatedInformerFactory)

	return aggregatedResourceInitializer, func(context genericserver.PostStartHookContext) error {
		aggregatedInformerFactory.Start(context.StopCh)
		return nil
	}
}
