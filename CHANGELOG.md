# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
  - Router with full Inertia.js protocol support
  - Server-side rendering (SSR) support
    - createInertiaSSRApp for server-side rendering
    - createSSRPage for HTML page template generation
    - XSS protection with automatic escaping
    - Version tracking in meta tags
  - Comprehensive test suite (56 tests, 100% passing)
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
- Go: 87% (77 tests passing)
- TypeScript (Vue): 100% (44 tests passing)
- Total: 121 tests passing

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

### In Progress (Next Phases)
- Cosan router integration (Phase 2)
- SSR support with V8 (Phase 3)
- ~~TypeScript code generation (Phase 4)~~ ✅ **Complete**
- Real-time WebSocket updates (Phase 5)
- Build tools integration (Phase 8)

## [0.1.0] - 2026-01-05

### Added
- Project initialization
- TDD approach with testify
- Core Inertia protocol types
