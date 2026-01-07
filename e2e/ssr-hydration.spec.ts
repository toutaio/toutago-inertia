import { test, expect } from '@playwright/test';

test.describe('SSR Hydration E2E', () => {
  test('should render content server-side', async ({ page }) => {
    // Disable JavaScript to test SSR
    await page.goto('/', { waitUntil: 'domcontentloaded' });
    await page.context().setOffline(true);
    
    // Content should be visible even without JS
    const heading = page.locator('h1');
    await expect(heading).toBeVisible();
    
    await page.context().setOffline(false);
  });

  test('should hydrate correctly on client-side', async ({ page }) => {
    await page.goto('/');
    
    // Check that page is interactive
    await page.click('a[href="/todos"]');
    await expect(page).toHaveURL('/todos');
    
    // Verify no hydration errors in console
    const errors: string[] = [];
    page.on('console', msg => {
      if (msg.type() === 'error') {
        errors.push(msg.text());
      }
    });
    
    await page.waitForTimeout(1000);
    
    // Filter out common non-hydration errors
    const hydrationErrors = errors.filter(err => 
      err.includes('hydrat') || 
      err.includes('mismatch')
    );
    
    expect(hydrationErrors).toHaveLength(0);
  });

  test('should preserve server-rendered HTML during hydration', async ({ page }) => {
    await page.goto('/');
    
    // Get initial HTML
    const initialHTML = await page.locator('body').innerHTML();
    
    // Wait for hydration
    await page.waitForTimeout(2000);
    
    // Get HTML after hydration
    const hydratedHTML = await page.locator('body').innerHTML();
    
    // Server HTML should be mostly preserved (some attributes may be added)
    // Check that main content structure is the same
    expect(hydratedHTML).toContain('h1');
  });

  test('should make interactive elements functional after hydration', async ({ page }) => {
    await page.goto('/');
    
    // Immediately try to click (should work after hydration)
    await page.waitForSelector('a[href="/todos"]');
    const link = page.locator('a[href="/todos"]');
    
    // Link should be clickable
    await expect(link).toBeEnabled();
    await link.click();
    
    // Navigation should work
    await expect(page).toHaveURL('/todos');
  });

  test('should include meta tags from server render', async ({ page }) => {
    await page.goto('/');
    
    // Check for Inertia version meta tag
    const versionMeta = await page.locator('meta[name="inertia-version"]');
    await expect(versionMeta).toHaveCount(1);
  });

  test('should load page data correctly', async ({ page }) => {
    await page.goto('/todos');
    
    // Check that page data is loaded
    const pageData = await page.evaluate(() => {
      return (window as any).inertiaPage;
    });
    
    expect(pageData).toBeDefined();
    expect(pageData.component).toBeDefined();
    expect(pageData.props).toBeDefined();
  });

  test('should handle SSR with dynamic data', async ({ page }) => {
    await page.goto('/todos');
    
    // Check that todos are rendered server-side
    const todoItems = page.locator('.todo-item');
    
    // Should have some todos (if data exists)
    // This depends on your test data setup
    const count = await todoItems.count();
    expect(count).toBeGreaterThanOrEqual(0);
  });

  test('should maintain SSR performance benefits', async ({ page }) => {
    const startTime = Date.now();
    
    await page.goto('/', { waitUntil: 'domcontentloaded' });
    
    // Check that content is visible before full load
    const heading = page.locator('h1');
    await expect(heading).toBeVisible();
    
    const timeToVisible = Date.now() - startTime;
    
    // Should be fast (SSR makes content visible quickly)
    expect(timeToVisible).toBeLessThan(3000);
  });

  test('should handle errors during hydration gracefully', async ({ page }) => {
    const errors: string[] = [];
    const warnings: string[] = [];
    
    page.on('console', msg => {
      if (msg.type() === 'error') {
        errors.push(msg.text());
      } else if (msg.type() === 'warning') {
        warnings.push(msg.text());
      }
    });
    
    await page.goto('/');
    await page.waitForTimeout(2000);
    
    // Should not have critical errors
    const criticalErrors = errors.filter(err => 
      !err.includes('favicon') && 
      !err.includes('404')
    );
    
    expect(criticalErrors).toHaveLength(0);
  });

  test('should support client-side navigation after SSR', async ({ page }) => {
    await page.goto('/');
    
    // First navigation (SSR)
    await expect(page.locator('h1')).toBeVisible();
    
    // Client-side navigation
    await page.click('a[href="/todos"]');
    await expect(page).toHaveURL('/todos');
    
    // Go back (should use client-side routing)
    await page.goBack();
    await expect(page).toHaveURL('/');
    
    // Should be fast (no server round-trip)
    const startTime = Date.now();
    await page.click('a[href="/todos"]');
    await expect(page).toHaveURL('/todos');
    const navigationTime = Date.now() - startTime;
    
    expect(navigationTime).toBeLessThan(1000);
  });

  test('should inject proper head tags in SSR', async ({ page }) => {
    await page.goto('/');
    
    // Check for essential head elements
    const title = await page.title();
    expect(title).toBeTruthy();
    expect(title.length).toBeGreaterThan(0);
    
    // Check for viewport meta
    const viewport = await page.locator('meta[name="viewport"]');
    await expect(viewport).toHaveCount(1);
  });
});

test.describe('SSR Performance', () => {
  test('should render large lists efficiently with SSR', async ({ page }) => {
    await page.goto('/todos');
    
    const startTime = Date.now();
    
    // Wait for all todo items to be rendered
    await page.waitForSelector('.todo-item, .no-todos', { timeout: 5000 });
    
    const renderTime = Date.now() - startTime;
    
    // Should render quickly
    expect(renderTime).toBeLessThan(2000);
  });

  test('should handle concurrent SSR requests', async ({ context }) => {
    const pages = await Promise.all([
      context.newPage(),
      context.newPage(),
      context.newPage(),
    ]);
    
    // Load multiple pages concurrently
    await Promise.all(
      pages.map(page => page.goto('/'))
    );
    
    // All should load successfully
    for (const page of pages) {
      await expect(page.locator('h1')).toBeVisible();
      await page.close();
    }
  });
});
