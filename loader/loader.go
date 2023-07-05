package loader

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	contextKey int

	// MockData represents the mocked data that will be returned by a mocked loader method.
	MockData struct {
		ContextKey contextKey
		Err        error
		Resource   any
	}
)

// GetObject loads an object from the cluster. This is a generic function that requires the object to be passed as an
// argument. The object is modified during the invocation.
func GetObject(name, namespace string, cli client.Client, ctx context.Context, object client.Object) error {
	return cli.Get(ctx, types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, object)
}

// GetMockedResourceAndErrorFromContext returns the mocked data found in the context passed as an argument. The data is
// to be found in the contextDataKey key. If not there, a panic will be raised.
func GetMockedResourceAndErrorFromContext[T any](ctx context.Context, contextKey contextKey, _ T) (T, error) {
	var resource T
	var err error

	value := ctx.Value(contextKey)
	if value == nil {
		panic("Mocked data not found in the context")
	}

	data, _ := value.(MockData)

	if data.Resource != nil {
		resource = data.Resource.(T)
	}

	if data.Err != nil {
		err = data.Err
	}

	return resource, err
}
