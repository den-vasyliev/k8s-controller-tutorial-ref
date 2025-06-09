package informer

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	testutil "github.com/yourusername/k8s-controller-tutorial/pkg/testutil"
)

func TestStartDeploymentInformer(t *testing.T) {
	_, clientset, cleanup := testutil.SetupEnv(t)
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	// Patch log to write to a buffer or just rely on test output
	added := make(chan string, 2)

	// Patch event handler for test
	factory := informers.NewSharedInformerFactoryWithOptions(
		clientset,
		30*time.Second,
		informers.WithNamespace("default"),
	)
	informer := factory.Apps().V1().Deployments().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if d, ok := obj.(metav1.Object); ok {
				added <- d.GetName()
			}
		},
	})

	go func() {
		defer wg.Done()
		factory.Start(ctx.Done())
		factory.WaitForCacheSync(ctx.Done())
		<-ctx.Done()
	}()

	// Wait for events
	found := map[string]bool{}
	for range 2 {
		select {
		case name := <-added:
			found[name] = true
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for deployment add events")
		}
	}

	require.True(t, found["sample-deployment-1"])
	require.True(t, found["sample-deployment-2"])

	cancel()
	wg.Wait()
}
