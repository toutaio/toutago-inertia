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
- GitHub Actions CI workflow (test, lint, build)
- golangci-lint configuration
- Comprehensive test suite (90.1% coverage)
- Go module setup with dependencies

### Features
- Inertia request detection (X-Inertia header)
- Version checking and conflict handling (409 responses)
- External redirect support (X-Inertia-Location)
- Partial reload support (X-Inertia-Partial-Data, RenderOnly)
- Shared data management (static and lazy function-based)
- Asset versioning
- JSON rendering with proper structure

### Completed Tasks (28/158)
- Repository setup and structure ✅
- Core type definitions ✅  
- Basic Inertia protocol (New, Render, RenderOnly) ✅
- Middleware implementation ✅
- CI/CD setup ✅
- Partial reload support ✅
- Shared data (static + lazy functions) ✅
- Test-driven development for core functionality ✅

### Test Coverage: 90.1%

### In Progress
- Cosan router integration
- SSR support with V8
- TypeScript code generation
- Real-time WebSocket updates
- HTMX alternative support

## [0.1.0] - 2026-01-05

### Added
- Project initialization
- TDD approach with testify
- Core Inertia protocol types
