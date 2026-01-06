## 2026-01-06 - [Refactoring PostProcess for Testability]
**Pattern:** [Hardcoded dependencies] The `internal/postprocess` package had a hardcoded dependency on `config.SettingsManager` and `http.Client`, making it impossible to test without side effects (file I/O, network).
**Refactor:** [Dependency Injection]
1. Introduced `SettingsProvider` interface to abstract configuration access.
2. Updated `New` to accept the interface.
3. Added `SetHTTPClient` to allow injecting a mock transport for the HTTP client.
**Action:** [Guidance for future tests] When testing components that depend on global configuration or network calls, always define an interface for the configuration and provide a way to inject or replace the network client. Use `httptest` or custom `RoundTripper` for network mocking.
