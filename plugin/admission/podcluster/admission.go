


package podclusteradmission

import (
	"context"
	aggregatedadmission "generals-space/kube-operator/plugin/admission"
	aggregatedinformerfactory "generals-space/kube-operator/pkg/client/informers_generated/externalversions"
	aggregatedclientset "generals-space/kube-operator/pkg/client/clientset_generated/clientset"
	genericadmissioninitializer "k8s.io/apiserver/pkg/admission/initializer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/apiserver/pkg/admission"
)

var _ admission.Interface 											= &podclusterPlugin{}
var _ admission.MutationInterface 									= &podclusterPlugin{}
var _ admission.ValidationInterface 								= &podclusterPlugin{}
var _ genericadmissioninitializer.WantsExternalKubeInformerFactory 	= &podclusterPlugin{}
var _ genericadmissioninitializer.WantsExternalKubeClientSet 		= &podclusterPlugin{}
var _ aggregatedadmission.WantsAggregatedResourceInformerFactory 	= &podclusterPlugin{}
var _ aggregatedadmission.WantsAggregatedResourceClientSet 			= &podclusterPlugin{}

func NewPodClusterPlugin() *podclusterPlugin {
	return &podclusterPlugin{
		Handler: admission.NewHandler(admission.Create, admission.Update),
	}
}

type podclusterPlugin struct {
	*admission.Handler
}

func (p *podclusterPlugin) ValidateInitialization() error {
	return nil
}

func (p *podclusterPlugin) Admit(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) error {
	return nil
}

func (p *podclusterPlugin) Validate(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) error {
	return nil
}

func (p *podclusterPlugin) SetAggregatedResourceInformerFactory(aggregatedinformerfactory.SharedInformerFactory) {}

func (p *podclusterPlugin) SetAggregatedResourceClientSet(aggregatedclientset.Interface) {}

func (p *podclusterPlugin) SetExternalKubeInformerFactory(informers.SharedInformerFactory) {}

func (p *podclusterPlugin) SetExternalKubeClientSet(kubernetes.Interface) {}
