---
name: generate-unit-tests
description: Generate a Go unit test file for a target source file, maximizing code coverage with one test per function/responsibility. Follows TDD philosophy — writes tests even for unimplemented code (expecting compile errors). Use when the user asks to create tests, generate tests, or write unit tests for a Go file.
---

# Generate Unit Tests

## Workflow

1. **Read the target file** to understand all exported and unexported functions, their signatures, dependencies, and error paths.
2. **Identify all dependencies** (interfaces, services, repositories) the target file uses.
3. **Locate existing mocks** by searching `**/*_mock*` in the codebase. Reuse them when available.
4. **Read at least one existing `_test.go`** in the project to confirm current conventions before writing.
5. **Create the test file** as `<target_filename>_test.go` in the same directory and package as the target file.

## Test Structure Conventions

### File setup

- Same package as the target (no `_test` suffix on the package name).
- Imports use the standard library `testing` package only — no third-party assertion libraries.
- Group imports: stdlib, then project packages, separated by a blank line.

### Helper constructors

Create `newTest<ServiceName>(...)` factory functions to build the service under test with mocks. Provide variants when tests need different initial state (e.g. `newTestAuthServiceEmptyUsers`, `newTestAuthServiceWithUsers`).

```go
func newTestUserSettingsService() *UserSettingsService {
    mockAI := &aiservice.AIServiceMock{}
    mockRepo := &user_repo.UserRepositoryMock{Users: make(map[string]*user_repo.User)}
    return NewUserSettingsService(mockRepo, mockAI)
}
```

### Naming

- Test function names: `Test<FunctionName><Scenario>` — e.g. `TestSignUpSuccess`, `TestSignUpInvalidEmail`.
- Table-driven sub-tests use `t.Run(tc.name, ...)` for each case.

### Assertions

- Use `t.Fatalf` when the test cannot continue (nil check before field access, unexpected error on setup).
- Use `t.Errorf` for non-fatal checks that should still report failures.
- Format: `t.Errorf("<FunctionName>(%q); Wanted <expected>, got <actual>", ...)`.
- Always compare with sentinel errors via `errors.Is(err, ExpectedError)` when applicable.

### Table-driven tests

Use table-driven tests when a function has multiple input/output combinations for the same logical assertion:

```go
func TestLoginInvalidCredentials(t *testing.T) {
    cases := []struct {
        name            string
        emailOrUsername string
        password        string
        expectedError   error
    }{
        {"wrong email", "bad@example.com", "Pass123", ErrInvalidCredentials},
        {"wrong password", "email@example.com", "wrong", ErrInvalidCredentials},
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            // ...
        })
    }
}
```

### Mocks with function fields

When a mock uses injectable function fields (like `AIServiceMock`), set the desired behavior at the start of each test or sub-test:

```go
mockAI.PromptFunc = func(prompt string) (string, error) {
    return "mocked response", nil
}
```

## Coverage Rules

Aim for at least **80% coverage of meaningful scenarios** while keeping the number of test cases minimal. Avoid redundant cases that test the same code path with trivially different inputs. Each case in a table-driven test must exercise a **distinct branch or behavior**.

For every function in the target file, generate tests for:

1. **Happy path** — expected successful behavior.
2. **Each distinct error return** — one test (or table-driven case) per error sentinel / error condition.
3. **Edge cases** — only when they exercise a different code path (boundary values, empty inputs, nil returns).

If two inputs trigger the exact same branch and validation, keep only the most representative one.

## TDD Philosophy

Analyze the target file and identify **functions that logically should exist but are not yet implemented**. For example, if a service has `SignUp` and `Logout` but no `Login`, the test file should include tests for `Login` as if it already existed — calling it with the expected signature and asserting on the expected return values. The resulting compile errors guide the next implementation step (Red → Green → Refactor).

This is about **completeness of the domain logic**, not about inventing arbitrary functions. Only write tests for functions that clearly belong in the file based on the existing patterns and responsibilities.

- Write the test calling the missing function with the expected signature.
- Assert on the expected return values.
- Do NOT stub or skip — let the compiler error guide the implementation.

## Checklist

Before finishing, verify:

- [ ] One test file created: `<target>_test.go` in the same package
- [ ] Every public function in the target has at least one test
- [ ] Every known error path has a dedicated test or table-driven case
- [ ] Helper constructors created for service instantiation with mocks
- [ ] Existing mocks reused — no duplicate mock types created
- [ ] No third-party test libraries used
- [ ] TDD: tests for unimplemented functions are present and expected to fail at compile time
