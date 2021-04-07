


package v1_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "generals-space/kube-operator/pkg/apis/kubegroup/v1"
	. "generals-space/kube-operator/pkg/client/clientset_generated/clientset/typed/kubegroup/v1"
)

var _ = Describe("PodCluster", func() {
	var instance PodCluster
	var expected PodCluster
	var client PodClusterInterface

	BeforeEach(func() {
		instance = PodCluster{}
		instance.Name = "instance-1"

		expected = instance
	})

	AfterEach(func() {
		client.Delete(context.TODO(), instance.Name, metav1.DeleteOptions{})
	})

	Describe("when sending a storage request", func() {
		Context("for a valid config", func() {
			It("should provide CRUD access to the object", func() {
				client = cs.KubegroupV1().PodClusters("podcluster-test-valid")

				By("returning success from the create request")
				actual, err := client.Create(context.TODO(), &instance, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				By("defaulting the expected fields")
				Expect(actual.Spec).To(Equal(expected.Spec))

				By("returning the item for list requests")
				result, err := client.List(context.TODO(), metav1.ListOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Items).To(HaveLen(1))
				Expect(result.Items[0].Spec).To(Equal(expected.Spec))

				By("returning the item for get requests")
				actual, err = client.Get(context.TODO(), instance.Name, metav1.GetOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(actual.Spec).To(Equal(expected.Spec))

				By("deleting the item for delete requests")
				err = client.Delete(context.TODO(), instance.Name, metav1.DeleteOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				result, err = client.List(context.TODO(), metav1.ListOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Items).To(HaveLen(0))
			})
		})
	})
})
