package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// StartTestManager sets up envtest, scheme, manager, and returns them with cleanup.
func StartTestManager(t *testing.T) (mgr manager.Manager, k8sClient client.Client, restCfg *rest.Config, cleanup func()) {
	t.Helper()
	testScheme := runtime.NewScheme()

	// Add the core Kubernetes schemes
	require.NoError(t, scheme.AddToScheme(testScheme))
	require.NoError(t, apiextensionsv1.AddToScheme(testScheme))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	env := &envtest.Environment{
		ErrorIfCRDPathMissing:    true,
		AttachControlPlaneOutput: false,
	}
	var startErr = make(chan error)
	var cfg *rest.Config
	var err error

	go func() {
		cfg, err = env.Start()
		startErr <- err
	}()

	// Wait for environment to start with timeout
	select {
	case err := <-startErr:
		require.NoError(t, err, "Failed to start test environment")
	case <-ctx.Done():
		t.Fatal("Timeout waiting for test environment to start")
	}

	require.NotNil(t, cfg)

	mgr, err = manager.New(cfg, manager.Options{Scheme: testScheme, LeaderElection: false})
	require.NoError(t, err)

	ctx, cancel = context.WithCancel(context.Background())
	go func() {
		_ = mgr.Start(ctx)
	}()

	k8sClient = mgr.GetClient()

	cleanup = func() {
		cancel()
		_ = env.Stop()
	}
	return mgr, k8sClient, cfg, cleanup
}

// SetupEnv starts envtest, creates a clientset, populates the cluster with sample Deployments, and returns env, clientset, and cleanup.
func SetupEnv(t *testing.T) (*envtest.Environment, *kubernetes.Clientset, func()) {
	t.Helper()
	testScheme := runtime.NewScheme()

	// Add the core Kubernetes schemes
	err := scheme.AddToScheme(testScheme)
	require.NoError(t, err)

	// Create a longer context timeout for environment startup
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Add CRD scheme
	err = apiextensionsv1.AddToScheme(testScheme)
	require.NoError(t, err)

	env := &envtest.Environment{
		ErrorIfCRDPathMissing:    true,
		AttachControlPlaneOutput: false,
	}
	var startErr = make(chan error)
	var cfg *rest.Config

	go func() {
		cfg, err = env.Start()
		startErr <- err
	}()

	// Wait for environment to start with timeout
	select {
	case err := <-startErr:
		require.NoError(t, err, "Failed to start test environment")
	case <-ctx.Done():
		t.Fatal("Timeout waiting for test environment to start")
	}

	require.NotNil(t, cfg)

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