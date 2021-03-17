package pod

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const queueProxyImage = "gcr.io/knative-releases/knative.dev/serving/cmd/queue@sha256:5481fae4ac4973bffb0d8c455a05f0e6c514cea7976d3add8f24081ab479e0b5"

func Knative(ns, name string) *corev1.Pod {
	var (
		grace int64 = 300
		tru         = true
		fal         = false
	)
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      name,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:            "user-container",
				Image:           basicApplicationImage,
				ImagePullPolicy: corev1.PullIfNotPresent,
				Env: []corev1.EnvVar{{
					Name:  "PORT",
					Value: "8080",
				}, {
					Name:  "K_REVISION",
					Value: "helloworld-go-00001",
				}, {
					Name:  "K_CONFIGURATION",
					Value: "helloworld-go",
				}, {
					Name:  "K_SERVICE",
					Value: "helloworld-go",
				}},
				Ports: []corev1.ContainerPort{{
					Name:          "user-port",
					ContainerPort: 8080,
				}},
				Lifecycle: &corev1.Lifecycle{
					PreStop: &corev1.Handler{
						HTTPGet: &corev1.HTTPGetAction{
							Port: intstr.FromInt(8022),
							Path: "/wait-for-drain",
						},
					},
				},
				TerminationMessagePolicy: corev1.TerminationMessageFallbackToLogsOnError,
			}, {
				Name:            "queue-proxy",
				Image:           queueProxyImage,
				ImagePullPolicy: corev1.PullIfNotPresent,
				Env: []corev1.EnvVar{{
					Name:  "SERVING_NAMESPACE",
					Value: ns,
				}, {
					Name:  "SERVING_SERVICE",
					Value: "helloworld-go",
				}, {
					Name:  "SERVING_CONFIGURATION",
					Value: "helloworld-go",
				}, {
					Name:  "SERVING_REVISION",
					Value: "helloworld-go-00001",
				}, {
					Name:  "QUEUE_SERVING_PORT",
					Value: "8012",
				}, {
					Name:  "CONTAINER_CONCURRENCY",
					Value: "0",
				}, {
					Name:  "REVISION_TIMEOUT_SECONDS",
					Value: "300",
				}, {
					Name: "SERVING_POD",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							FieldPath: "metadata.name",
						},
					},
				}, {
					Name: "SERVING_POD_IP",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							FieldPath: "status.podIP",
						},
					},
				}, {
					Name:  "SERVING_LOGGING_CONFIG",
					Value: "",
				}, {
					Name:  "SERVING_LOGGING_LEVEL",
					Value: "",
				}, {
					Name: "SERVING_REQUEST_LOG_TEMPLATE",
					Value: `{"httpRequest": {"requestMethod": "{{.Request.Method}}", "requestUrl":
					"{{js .Request.RequestURI}}", "requestSize": "{{.Request.ContentLength}}",
					"status": {{.Response.Code}}, "responseSize": "{{.Response.Size}}", "userAgent":
					"{{js .Request.UserAgent}}", "remoteIp": "{{js .Request.RemoteAddr}}", "serverIp":
					"{{.Revision.PodIP}}", "referer": "{{js .Request.Referer}}", "latency":
					"{{.Response.Latency}}s", "protocol": "{{.Request.Proto}}"}, "traceId":
					"{{index .Request.Header "X-B3-Traceid"}}"}`,
				}, {
					Name:  "SERVING_ENABLE_REQUEST_LOG",
					Value: "false",
				}, {
					Name:  "SERVING_REQUEST_METRICS_BACKEND",
					Value: "prometheus",
				}, {
					Name:  "TRACING_CONFIG_BACKEND",
					Value: "none",
				}, {
					Name:  "TRACING_CONFIG_ZIPKIN_ENDPOINT",
					Value: "",
				}, {
					Name:  "TRACING_CONFIG_STACKDRIVER_PROJECT_ID",
					Value: "",
				}, {
					Name:  "TRACING_CONFIG_DEBUG",
					Value: "false",
				}, {
					Name:  "TRACING_CONFIG_SAMPLE_RATE",
					Value: "0.1",
				}, {
					Name:  "USER_PORT",
					Value: "8080",
				}, {
					Name:  "SYSTEM_NAMESPACE",
					Value: "knative-serving",
				}, {
					Name:  "METRICS_DOMAIN",
					Value: "knative.dev/internal/serving",
				}, {
					Name:  "SERVING_READINESS_PROBE",
					Value: `{"tcpSocket":{"port":8080,"host":"127.0.0.1"},"successThreshold":1}`,
				}, {
					Name:  "ENABLE_PROFILING",
					Value: "false",
				}, {
					Name:  "SERVING_ENABLE_PROBE_REQUEST_LOG",
					Value: "false",
				}, {
					Name:  "METRICS_COLLECTOR_ADDRESS",
					Value: "",
				}},
				Ports: []corev1.ContainerPort{{
					Name:          "http-queueadm",
					ContainerPort: 8022,
				}, {
					Name:          "http-autometric",
					ContainerPort: 9090,
				}, {
					Name:          "http-usermetric",
					ContainerPort: 9091,
				}, {
					Name:          "queue-port",
					ContainerPort: 8012,
				}},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU: resource.MustParse("25m"),
					},
				},
				ReadinessProbe: &corev1.Probe{
					Handler: corev1.Handler{
						Exec: &corev1.ExecAction{
							Command: []string{"/ko-app/queue", "-probe-period", "0"},
						},
					},
					TimeoutSeconds:   10,
					FailureThreshold: 3,
					SuccessThreshold: 1,
					PeriodSeconds:    10,
				},
				SecurityContext: &corev1.SecurityContext{
					AllowPrivilegeEscalation: &fal,
					ReadOnlyRootFilesystem:   &tru,
					RunAsNonRoot:             &tru,
					Capabilities: &corev1.Capabilities{
						Drop: []corev1.Capability{"all"},
					},
				},
			}},
			TerminationGracePeriodSeconds: &grace,
		},
	}
}
