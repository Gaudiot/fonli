---
name: generate-mock
description: Generate or update a Go mock file from an interface definition. Creates a mock struct with injectable function fields, default zero-value returns, and call counters for each method. Use when the user asks to create a mock, generate a mock, update a mock, or write a mock for a Go interface.
---

# Generate Mock

## Workflow

1. **Read the source file** containing the interface to mock.
2. **Identify the interface** and all its methods, including parameter types and return types.
3. **Determine the file name** from the source file pattern: if the source is `user.repository.go`, the mock file is `user.repository_mock.go` in the same directory.
4. **Check if the mock file already exists** — if so, update it preserving any custom helper functions defined outside the mock struct and its methods.
5. **Generate or update the mock file** following the structure below.
6. **Run `go vet`** on the generated file to verify correctness.

## File Naming

Derive the mock file name from the source file:

| Source file | Mock file |
|---|---|
| `user.repository.go` | `user.repository_mock.go` |
| `ai_service.go` | `ai_service_mock.go` |
| `payment.gateway.go` | `payment.gateway_mock.go` |

Rule: insert `_mock` before `.go` in the source file name.

## Mock Structure

The mock struct has three kinds of fields per interface method:

1. **`<Method>Func`** — injectable `func` field matching the method signature
2. **`<Method>CallCount`** — `int` counter incremented on every call

Each method implementation:
1. Increments `CallCount`
2. Delegates to `Func` if set
3. Returns zero values otherwise

### Zero-Value Defaults

When `<Method>Func` is nil, return the zero value for each return type:

| Type | Zero value |
|---|---|
| `int`, `int8/16/32/64`, `uint`, `float32/64` | `0` |
| `string` | `""` |
| `bool` | `false` |
| Pointer, slice, map, interface, func, chan | `nil` |
| `error` | `nil` |
| Named struct (value type) | `StructName{}` |

## Template

Given this interface:

```go
package mypackage

type MyService interface {
    GetByID(id string) (*Item, error)
    Delete(id string) error
    Count() int
}
```

Generate this mock:

```go
package mypackage

type MyServiceMock struct {
    GetByIDFunc    func(id string) (*Item, error)
    GetByIDCallCount int

    DeleteFunc    func(id string) error
    DeleteCallCount int

    CountFunc    func() int
    CountCallCount int
}

func (m *MyServiceMock) GetByID(id string) (*Item, error) {
    m.GetByIDCallCount++
    if m.GetByIDFunc != nil {
        return m.GetByIDFunc(id)
    }
    return nil, nil
}

func (m *MyServiceMock) Delete(id string) error {
    m.DeleteCallCount++
    if m.DeleteFunc != nil {
        return m.DeleteFunc(id)
    }
    return nil
}

func (m *MyServiceMock) Count() int {
    m.CountCallCount++
    if m.CountFunc != nil {
        return m.CountFunc()
    }
    return 0
}
```

## Rules

- Same package as the interface — no separate `_mock` or `mocks` package.
- Mock struct name: `<InterfaceName>Mock` (e.g. `UserRepository` -> `UserRepositoryMock`).
- Receiver name: `m` (consistent across all mocks).
- Only add imports that the mock methods actually need. Do not import packages only used by the interface file.
- Do NOT add comments that narrate what the code does.
- If the source file contains types (structs, error vars, etc.) alongside the interface, do NOT duplicate them in the mock file.
- Separate each group of fields (Func + CallCount) with a blank line inside the struct.
- Method implementations follow the order: increment counter, check Func, return zero values.

## Updating Existing Mocks

When a mock file already exists:

1. Read the current mock file.
2. Compare methods in the interface vs methods in the mock.
3. **Add** new methods that exist in the interface but not in the mock.
4. **Remove** methods that no longer exist in the interface.
5. **Update** method signatures that changed (parameter types, return types).
6. **Preserve** any helper functions or variables defined in the mock file outside the struct and its method implementations.

## Checklist

Before finishing, verify:

- [ ] Mock file created/updated in the same directory as the interface
- [ ] File name follows `<source>_mock.go` convention
- [ ] Every interface method has a corresponding `Func` field, `CallCount` field, and method implementation
- [ ] Zero-value defaults are correct for all return types
- [ ] `go vet` passes on the generated file
- [ ] No duplicate types or imports
