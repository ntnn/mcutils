package mcutils

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	mccontroller "sigs.k8s.io/multicluster-runtime/pkg/controller"
	mchandler "sigs.k8s.io/multicluster-runtime/pkg/handler"
	"sigs.k8s.io/multicluster-runtime/pkg/multicluster"
	mcreconcile "sigs.k8s.io/multicluster-runtime/pkg/reconcile"
	mcsource "sigs.k8s.io/multicluster-runtime/pkg/source"
)

// UnmanagedController creates an unmanaged controller and configures a multicluster watch for the given object and predicates.
func UnmanagedController(
	name string,
	obj client.Object,
	clName multicluster.ClusterName,
	predicates []predicate.TypedPredicate[client.Object],
	opts mccontroller.Options,
) (mccontroller.TypedController[mcreconcile.Request], error) {
	handler := mchandler.TypedEnqueueRequestForObject[client.Object]()
	source := mcsource.TypedKind(obj, handler, predicates...)

	if clName.String() != "" {
		source = source.WithClusterFilter(func(name multicluster.ClusterName, _ cluster.Cluster) bool {
			return name == clName
		})
	}

	c, err := mccontroller.NewUnmanaged(name, nil, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create unmanaged controller for node %s: %w", name, err)
	}

	if err := c.MultiClusterWatch(source); err != nil {
		return nil, fmt.Errorf("failed to start watch: %w", err)
	}

	return c, nil
}
