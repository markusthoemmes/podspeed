package template

import (
	"io"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func PodConstructorFromYAML(content io.Reader) (func(string, string) *corev1.Pod, error) {
	pod := &corev1.Pod{}

	decoder := yaml.NewYAMLToJSONDecoder(content)
	decoder.Decode(pod)

	return func(ns, name string) *corev1.Pod {
		// Reset metadata
		pod.ObjectMeta = metav1.ObjectMeta{
			Namespace: ns,
			Name:      name,
		}

		// Reset Status
		pod.Status = corev1.PodStatus{}

		return pod
	}, nil
}
