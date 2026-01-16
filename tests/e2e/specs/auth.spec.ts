import { test, expect } from '@playwright/test';

test.describe('Authentication', () => {
  test('should show landing page', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('Monitor Your SEO')).toBeVisible();
    await expect(page.getByText('Like a Pro')).toBeVisible();
  });

  test('should navigate to login from landing page', async ({ page }) => {
    await page.goto('/');
    await page.click('text=Sign In');
    await expect(page).toHaveURL('/login');
    await expect(page.getByText('Welcome back')).toBeVisible();
  });

  test('should navigate to register from landing page', async ({ page }) => {
    await page.goto('/');
    await page.click('text=Get Started');
    await expect(page).toHaveURL('/register');
    await expect(page.getByText('Create an account')).toBeVisible();
  });

  test('should register new user', async ({ page }) => {
    const uniqueEmail = `test${Date.now()}@example.com`;

    await page.goto('/register');
    await page.fill('#name', 'Test User');
    await page.fill('#email', uniqueEmail);
    await page.fill('#password', 'password123');
    await page.fill('#confirmPassword', 'password123');
    await page.click('button[type="submit"]');

    // Should redirect to login after successful registration
    await expect(page).toHaveURL('/login', { timeout: 10000 });
  });

  test('should show error for mismatched passwords', async ({ page }) => {
    await page.goto('/register');
    await page.fill('#name', 'Test User');
    await page.fill('#email', 'test@example.com');
    await page.fill('#password', 'password123');
    await page.fill('#confirmPassword', 'differentpassword');
    await page.click('button[type="submit"]');

    // Should show toast error about mismatched passwords
    await expect(page.getByText("Passwords don't match")).toBeVisible({ timeout: 5000 });
  });

  test('should show error for short password', async ({ page }) => {
    await page.goto('/register');
    await page.fill('#name', 'Test User');
    await page.fill('#email', 'test@example.com');
    await page.fill('#password', 'short');
    await page.fill('#confirmPassword', 'short');
    await page.click('button[type="submit"]');

    // Should show toast error about short password
    await expect(page.getByText('Password too short')).toBeVisible({ timeout: 5000 });
  });

  test('should login and redirect to dashboard', async ({ page }) => {
    // First register a test user via the form
    const uniqueEmail = `testlogin${Date.now()}@example.com`;

    await page.goto('/register');
    await page.fill('#name', 'Login Test User');
    await page.fill('#email', uniqueEmail);
    await page.fill('#password', 'password123');
    await page.fill('#confirmPassword', 'password123');
    await page.click('button[type="submit"]');

    // Wait for redirect to login
    await expect(page).toHaveURL('/login', { timeout: 10000 });

    // Now login with the registered user
    await page.fill('#email', uniqueEmail);
    await page.fill('#password', 'password123');
    await page.click('button[type="submit"]');

    // Should redirect to backlinks dashboard
    await expect(page).toHaveURL('/backlinks', { timeout: 10000 });
    await expect(page.getByText('Backlinks')).toBeVisible();
  });

  test('should show error for wrong password', async ({ page }) => {
    await page.goto('/login');
    await page.fill('#email', 'nonexistent@example.com');
    await page.fill('#password', 'wrongpassword');
    await page.click('button[type="submit"]');

    // Should show toast error
    await expect(page.getByText('Login failed')).toBeVisible({ timeout: 5000 });
  });

  test('should protect dashboard routes', async ({ page }) => {
    // Clear any cookies first
    await page.context().clearCookies();

    await page.goto('/backlinks');
    // Should redirect to login
    await expect(page).toHaveURL('/login');
  });

  test('should protect index-checker route', async ({ page }) => {
    await page.context().clearCookies();
    await page.goto('/index-checker');
    await expect(page).toHaveURL('/login');
  });

  test('should protect site-health route', async ({ page }) => {
    await page.context().clearCookies();
    await page.goto('/site-health');
    await expect(page).toHaveURL('/login');
  });

  test('should logout user', async ({ page }) => {
    // First register and login
    const uniqueEmail = `testlogout${Date.now()}@example.com`;

    await page.goto('/register');
    await page.fill('#name', 'Logout Test User');
    await page.fill('#email', uniqueEmail);
    await page.fill('#password', 'password123');
    await page.fill('#confirmPassword', 'password123');
    await page.click('button[type="submit"]');

    await expect(page).toHaveURL('/login', { timeout: 10000 });

    await page.fill('#email', uniqueEmail);
    await page.fill('#password', 'password123');
    await page.click('button[type="submit"]');

    await expect(page).toHaveURL('/backlinks', { timeout: 10000 });

    // Now test logout - click on user icon dropdown
    await page.click('button:has([data-lucide="user"])');

    // Wait for dropdown to appear and click logout
    await page.click('text=Log out');

    // Should redirect to login
    await expect(page).toHaveURL('/login', { timeout: 10000 });
  });

  test('should redirect authenticated user from login page', async ({ page }) => {
    // First register and login
    const uniqueEmail = `testauthredirect${Date.now()}@example.com`;

    await page.goto('/register');
    await page.fill('#name', 'Auth Redirect Test');
    await page.fill('#email', uniqueEmail);
    await page.fill('#password', 'password123');
    await page.fill('#confirmPassword', 'password123');
    await page.click('button[type="submit"]');

    await expect(page).toHaveURL('/login', { timeout: 10000 });

    await page.fill('#email', uniqueEmail);
    await page.fill('#password', 'password123');
    await page.click('button[type="submit"]');

    await expect(page).toHaveURL('/backlinks', { timeout: 10000 });

    // Try to navigate to login while authenticated
    await page.goto('/login');

    // Should redirect back to dashboard
    await expect(page).toHaveURL('/backlinks', { timeout: 5000 });
  });

  test('should navigate between auth pages', async ({ page }) => {
    await page.goto('/login');
    await expect(page.getByText('Welcome back')).toBeVisible();

    // Click sign up link
    await page.click('text=Sign up');
    await expect(page).toHaveURL('/register');

    // Click sign in link
    await page.click('text=Sign in');
    await expect(page).toHaveURL('/login');
  });
});
