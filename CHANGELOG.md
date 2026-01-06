# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project structure with MIT License
- README with badges, quick start, and documentation links
- Contributing guidelines with TDD approach
- Core Inertia types (Response, Config, Page, SharedDataFunc)
- Inertia instance with shared data management  
- Middleware implementation with Inertia protocol support
- Redirect and navigation helpers (Location, Back, Redirect)
- Error handling and validation support
- **Context wrapper for router integration**
- **ContextInterface for router-agnostic design**
- **InertiaContext with all Inertia methods**
- **NPM monorepo structure for client packages**
- **@toutaio/inertia-vue package (v0.1.0)**
  - Vue 3 adapter with full TypeScript support
  - createInertiaApp function for app initialization
  - **Link component with external link detection** ✨
  - **useForm composable with dirty tracking** ✨
  - **usePageProps and usePage composables**
  - Router with full Inertia.js protocol support
  - **Server-side rendering (SSR) support** ✨
    - createInertiaSSRApp for server-side rendering
    - createSSRPage for HTML page template generation
    - XSS protection with automatic escaping
    - Version tracking in meta tags
  - Comprehensive test suite (69 tests, 100% passing) ✨
  - Build pipeline with tsup (CJS, ESM, TypeScript declarations)
- GitHub Actions CI workflow (test, lint, build)
- golangci-lint configuration
- Comprehensive test suite (87.4% coverage Go, 69 tests passing TypeScript) ✨
- Go module setup with dependencies
- Getting started documentation
- Basic example application
- **HTTP context wrapper example**

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
- **Client-side routing with history API**
- **Method spoofing for PUT/PATCH/DELETE**
- **Automatic external link detection** ✨
- **Form helper with dirty state tracking** ✨
- **Automatic preventDefault handling** ✨

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
- Go: 87.4% (45 tests passing)
- TypeScript (Vue): 100% (58 tests passing) ✨

### Ready for v0.1.0
The core Inertia.js adapter is complete and functional. Ready for:
- Integration with Cosan router
- SSR implementation
- Production usage (without SSR)

### In Progress (Next Phases)
- Cosan router integration (Phase 2)
- SSR support with V8 (Phase 3)
- TypeScript code generation (Phase 4)
- Real-time WebSocket updates (Phase 5)
- Build tools integration (Phase 8)

## [0.1.0] - 2026-01-05

### Added
- Project initialization
- TDD approach with testify
- Core Inertia protocol types
