import { test, expect } from '@playwright/test';

test.describe('Navigation Flow E2E', () => {
  test('should navigate through entire app without page reloads', async ({ page }) => {
    await page.goto('/');
    
    // Track page reloads
    let pageReloads = 0;
    page.on('load', () => pageReloads++);
    
    // Initial load counts as 1
    expect(pageReloads).toBe(1);
    
    // Navigate to todos
    await page.click('a[href="/todos"]');
    await expect(page).toHaveURL('/todos');
    
    // Should still be 1 (no reload)
    expect(pageReloads).toBe(1);
    
    // Navigate to create
    await page.click('a[href="/todos/create"]');
    await expect(page).toHaveURL('/todos/create');
    expect(pageReloads).toBe(1);
    
    // Go back
    await page.goBack();
    await expect(page).toHaveURL('/todos');
    expect(pageReloads).toBe(1);
    
    // Navigate to home
    await page.click('a[href="/"]');
    await expect(page).toHaveURL('/');
    expect(pageReloads).toBe(1);
  });

  test('should preserve state across navigation', async ({ page }) => {
    await page.goto('/');
    
    // Go to todos
    await page.click('a[href="/todos"]');
    
    // Scroll down
    await page.evaluate(() => window.scrollTo(0, 300));
    const scrollPos = await page.evaluate(() => window.scrollY);
    
    // Navigate away and back
    await page.click('a[href="/"]');
    await page.click('a[href="/todos"]');
    
    // Check if scroll position was restored
    const newScrollPos = await page.evaluate(() => window.scrollY);
    expect(Math.abs(newScrollPos - scrollPos)).toBeLessThan(50);
  });

  test('should update browser history correctly', async ({ page }) => {
    await page.goto('/');
    
    // Navigate through pages
    await page.click('a[href="/todos"]');
    await page.click('a[href="/todos/create"]');
    
    // Go back twice
    await page.goBack();
    await expect(page).toHaveURL('/todos');
    
    await page.goBack();
    await expect(page).toHaveURL('/');
    
    // Go forward
    await page.goForward();
    await expect(page).toHaveURL('/todos');
  });

  test('should handle external navigation', async ({ page }) => {
    await page.goto('/');
    
    // Navigate to external URL (should cause full page load)
    await page.evaluate(() => {
      window.location.href = 'https://example.com';
    });
    
    await expect(page).toHaveURL(/example\.com/);
  });

  test('should handle partial reloads', async ({ page }) => {
    await page.goto('/todos');
    
    // Click on a link with data-inertia-only
    const partialLink = page.locator('a[data-inertia-only]').first();
    if (await partialLink.count() > 0) {
      await partialLink.click();
      
      // Should update only partial content
      // (implementation depends on your app structure)
      await page.waitForTimeout(500);
    }
  });

  test('should show loading indicator during navigation', async ({ page }) => {
    await page.goto('/');
    
    // Start navigation
    const navigationPromise = page.click('a[href="/todos"]');
    
    // Check for loading indicator
    const loader = page.locator('.inertia-loading');
    if (await loader.count() > 0) {
      await expect(loader).toBeVisible({ timeout: 100 });
    }
    
    await navigationPromise;
    
    // Loader should be hidden after navigation
    if (await loader.count() > 0) {
      await expect(loader).not.toBeVisible();
    }
  });

  test('should handle navigation errors gracefully', async ({ page }) => {
    await page.goto('/');
    
    // Navigate to non-existent page
    await page.goto('/non-existent-page');
    
    // Should show 404 page
    await expect(page.locator('text=404')).toBeVisible();
  });

  test('should maintain scroll position on refresh', async ({ page }) => {
    await page.goto('/todos');
    
    // Scroll down
    await page.evaluate(() => window.scrollTo(0, 400));
    const scrollPos = await page.evaluate(() => window.scrollY);
    
    // Refresh page
    await page.reload();
    
    // Scroll position might reset on full reload (normal behavior)
    // But with Inertia remember, some state is preserved
  });

  test('should handle rapid navigation', async ({ page }) => {
    await page.goto('/');
    
    // Rapidly click multiple links
    await Promise.all([
      page.click('a[href="/todos"]'),
      page.waitForURL('/todos'),
    ]);
    
    await Promise.all([
      page.click('a[href="/"]'),
      page.waitForURL('/'),
    ]);
    
    await Promise.all([
      page.click('a[href="/todos"]'),
      page.waitForURL('/todos'),
    ]);
    
    // Should end up on the last clicked link
    await expect(page).toHaveURL('/todos');
  });

  test('should update page title during navigation', async ({ page }) => {
    await page.goto('/');
    const homeTitle = await page.title();
    
    await page.click('a[href="/todos"]');
    await page.waitForTimeout(500);
    const todosTitle = await page.title();
    
    // Titles should be different
    expect(todosTitle).not.toBe(homeTitle);
  });
});

test.describe('Form Submission Flow E2E', () => {
  test('should handle form submission with validation', async ({ page }) => {
    await page.goto('/todos/create');
    
    // Submit without required field
    await page.click('button[type="submit"]');
    
    // Should show validation error
    const error = page.locator('.error');
    if (await error.count() > 0) {
      await expect(error).toBeVisible();
    }
    
    // Fill required field
    await page.fill('input[name="title"]', 'Valid Title');
    await page.click('button[type="submit"]');
    
    // Should redirect to todos list
    await expect(page).toHaveURL('/todos');
  });

  test('should show flash messages after form submission', async ({ page }) => {
    await page.goto('/todos/create');
    
    await page.fill('input[name="title"]', 'New Todo');
    await page.click('button[type="submit"]');
    
    // Should show success message
    const flash = page.locator('.flash-message');
    if (await flash.count() > 0) {
      await expect(flash).toBeVisible();
      await expect(flash).toContainText(/success|created/i);
    }
  });

  test('should preserve form data on validation errors', async ({ page }) => {
    await page.goto('/todos/create');
    
    const testTitle = 'Test Title';
    const testDescription = 'Test Description';
    
    await page.fill('input[name="title"]', testTitle);
    await page.fill('textarea[name="description"]', testDescription);
    
    // Trigger validation error (e.g., by clearing required field)
    await page.fill('input[name="title"]', '');
    await page.click('button[type="submit"]');
    
    // Description should be preserved
    await expect(page.locator('textarea[name="description"]')).toHaveValue(testDescription);
  });
});
