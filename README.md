# operator-toolkit

Operator-toolkit is a framework and collection of helpers designed to facilitate the development of Kubernetes
operators. With the goal of promoting code homogeneity across different operators, the toolkit offers a comprehensive
set of tools and best practices. By leveraging operator-toolkit, you can streamline the construction of Kubernetes operators and ensure consistent coding
patterns and practices across your projects.

For detailed documentation, guides, and examples, please refer to the
[project wiki](https://github.com/konflux-ci/operator-toolkit/wiki). The wiki serves as a
central resource for comprehensive information on the toolkit's various components and features. It provides detailed
explanations, code samples, and best practices that will help you effectively utilize the toolkit to construct
high-quality Kubernetes operators.

Key Features of Operator Toolkit:
- Abstractions for defining controllers, and reconcilers.
- Predicates for filtering and handling resource updates.
- Utilities for managing conditions, metrics, and event filtering.
- Seamless integration with [Operator SDK](https://sdk.operatorframework.io/) and [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime).
- Extensive documentation and examples.


## Testing with Mock Client

The operator-toolkit provides comprehensive mocking capabilities for unit testing. You can mock both loader operations (Get) and all client operations (Create, Update, Patch, Delete).

### Basic Usage

#### Mocking Client Operations

```go
import (
    "github.com/konflux-ci/operator-toolkit/loader"
    "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/errors"
)

// Mock a Create operation that returns an error
mockCtx, mockClient := loader.GetMockedContextWithClient(
    ctx,
    realK8sClient,
    []loader.MockData{}, // loader mocks (legacy support)
    []loader.ClientCallMock{
        {
            Operation:  loader.OperationCreate,
            ObjectType: &v1.Pod{},
            Err:        errors.NewBadRequest("pod creation failed"),
        },
    },
)

// Now when code calls Create on a Pod, it will return the mocked error
err := mockClient.Create(mockCtx, myPod)
// err will be "pod creation failed"
```

#### Supported Operations

- `OperationCreate` - Mock `client.Create()`
- `OperationUpdate` - Mock `client.Update()`
- `OperationPatch` - Mock `client.Patch()`
- `OperationDelete` - Mock `client.Delete()`
- `OperationGet` - Mock `client.Get()`
- `OperationList` - Mock `client.List()`
- `OperationDeleteAllOf` - Mock `client.DeleteAllOf()`

#### Mocking Status Updates

```go
mockCtx, mockClient := loader.GetMockedContextWithClient(
    ctx,
    realK8sClient,
    []loader.MockData{},
    []loader.ClientCallMock{
        {
            Operation:       loader.OperationUpdate,
            ObjectType:      &v1.Pod{},
            SubResourceName: "status",
            Err:             errors.NewForbidden(schema.GroupResource{}, "pod", nil),
        },
    },
)

// This will return the mocked error
err := mockClient.Status().Update(mockCtx, myPod)
```

#### Combining Loader and Client Mocks

```go
mockCtx, mockClient := loader.GetMockedContextWithClient(
    ctx,
    realK8sClient,
    []loader.MockData{
        {
            ContextKey: loader.ApplicationContextKey,
            Resource:   myApplication,
        },
    },
    []loader.ClientCallMock{
        {
            Operation:  loader.OperationCreate,
            ObjectType: &v1beta2.IntegrationTestScenario{},
            Err:        errors.NewServerTimeout("timeout"),
        },
    },
)

// Loader mock works
app, err := loader.GetMockedResourceAndErrorFromContext(mockCtx, loader.ApplicationContextKey, myApplication)

// Client mock works
err = mockClient.Create(mockCtx, newScenario)
```

### Testing Different Object Types

The mock client uses type matching, so you can mock different errors for different resource types:

```go
mockCtx, mockClient := loader.GetMockedContextWithClient(
    ctx,
    realK8sClient,
    []loader.MockData{},
    []loader.ClientCallMock{
        {
            Operation:  loader.OperationCreate,
            ObjectType: &v1.Pod{},
            Err:        errors.NewBadRequest("pod error"),
        },
        {
            Operation:  loader.OperationCreate,
            ObjectType: &v1.ConfigMap{},
            Err:        errors.NewConflict(schema.GroupResource{}, "cm", nil),
        },
    },
)

// Each type gets its specific mocked error
err1 := mockClient.Create(mockCtx, myPod)        // returns "pod error"
err2 := mockClient.Create(mockCtx, myConfigMap)  // returns conflict error
```

### Fallback to Real Client

When no mock matches an operation, the real client is used. This allows you to mock only specific operations while letting others proceed normally.
