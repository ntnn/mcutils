package mcutils

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ntnn/mcutils/mctest"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	mctrl "sigs.k8s.io/multicluster-runtime"
	mccontroller "sigs.k8s.io/multicluster-runtime/pkg/controller"
	"sigs.k8s.io/multicluster-runtime/pkg/multicluster"
	mcreconcile "sigs.k8s.io/multicluster-runtime/pkg/reconcile"
)

func TestUnmanagedController(t *testing.T) {
	env := mctest.EnvTest(t, nil)

	watchObj := &unstructured.Unstructured{}
	watchObj.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Secret",
	})

	hit := atomic.Bool{}
	opts := mccontroller.Options{
		Reconciler: mcreconcile.Func(func(_ context.Context, req mctrl.Request) (mctrl.Result, error) {
			t.Logf("Reconciling %s/%s", req.Namespace, req.Name)
			if req.Name == "test-secret" && req.Namespace == "default" {
				hit.Store(true)
			}
			return mctrl.Result{}, nil
		}),
	}

	t.Log("Creating manager")
	mgr, err := mctrl.NewManager(env.Config, nil, SilentManagerOpts(mctrl.Options{}))
	require.NoError(t, err)

	t.Log("Creating unmanaged controller")
	c, err := UnmanagedController("test-controller", watchObj, "", nil, opts)
	require.NoError(t, err)

	// The controller shouldn't be added as a runnable to the manager
	// since it is supposed to be unmanaged, but this is the easiest way
	// to kickstart the process.
	mgr.Add(c)

	go func() {
		t.Log("Starting manager")
		require.NoError(t, mgr.Start(t.Context()))
	}()
	t.Log("Waiting for manager to be elected")
	<-mgr.Elected()

	t.Log("Emulating a provider by creating a cluster and engaging it with the manager")
	cl, err := cluster.New(env.Config, func(o *cluster.Options) {})
	require.NoError(t, err)

	go func() {
		require.NoError(t, cl.Start(t.Context()))
	}()

	clusterName := multicluster.ClusterName("test-cluster")
	err = mgr.Engage(t.Context(), clusterName, cl)
	require.NoError(t, err)

	t.Log("Creating test secret in the provided cluster")
	require.NoError(t, cl.GetClient().Create(t.Context(), &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name":      "test-secret",
				"namespace": "default",
			},
		},
	}))

	t.Log("Waiting for reconciler to be hit")
	require.Eventually(t, func() bool {
		return hit.Load()
	}, wait.ForeverTestTimeout, 100*time.Millisecond)
}
