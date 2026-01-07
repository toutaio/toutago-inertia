import { test, expect } from '@playwright/test';

test.describe('Chat App E2E - Real-time Updates', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/chat');
  });

  test('should display chat interface', async ({ page }) => {
    await expect(page.locator('h1')).toContainText('Chat');
    await expect(page.locator('input[name="message"]')).toBeVisible();
    await expect(page.locator('button[type="submit"]')).toBeVisible();
  });

  test('should send a message', async ({ page }) => {
    const messageText = `Test message ${Date.now()}`;
    
    await page.fill('input[name="message"]', messageText);
    await page.click('button[type="submit"]');
    
    // Message should appear in the chat
    await expect(page.locator(`text=${messageText}`)).toBeVisible();
    
    // Input should be cleared
    await expect(page.locator('input[name="message"]')).toHaveValue('');
  });

  test('should receive real-time messages in multiple tabs', async ({ context, page }) => {
    // Open second tab
    const page2 = await context.newPage();
    await page2.goto('/chat');
    
    // Send message from first tab
    const messageText = `Multi-tab message ${Date.now()}`;
    await page.fill('input[name="message"]', messageText);
    await page.click('button[type="submit"]');
    
    // Message should appear in first tab
    await expect(page.locator(`text=${messageText}`)).toBeVisible();
    
    // Message should also appear in second tab (real-time update)
    await expect(page2.locator(`text=${messageText}`)).toBeVisible({ timeout: 5000 });
    
    await page2.close();
  });

  test('should show connection status', async ({ page }) => {
    // Should show connected status
    const status = page.locator('.connection-status');
    await expect(status).toContainText(/connected/i);
  });

  test('should reconnect on connection loss', async ({ page, context }) => {
    // Simulate offline
    await context.setOffline(true);
    
    // Wait for disconnected status
    const status = page.locator('.connection-status');
    await expect(status).toContainText(/disconnected/i, { timeout: 5000 });
    
    // Restore connection
    await context.setOffline(false);
    
    // Should reconnect
    await expect(status).toContainText(/connected/i, { timeout: 10000 });
  });

  test('should display message history', async ({ page }) => {
    // Send multiple messages
    for (let i = 1; i <= 3; i++) {
      await page.fill('input[name="message"]', `Message ${i}`);
      await page.click('button[type="submit"]');
      await page.waitForTimeout(100);
    }
    
    // All messages should be visible
    await expect(page.locator('text=Message 1')).toBeVisible();
    await expect(page.locator('text=Message 2')).toBeVisible();
    await expect(page.locator('text=Message 3')).toBeVisible();
  });

  test('should handle empty message submission', async ({ page }) => {
    const messagesCount = await page.locator('.message').count();
    
    // Try to submit empty message
    await page.click('button[type="submit"]');
    
    // No new message should be added
    const newMessagesCount = await page.locator('.message').count();
    expect(newMessagesCount).toBe(messagesCount);
  });

  test('should auto-scroll to latest message', async ({ page }) => {
    // Send enough messages to cause scroll
    for (let i = 1; i <= 20; i++) {
      await page.fill('input[name="message"]', `Scroll test ${i}`);
      await page.click('button[type="submit"]');
      await page.waitForTimeout(50);
    }
    
    // Last message should be visible (auto-scrolled)
    await expect(page.locator('text=Scroll test 20')).toBeVisible();
  });
});

test.describe('Chat App E2E - WebSocket Features', () => {
  test('should handle user typing indicators', async ({ context, page }) => {
    // Open second tab
    const page2 = await context.newPage();
    await page2.goto('/chat');
    
    // Type in first tab
    await page.fill('input[name="message"]', 'typing...');
    
    // Second tab should show typing indicator
    const typingIndicator = page2.locator('.typing-indicator');
    await expect(typingIndicator).toBeVisible({ timeout: 2000 });
    
    // Clear input in first tab
    await page.fill('input[name="message"]', '');
    
    // Typing indicator should disappear
    await expect(typingIndicator).not.toBeVisible({ timeout: 2000 });
    
    await page2.close();
  });

  test('should display user list', async ({ page }) => {
    const userList = page.locator('.user-list');
    await expect(userList).toBeVisible();
    
    // Should show at least current user
    const users = page.locator('.user-list .user');
    await expect(users).not.toHaveCount(0);
  });

  test('should update user count in real-time', async ({ context, page }) => {
    const userCountBefore = await page.locator('.user-count').textContent();
    
    // Open new tab (new user)
    const page2 = await context.newPage();
    await page2.goto('/chat');
    
    // Wait for user count to update
    await page.waitForTimeout(1000);
    const userCountAfter = await page.locator('.user-count').textContent();
    
    // Count should increase
    expect(userCountAfter).not.toBe(userCountBefore);
    
    await page2.close();
  });
});
