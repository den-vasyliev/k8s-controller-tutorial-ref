package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

// SetupEnv starts envtest, creates a clientset, populates the cluster with sample Deployments, and returns env, clientset, and cleanup.
func SetupEnv(t *testing.T) (*envtest.Environment, *kubernetes.Clientset, func()) {
	t.Helper()
	ctx := context.Background()
	env := &envtest.Environment{}

	cfg, err := env.Start()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Write kubeconfig to /tmp/envtest.kubeconfig
	kubeconfig := clientcmdapi.NewConfig()
	kubeconfig.Clusters["envtest"] = &clientcmdapi.Cluster{
		Server:                   cfg.Host,
		CertificateAuthorityData: cfg.CAData,
	}
	kubeconfig.AuthInfos["envtest-user"] = &clientcmdapi.AuthInfo{
		ClientCertificateData: cfg.CertData,
		ClientKeyData:         cfg.KeyData,
	}
	kubeconfig.Contexts["envtest-context"] = &clientcmdapi.Context{
		Cluster:  "envtest",
		AuthInfo: "envtest-user",
	}
	kubeconfig.CurrentContext = "envtest-context"

	kubeconfigBytes, err := clientcmd.Write(*kubeconfig)
	require.NoError(t, err)
	err = os.WriteFile("/tmp/envtest.kubeconfig", kubeconfigBytes, 0644)
	require.NoError(t, err)

	clientset, err := kubernetes.NewForConfig(cfg)
	require.NoError(t, err)

	// Create sample Deployments
	for i := 1; i <= 2; i++ {
		dep := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("sample-deployment-%d", i),
				Namespace: "default",
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "test"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "test"}},
					Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "nginx", Image: "nginx"}}},
				},
			},
		}
		_, err := clientset.AppsV1().Deployments("default").Create(ctx, dep, metav1.CreateOptions{})
		require.NoError(t, err)
	}

	cleanup := func() {
		_ = env.Stop()
	}
	return env, clientset, cleanup
}

func int32Ptr(i int32) *int32 { return &i }
