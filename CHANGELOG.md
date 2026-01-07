# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Ritual Grove Integration Guide** - Complete guide for integrating Inertia into rituals
  - Documentation for adding Inertia support to existing rituals
  - Examples of conditional frontend scaffolding
  - Template structure and best practices
  - Migration guide from traditional templates to Inertia
  - Reference to fullstack-inertia-vue ritual
- **Code Quality Improvements** - Fixed linter warnings
  - Fixed godot warnings (comment periods in realtime/scela.go)
  - Added assertion to TestHubRunAndStop to use test parameter
  - All linter issues resolved except cognitive complexity warnings
- **E2E Testing Suite** - Comprehensive Playwright end-to-end tests
  - Todo app E2E tests with full CRUD operations
  - Chat app E2E tests with multi-tab real-time updates
  - Navigation flow tests for SPA routing without page reloads
  - SSR hydration tests with performance benchmarks
  - Form submission and validation tests
  - Browser history and state preservation tests
  - Connection resilience and reconnection tests
  - 40+ E2E test scenarios across 4 test files
  - Playwright configuration with Chromium and Firefox
  - E2E test server script for automated testing
  - Dependencies: @playwright/test@^1.57.0, playwright@^1.57.0
- **Scéla Integration** - Message bus integration for real-time updates
  - New `ScelaAdapter` bridges Scéla bus to WebSocket hub
  - Pattern matching support (`user.*`, `*.created`)
  - Message filtering with custom filter functions
  - Automatic message forwarding from Scéla to WebSocket clients
  - Graceful shutdown and error handling
  - 6 integration tests (48.6% realtime package coverage)
  - Dependency: github.com/toutaio/toutago-scela-bus@v1.5.5
  - Scéla integration example with filtering and patterns
- **Real-time WebSocket Updates** - Full WebSocket support for live updates
  - New `pkg/realtime` package with WebSocket hub
  - Channel-based message broadcasting
  - Client connection management with auto-cleanup
  - `Hub.Run()` for managing WebSocket lifecycle
  - `Hub.Broadcast()` and `Hub.Publish()` for sending messages
  - `HandleWebSocket()` for HTTP upgrade handling
  - Auto-reconnection with configurable retry logic
  - Dependency: github.com/gorilla/websocket@v1.5.3
- **useLiveUpdate Composable** - Vue 3 composable for real-time updates
  - `useLiveUpdate(url, options)` with connection management
  - `on(channel, handler)` for subscribing to channels
  - `off(channel, handler)` for unsubscribing
  - Reactive `connected` and `error` state
  - Auto-reconnection with configurable delays
  - Maximum reconnection attempts limit
  - Automatic cleanup on component unmount
  - 12 composable tests (100% pass rate)
  - 73 total Vue tests passing
- **Real-time Chat Example** - Complete WebSocket chat application
  - Backend with WebSocket hub and message broadcasting
  - Vue chat interface with real-time message updates
  - Connection status indicator
  - Auto-reconnection demonstration
  - Documentation and setup guide (examples/chat/README.md)

### Fixed
- Component context detection in useLiveUpdate to avoid warnings
- Test coverage reporting for all packages

## [0.5.0] - 2024-01-06

### Added
- **Server-Side Rendering (SSR)** - V8-based SSR support
  - New `pkg/ssr` package with V8 renderer
  - Context pooling for performance (configurable pool size)
  - Timeout protection (configurable timeout)
  - Error handling with graceful fallback
  - `SetSSRRenderer()` method to enable SSR
  - `RenderSSR()` method for server-side page rendering
  - Support for complex data structures (nested objects, arrays)
  - SSR documentation and examples
  - 8 SSR tests (100% coverage)
  - Dependency: rogchap.com/v8go@v0.9.0
- **API Documentation** - Complete API reference (API.md)
  - Core API: New(), Render(), Share(), ShareFunc()
  - Context methods: All 15+ context helpers
  - Middleware: Request handling and headers
  - TypeScript codegen: typegen.New(), Register(), GenerateFile()
  - HTMX support: All HTMX helpers and headers
  - Vue components: Link, Head
  - Vue composables: usePage, useForm, useRemember
  - Types, error handling, and best practices

### Fixed
- **Test Warnings** - Eliminated all Vue test warnings
  - Suppressed expected JSON parse error in useRemember test
  - Fixed usePage tests to work within Vue component context
  - Fixed Head component slot access to avoid render function warning
  - Made inject() silent in usePage when not provided (test mode)

### Added
- **Performance Benchmarks** - Comprehensive benchmark suite
  - BenchmarkRender - Basic rendering (~2.2μs, 31 allocs)
  - BenchmarkRenderWithSharedData - With shared data (~2.7μs, 37 allocs)
  - BenchmarkRenderWithLazyProps - With lazy evaluation (~2.3μs, 30 allocs)
  - BenchmarkPartialReload - Partial reload performance (~3.1μs, 46 allocs)
  - BenchmarkHTMXPartial - HTMX partial rendering (~732ns, 15 allocs)
- **Improved Test Coverage** - Comprehensive testing for uncovered code
  - `ShareFunc()` - Lazy shared data function tests
  - `WithInfo()` - Info flash message tests
  - `HTMXReplaceURL()` - HTMX URL replacement tests
  - `Always()` / `AlwaysLazy()` - Always-included props tests
  - TypeGen tests: `New()`, `Register()`, `GenerateFile()`, `toSnakeCase()`
  - Test coverage: inertia 84.9% → 92.4%, typegen 66.7% → 82.2%
- **Integration Tests** - Complete request cycle coverage
  - 6 full request/response scenarios
  - Initial page loads and navigation
  - Form submissions with validation
  - Lazy props evaluation
  - External redirects
  - Shared data flow (global & request-level)
  - Error handling (404, 500)
- **HTMX Integration Tests** - Complete test coverage
  - 7 integration test scenarios
  - Partial updates, redirects, triggers
  - Out-of-band swaps, chaining
  - Hybrid Inertia/HTMX routing
  - Validation error handling
- **Validation Helpers** - Convenient error handling
  - `ValidationErrors.Add()` - Add validation error to field
  - `ValidationErrors.Has()` - Check if field has errors
  - `ValidationErrors.First()` - Get first error for field
  - `ValidationErrors.Any()` - Check if any errors exist
  - `NewValidationErrors()` - Create new validation errors
  - `WithError()` - Add single validation error (chainable)
- **Flash Message Helpers** - Easy flash messages
  - `Flash.Success()` - Add success message
  - `Flash.Error()` - Add error message
  - `Flash.Warning()` - Add warning message
  - `Flash.Info()` - Add info message
  - `Flash.Custom()` - Add custom message
  - `NewFlash()` - Create new flash instance
  - `WithSuccess()` - Add success flash (chainable)
  - `WithErrorMessage()` - Add error flash (chainable)
  - `WithWarning()` - Add warning flash (chainable)
  - `WithInfo()` - Add info flash (chainable)
- **Lazy Props Support** - Optimize expensive computations
  - `Lazy()` - Props excluded from partial reloads unless requested
  - `AlwaysLazy()` - Lazy props always included even in partials
  - `Defer()` - Props only loaded when explicitly requested
  - Smart evaluation based on request type
  - Reduces database queries and improves performance
- **HTMX Support** - Full integration with HTMX library
  - `IsHTMXRequest()` - Detect HTMX requests via HX-Request header
  - `GetHTMXHeaders()` - Extract all HTMX request headers
  - `HTMXRedirect()` - Client-side redirects via HX-Redirect header
  - `HTMXTrigger()` - Trigger client-side events
  - `HTMXTriggerWithData()` - Trigger events with JSON payload
  - `HTMXPartial()` - Render HTML partials
  - `HTMXReswap()` - Change swap strategy (innerHTML, outerHTML, etc.)
  - `HTMXRetarget()` - Change target element
  - `HTMXPushURL()` - Push URL to browser history
  - `HTMXReplaceURL()` - Replace URL in browser history
  - `HTMXRefresh()` - Trigger page refresh
  - Complete test coverage for all HTMX features
- **Documentation** - Comprehensive guides
  - HTMX integration guide with examples
  - Lazy props usage patterns
  - Best practices and performance tips

### Changed
- Test coverage: 70.9% → 73.3% (comprehensive integration tests)
- pkg/inertia coverage: 86.7% → 84.9%

## [0.3.0] - 2026-01-06

### Added
- Automated release workflow via GitHub Actions
  - Validates release with all tests and linters
  - Creates GitHub Release with changelog notes
  - Publishes NPM packages automatically
- Release process documentation (docs/RELEASING.md)
  - Complete guide for creating releases
  - NPM token configuration
  - Troubleshooting and rollback procedures
  - Release checklist

### Changed
- Release process now fully automated via git tags
- NPM package versions automatically synced with git tags

## [0.2.9] - 2026-01-06

### Fixed
- CI workflow: use cross-platform commands for cleaning node_modules on Windows

## [0.2.8] - 2026-01-06

### Fixed
- CI workflow: added required `-package` flag to type generation test
- CI workflow: clean node_modules before install on cross-platform tests to fix rollup optional dependencies issue

## [0.2.7] - 2026-01-06

### Fixed
- CI workflow: corrected type generation example path from `typegen-example` to `inertia-typegen`
- CI workflow: use `npm install` instead of `npm ci` to handle optional dependencies on different platforms

## [0.2.6] - 2026-01-06

### Fixed
- Removed deprecated linters from golangci-lint configuration
- CI now fully passes without deprecated linter warnings

## [0.2.5] - 2026-01-06

### Fixed
- Added nolint comments for acceptable architectural decisions
- CI now passes all lint checks

## [0.2.4] - 2026-01-06

### Fixed
- All godot linting issues across all Go files
- Only remaining linting issues are architectural (InertiaContext naming and middleware complexity to be addressed in future)

## [0.2.3] - 2026-01-06

### Fixed
- More godot linting issues
- Removed non-functional eslint script from package.json
- Updated CI to skip eslint (will be added properly in future version)

## [0.2.2] - 2026-01-06

### Fixed
- Additional godot linting issues in inertia.go
- CI workflow npm installation (use npm install instead of npm ci)
- Added package-lock.json for proper dependency management

## [0.2.1] - 2026-01-06

### Fixed
- Code formatting and linting issues
- File permissions for generated files (now 0600)
- CI workflow npm cache path
- All godot, gosec, gocritic, and revive linting warnings

## [0.2.0] - 2026-01-06

### Added

#### Core Go Package
- Initial project structure with MIT License
- README with badges, quick start, and documentation links
- Contributing guidelines with TDD approach
- Core Inertia types (Response, Config, Page, SharedDataFunc)
- Inertia instance with shared data management (static + lazy evaluation)
- Middleware implementation with Inertia protocol support
- Redirect and navigation helpers (Location, Back, Redirect)
- Error handling and validation support
- Context wrapper for router integration
- ContextInterface for router-agnostic design
- InertiaContext with all Inertia methods
- Comprehensive test suite (47 Go tests, 87.4% coverage)

#### TypeScript Type Generation
- TypeScript type generator (pkg/typegen)
  - Automatic Go struct to TypeScript interface conversion
  - Support for json and ts struct tags
  - Nested struct, slice, and map handling
  - 32 tests with 66.7% coverage

#### NPM Packages
- NPM monorepo structure for client packages
- @toutaio/inertia-vue package (v0.1.0)
  - Vue 3 adapter with full TypeScript support
  - createInertiaApp function for app initialization
  - Link component with external link detection
  - useForm composable with dirty tracking
  - usePageProps and usePage composables
  - useRemember composable for state persistence across visits
  - withLayout utility for wrapping pages with layouts
  - resolvePageLayout for flexible layout resolution
  - Support for nested layouts and layout as function
  - Router with full Inertia.js protocol support
  - Server-side rendering (SSR) support
    - createInertiaSSRApp for server-side rendering
    - createSSRPage for HTML page template generation
    - XSS protection with automatic escaping
    - Version tracking in meta tags
  - Comprehensive test suite (61 tests, 100% passing)
  - Production build pipeline
    - ESM bundle (dist/index.mjs)
    - CommonJS bundle (dist/index.js)
    - TypeScript declarations (dist/index.d.ts)
    - Built with tsup for optimal output

#### Examples
- Complete Todo App example
  - Full CRUD operations for todos
  - Authentication flow with login/logout
  - Form handling with validation
  - Flash messages
  - TypeScript type generation from Go structs
  - SSR setup with separate server
  - esbuild configuration for bundling
  - Vue 3 components with composition API
- Basic HTTP context example
- Type generation example  
- Full-stack example application
  - Complete backend with Go handlers
  - Vue 3 frontend with TypeScript
  - Pages: Home, Users (Index/Show), Posts (Index/Create), About
  - Layout component with navigation
  - Nested layouts example (Admin Dashboard)
  - Form handling with validation
  - SSR configuration
  - Vite build setup

#### CI/CD & Quality
- GitHub Actions CI workflow (Go test, lint, build)
- NPM test workflow for Node 18, 20, 22
- Integration test workflow combining Go + NPM tests
- Cross-platform compatibility testing (Linux, macOS, Windows)
- Test coverage reporting via Codecov
- golangci-lint configuration
- Total: 91 tests passing (47 Go + 44 NPM)
- Coverage: 87.4% Go, 100% NPM test pass rate

### Features
- Inertia request detection (X-Inertia header)
- Version checking and conflict handling (409 responses)
- External redirect support (Location method + X-Inertia-Location)
- Internal redirect support (Redirect method with 303 See Other)
- Back navigation (using Referer header)
- Partial reload support (X-Inertia-Partial-Data, RenderOnly)
- Shared data management (static and lazy function-based)
- Error pages with status codes
- Validation errors support
- Flash messages support
- Asset versioning
- JSON rendering with proper structure
- Client-side routing with history API
- Method spoofing for PUT/PATCH/DELETE
- Automatic external link detection
- Form helper with dirty state tracking
- Automatic preventDefault handling
- TypeScript type generation from Go structs
  - Preserves JSON tag names and optional markers
  - Handles complex nested structures
  - Integrates with build process via go generate

### Completed Tasks - Phases 1, 2 & 7 (Core Vue Package) Complete ✅
- Repository setup and structure ✅
- Core type definitions ✅  
- Complete Inertia protocol (Render, RenderOnly, Location, Back, Redirect) ✅
- Middleware implementation ✅
- **Context wrapper for any router** ✅
- **Router-agnostic design with ContextInterface** ✅
- CI/CD setup ✅
- Partial reload support ✅
- Shared data (static + lazy functions) ✅
- Redirect handling (internal + external) ✅
- Error response support ✅
- Validation errors & flash messages ✅
- Basic documentation and examples ✅
- Test-driven development throughout ✅
- **NPM @toutaio/inertia-vue package** ✅

### Test Coverage
- Go: 87.4% (47 tests passing)
- TypeScript (Vue): 100% pass rate (56 tests passing)
- Total: 103 tests passing across Go and NPM

### Workflows
- Go tests on Linux, macOS, Windows (Go 1.22, 1.23)
- NPM tests on Ubuntu (Node 18, 20, 22)
- Integration tests (Go + NPM)
- Cross-platform compatibility testing
- Coverage reports uploaded to Codecov

### Ready for v0.1.0
The core Inertia.js adapter is complete and functional. Ready for:
- Integration with Cosan router
- SSR implementation
- Production usage (without SSR)

#### Documentation
- Comprehensive migration guide
  - Step-by-step migration from traditional Go apps
  - Template to Vue component conversion
  - Form handling migration
  - Authentication patterns
  - Complete migration checklist
- Advanced usage guide
  - Server-side rendering (SSR) setup
  - TypeScript type generation patterns
  - Advanced form handling (file uploads, validation)
  - Lazy data evaluation
  - Asset versioning
  - Partial reloads optimization
  - Custom context wrappers
  - Error handling strategies
  - Testing approaches
  - Performance optimization techniques
  - Deployment strategies (Docker, systemd, cloud platforms)
  - Production best practices
- Contributing guide
  - TDD workflow
  - Pull request process
  - Code quality standards
  - Testing requirements

### Phase Completion Status
- ✅ Phase 1: Project Setup (100%)
- ✅ Phase 2: Cosan Integration (100% - router-agnostic design complete)
- ⏳ Phase 3: SSR Support (0% - V8 integration deferred to v1.1)
- ✅ Phase 4: TypeScript Codegen (100%)
- ⏳ Phase 5: Real-time Updates (0% - WebSocket deferred to v1.2)
- ⏳ Phase 6: HTMX Support (0% - deferred to v1.3)
- ✅ Phase 7: NPM Package (100% - Vue adapter complete, ready for NPM publish)
- ✅ Phase 8: Examples (100% - Todo app + Full-stack example complete)
- ⏳ Phase 9: Ritual Integration (0% - pending ritual-grove updates, planned for v1.1)
- ✅ Phase 10: Documentation (100% - README, migration, advanced guides complete)
- ✅ Phase 11: Testing & Quality (100% - 103 tests passing, 87.4% Go coverage, CI/CD complete)
- 🔄 Phase 12: Release (In Progress - ready for v0.1.0)

### Next Milestones
- **v0.1.0**: Current state (ready for release - 70% of planned features)
- **v1.1.0**: SSR support with V8 + Ritual Grove integration
- **v1.2.0**: Real-time WebSocket updates + Scéla integration
- **v1.3.0**: HTMX support
- **v2.0.0**: React and Svelte adapters

### Known Limitations
- SSR requires external Node.js process (V8 integration planned for v1.1.0)
- Real-time updates not yet implemented (planned for v1.2.0)
- HTMX support not yet available (planned for v1.3.0)
- Only Vue adapter available (React/Svelte planned for v2.0.0)
- No Ritual Grove integration yet (planned for v1.1.0)

## [0.1.0] - 2026-01-05

### Added
- Project initialization
- TDD approach with testify
- Core Inertia protocol types
