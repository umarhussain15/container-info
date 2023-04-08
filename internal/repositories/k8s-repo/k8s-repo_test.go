package k8srepo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func setupBefore() (*testclient.Clientset, K8sRepository, []v1.Pod) {
	clientset := testclient.NewSimpleClientset()
	repository := K8sRepository{client: clientset}

	pods := []v1.Pod{
		{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Pod",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod",
				Namespace: "default",
				Labels: map[string]string{
					"app.kubernetes.io/instance": "app1",
					"customLabel":                "label1",
				},
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:            "nginx",
						Image:           "nginx",
						ImagePullPolicy: "Always",
					},
				},
			},
		},
		{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Pod",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod2",
				Namespace: "default",
				Labels: map[string]string{
					"app.kubernetes.io/instance": "app2",
					"customLabel":                "label2",
				},
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "nginx",
						Image: "nginx",
						Resources: v1.ResourceRequirements{
							Limits: map[v1.ResourceName]resource.Quantity{
								v1.ResourceCPU:    resource.MustParse("4"),
								v1.ResourceMemory: resource.MustParse("8Gi"),
							},
							Requests: map[v1.ResourceName]resource.Quantity{
								v1.ResourceCPU:    resource.MustParse("2"),
								v1.ResourceMemory: resource.MustParse("3Gi"),
							},
						},
					},
					{
						Name:  "nginx2",
						Image: "nginx",
						Resources: v1.ResourceRequirements{
							Limits: map[v1.ResourceName]resource.Quantity{
								v1.ResourceCPU:    resource.MustParse("2"),
								v1.ResourceMemory: resource.MustParse("2Gi"),
							},
							Requests: map[v1.ResourceName]resource.Quantity{
								v1.ResourceCPU:    resource.MustParse("1"),
								v1.ResourceMemory: resource.MustParse("1Gi"),
							},
						},
					},
				},
			},
		},
	}

	for _, pod := range pods {
		create, err := clientset.CoreV1().Pods(pod.Namespace).Create(context.TODO(), &pod, metav1.CreateOptions{})
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Pod created %s\n", create.Name)
	}

	return clientset, repository, pods
}

func cleanUpAfter(clientset kubernetes.Interface) {
	err := clientset.CoreV1().Pods("").Delete(context.TODO(), "", metav1.DeleteOptions{})
	if err != nil {
		return
	}
}

func TestK8sRepository_SearchPods(t *testing.T) {
	searchLabel := "app.kubernetes.io/instance = app1"
	clientset, repository, pods := setupBefore()

	containers, err := repository.SearchPods(searchLabel)
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, 1, len(containers))

	searchMatch := "app.kubernetes.io/instance in (app1,app2)"
	containers, err = repository.SearchPods(searchMatch)
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, 3, len(containers))

	searchWithValues := "app.kubernetes.io/instance = app2"
	containers, err = repository.SearchPods(searchWithValues)
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, 2, len(containers))
	expectedPod := pods[1]
	for i, container := range containers {
		assert.NotEqual(t, expectedPod.Spec.Containers[i].Resources.Limits.Cpu(), container.CPULimit)
		assert.NotEqual(t, expectedPod.Spec.Containers[i].Resources.Requests.Cpu(), container.CPURequest)
		assert.NotEqual(t, expectedPod.Spec.Containers[i].Resources.Limits.Memory(), container.MemoryLimit)
		assert.NotEqual(t, expectedPod.Spec.Containers[i].Resources.Requests.Memory(), container.MemoryRequest)
	}

	cleanUpAfter(clientset)
}

func TestK8sRepository_SearchPods_MultipleLabels(t *testing.T) {
	multiLabels := "app.kubernetes.io/instance = app1,customLabel = label1"
	clientset, repository, _ := setupBefore()

	pods, err := repository.SearchPods(multiLabels)
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, 1, len(pods))
	multiLabels = "app.kubernetes.io/instance = app1,customLabel = label2"

	pods, err = repository.SearchPods(multiLabels)
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, 0, len(pods))

	cleanUpAfter(clientset)
}

func TestK8sRepository_SearchPods_NoMatch(t *testing.T) {
	searchLabel := "app.kubernetes.io/instance notin (app1,app2)"
	clientset, repository, _ := setupBefore()

	pods, err := repository.SearchPods(searchLabel)
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, 0, len(pods))

	cleanUpAfter(clientset)
}

func TestK8sRepository_SearchPods_ValuesNotPresent(t *testing.T) {
	searchLabel := "app.kubernetes.io/instance = app1"
	clientset, repository, _ := setupBefore()

	pods, err := repository.SearchPods(searchLabel)
	if err != nil {
		log.Fatalln(err)
	}
	assert.Equal(t, 1, len(pods))

	assert.Equal(t, "0", pods[0].CPULimit)
	assert.Equal(t, "0", pods[0].CPURequest)
	assert.Equal(t, "0", pods[0].MemoryLimit)
	assert.Equal(t, "0", pods[0].MemoryRequest)

	cleanUpAfter(clientset)
}
