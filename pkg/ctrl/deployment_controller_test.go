package ctrl

import (
	context "context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	testutil "github.com/yourusername/k8s-controller-tutorial/pkg/testutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestDeploymentReconciler_BasicFlow(t *testing.T) {
	mgr, k8sClient, _, cleanup := testutil.StartTestManager(t)
	defer cleanup()

	// Register the controller before starting the manager
	err := AddDeploymentController(mgr)
	require.NoError(t, err)

	go func() {
		_ = mgr.Start(context.Background())
	}()

	ns := "default"
	ctx := context.Background()
	name := "test-deployment"

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
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
	if err := k8sClient.Create(ctx, dep); err != nil {
		t.Fatalf("Failed to create Deployment: %v", err)
	}

	// Wait a bit to allow reconcile to be triggered
	time.Sleep(1 * time.Second)

	// Just check the Deployment still exists (reconcile didn't error or delete it)
	var got appsv1.Deployment
	err = k8sClient.Get(ctx, client.ObjectKey{Name: name, Namespace: ns}, &got)
	require.NoError(t, err)
}

func int32Ptr(i int32) *int32 { return &i }
