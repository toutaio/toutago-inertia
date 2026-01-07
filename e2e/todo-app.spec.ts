import { test, expect } from '@playwright/test';

test.describe('Todo App E2E', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should display homepage', async ({ page }) => {
    await expect(page).toHaveTitle(/Toutā Inertia/);
    const heading = page.locator('h1');
    await expect(heading).toBeVisible();
    await expect(heading).toContainText('Welcome');
  });

  test('should navigate to todos page', async ({ page }) => {
    await page.click('a[href="/todos"]');
    await expect(page).toHaveURL('/todos');
    await expect(page.locator('h1')).toContainText('Todos');
  });

  test('should create a new todo', async ({ page }) => {
    await page.click('a[href="/todos"]');
    await page.click('a[href="/todos/create"]');
    
    await page.fill('input[name="title"]', 'Test Todo');
    await page.fill('textarea[name="description"]', 'This is a test todo item');
    await page.click('button[type="submit"]');
    
    await expect(page).toHaveURL('/todos');
    await expect(page.locator('text=Test Todo')).toBeVisible();
  });

  test('should edit a todo', async ({ page }) => {
    await page.click('a[href="/todos"]');
    
    // Assuming there's at least one todo, click edit
    const editButton = page.locator('a:has-text("Edit")').first();
    await editButton.click();
    
    await page.fill('input[name="title"]', 'Updated Todo');
    await page.click('button[type="submit"]');
    
    await expect(page).toHaveURL('/todos');
    await expect(page.locator('text=Updated Todo')).toBeVisible();
  });

  test('should delete a todo', async ({ page }) => {
    await page.click('a[href="/todos"]');
    
    // Count todos before delete
    const todosBeforeCount = await page.locator('.todo-item').count();
    
    // Delete first todo
    const deleteButton = page.locator('button:has-text("Delete")').first();
    await deleteButton.click();
    
    // Confirm deletion if there's a confirm dialog
    page.on('dialog', dialog => dialog.accept());
    
    // Count todos after delete
    const todosAfterCount = await page.locator('.todo-item').count();
    expect(todosAfterCount).toBe(todosBeforeCount - 1);
  });

  test('should handle validation errors', async ({ page }) => {
    await page.click('a[href="/todos"]');
    await page.click('a[href="/todos/create"]');
    
    // Submit without filling required fields
    await page.click('button[type="submit"]');
    
    // Should show validation error
    await expect(page.locator('.error')).toBeVisible();
  });

  test('should navigate using browser back button', async ({ page }) => {
    await page.click('a[href="/todos"]');
    await expect(page).toHaveURL('/todos');
    
    await page.goBack();
    await expect(page).toHaveURL('/');
  });

  test('should preserve scroll position on back navigation', async ({ page }) => {
    await page.click('a[href="/todos"]');
    
    // Scroll down
    await page.evaluate(() => window.scrollTo(0, 500));
    const scrollPosition = await page.evaluate(() => window.scrollY);
    
    // Navigate away
    const editButton = page.locator('a:has-text("Edit")').first();
    await editButton.click();
    
    // Go back
    await page.goBack();
    
    // Check scroll position restored (with some tolerance)
    const newScrollPosition = await page.evaluate(() => window.scrollY);
    expect(Math.abs(newScrollPosition - scrollPosition)).toBeLessThan(50);
  });
});
