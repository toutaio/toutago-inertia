# E2E Tests

Comprehensive end-to-end tests for toutago-inertia using Playwright.

## Overview

The E2E test suite validates the entire application stack including:

- **Todo App**: Full CRUD operations, navigation, and form handling
- **Chat App**: Real-time WebSocket updates with multi-tab synchronization
- **Navigation**: SPA routing, browser history, and state preservation
- **SSR**: Server-side rendering and client-side hydration

## Test Files

### `todo-app.spec.ts`
Tests for the todo application example:
- Homepage display and navigation
- Todo creation, editing, and deletion
- Form validation and error handling
- Browser back/forward navigation
- Scroll position preservation

### `chat.spec.ts`
Tests for real-time chat functionality:
- Chat interface rendering
- Message sending and receiving
- Multi-tab real-time synchronization
- WebSocket connection management
- Auto-reconnection on connection loss
- Typing indicators and user lists

### `navigation.spec.ts`
Tests for SPA navigation behavior:
- Navigation without page reloads
- State preservation across routes
- Browser history updates
- External navigation handling
- Partial reloads
- Loading indicators
- Form submission flows
- Flash messages

### `ssr-hydration.spec.ts`
Tests for server-side rendering:
- Server-rendered content availability
- Client-side hydration correctness
- HTML preservation during hydration
- Interactive elements after hydration
- Meta tags injection
- Page data loading
- SSR performance benchmarks
- Error handling during hydration
- Client-side navigation after SSR

## Running Tests

### Prerequisites

Install Playwright browsers:

```bash
npx playwright install --with-deps chromium firefox
```

### Run All Tests

```bash
npm run test:e2e
```

### Run with UI

```bash
npm run test:e2e:ui
```

### Run Specific Test File

```bash
npx playwright test e2e/todo-app.spec.ts
```

### Run in Specific Browser

```bash
npx playwright test --project=chromium
npx playwright test --project=firefox
```

### Debug Mode

```bash
npx playwright test --debug
```

## Test Server

The E2E tests require a running server. The test server is configured in `playwright.config.ts` and automatically starts before tests run.

Manual server start:

```bash
npm run test:e2e:serve
```

The server runs the `examples/fullstack` application on `http://localhost:3000`.

## Writing New Tests

### Basic Structure

```typescript
import { test, expect } from '@playwright/test';

test.describe('Feature Name', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should do something', async ({ page }) => {
    // Arrange
    await page.click('button');
    
    // Assert
    await expect(page.locator('h1')).toBeVisible();
  });
});
```

### Best Practices

1. **Use descriptive test names**: `should create a new todo` instead of `test1`
2. **Wait for elements**: Use `waitForSelector` or `expect().toBeVisible()` instead of fixed timeouts
3. **Isolate tests**: Each test should be independent and not rely on state from other tests
4. **Clean up**: Reset state between tests using `beforeEach` or `afterEach`
5. **Use page object pattern**: For complex pages, extract selectors and actions into helper functions
6. **Test user flows**: Focus on end-to-end user journeys, not just individual features
7. **Handle async**: Always await async operations and use proper error handling

### Common Patterns

```typescript
// Navigation
await page.click('a[href="/todos"]');
await expect(page).toHaveURL('/todos');

// Form submission
await page.fill('input[name="title"]', 'Test');
await page.click('button[type="submit"]');

// Multi-tab testing
const page2 = await context.newPage();
await page2.goto('/chat');
// ... test cross-tab synchronization
await page2.close();

// Error checking
const errors: string[] = [];
page.on('console', msg => {
  if (msg.type() === 'error') errors.push(msg.text());
});
// ... run test
expect(errors).toHaveLength(0);
```

## CI Integration

Tests run automatically in GitHub Actions. See `.github/workflows/test.yml` for configuration.

## Test Coverage

The E2E suite includes 40+ test scenarios covering:

- ✅ Navigation flows (10 tests)
- ✅ CRUD operations (5 tests)
- ✅ Real-time updates (8 tests)
- ✅ WebSocket features (3 tests)
- ✅ SSR hydration (11 tests)
- ✅ Form validation (3 tests)
- ✅ Performance benchmarks (2 tests)

## Troubleshooting

### Tests Timeout

Increase timeout in `playwright.config.ts`:

```typescript
use: {
  actionTimeout: 10000, // 10 seconds per action
}
```

### Server Won't Start

Check if port 3000 is already in use:

```bash
lsof -i :3000
kill -9 <PID>
```

### Flaky Tests

Add explicit waits or use auto-waiting selectors:

```typescript
await page.waitForLoadState('networkidle');
await page.waitForSelector('.element');
```

### Browser Installation Issues

Install dependencies manually:

```bash
npx playwright install-deps
npx playwright install
```

## Reports

After running tests, view the HTML report:

```bash
npx playwright show-report
```

## Resources

- [Playwright Documentation](https://playwright.dev)
- [Playwright Test API](https://playwright.dev/docs/api/class-test)
- [Best Practices](https://playwright.dev/docs/best-practices)
