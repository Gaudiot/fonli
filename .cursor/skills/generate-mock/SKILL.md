---
name: generate-mock
description: Generate or update a Go mock file from an interface definition. Creates a mock struct with injectable function fields, default zero-value returns, and call counters for each method. Use when the user asks to create a mock, generate a mock, update a mock, or write a mock for a Go interface.
---

# Generate Mock

## Workflow

1. **Read the source file** containing the interface to mock.
2. **Identify the interface** and all its methods, including parameter types and return types.
3. **Determine the file name** from the source file pattern (see File Naming section below).
4. **Check if the mock file already exists** — if so, update it preserving any custom helper functions defined outside the mock struct and its methods.
5. **Generate or update the mock file** following the structure below.
6. **Run `go vet`** on the generated file to verify correctness.

## File Naming

Derive the mock file name from the source file using two rules:

1. If the source file ends with `.interface.go`, replace `.interface.go` with `.mock.go`.
2. Otherwise, insert `_mock` before `.go`.

| Source file | Mock file | Rule applied |
|---|---|---|
| `user_repository.interface.go` | `user_repository.mock.go` | `.interface.go` → `.mock.go` |
| `token_service.interface.go` | `token_service.mock.go` | `.interface.go` → `.mock.go` |
| `ai_service.interface.go` | `ai_service.mock.go` | `.interface.go` → `.mock.go` |
| `payment.gateway.go` | `payment.gateway_mock.go` | insert `_mock` before `.go` |

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
| `error` | **see below** |
| Named struct (value type) | `StructName{}` |

### Error Return Behavior

If any return type is `error`, the default (when `Func` is nil) must return `errors.New("[Mock] not implemented")` for that error value. All other return values stay at their zero value.

This makes forgotten mock setups immediately visible — instead of silently returning `nil` and letting tests pass incorrectly, the mock loudly fails with a clear message.

| Signature | Default return |
|---|---|
| `Delete(id string) error` | `errors.New("[Mock] not implemented")` |
| `GetByID(id string) (*Item, error)` | `nil, errors.New("[Mock] not implemented")` |
| `Count() int` | `0` (no error in signature — unchanged) |

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

import "errors"

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
    return nil, errors.New("[Mock] not implemented")
}

func (m *MyServiceMock) Delete(id string) error {
    m.DeleteCallCount++
    if m.DeleteFunc != nil {
        return m.DeleteFunc(id)
    }
    return errors.New("[Mock] not implemented")
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
- [ ] File name follows the naming convention (`.interface.go` → `.mock.go`, otherwise `_mock.go`)
- [ ] Every interface method has a corresponding `Func` field, `CallCount` field, and method implementation
- [ ] Zero-value defaults are correct for all return types
- [ ] `go vet` passes on the generated file
- [ ] No duplicate types or imports
