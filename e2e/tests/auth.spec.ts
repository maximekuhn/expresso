import test, { expect } from "@playwright/test";
import { createAccount } from "./auth.utils";


test.beforeEach(async ({ page }) => {
    await page.goto("/register");
});

test("successful account creation should redirect to /login", async ({ page, browserName }) => {
    const email = `jeff-${browserName}@gmail.com`
    const name = `Jeff ${browserName}`

    await page.fill('input[name="email"]', email);
    await page.fill('input[name="name"]', name);
    await page.fill('input[name="password"]', "SecurePassword123");
    await page.fill('input[name="password-confirm"]', "SecurePassword123");

    await page.click('button[type="submit"]');

    await expect(page).toHaveURL("/login");
});

test("should show an error message when passwords don't match", async ({ page, browserName }) => {
    const email = `albi-${browserName}@gmail.com`
    const name = `Albi ${browserName}`

    await page.fill('input[name="email"]', email);
    await page.fill('input[name="name"]', name);
    await page.fill('input[name="password"]', "SecurePassword1234");
    await page.fill('input[name="password-confirm"]', "AnotherSecurePassword123");

    await page.click('button[type="submit"]');

    const errorMessage = page.locator("#register-form-error-box");
    await expect(errorMessage).toBeVisible();
    await expect(errorMessage).toContainText("Password and confirmation must match");
});

test("should be able to re-submit when account creation failed", async ({ page, browserName }) => {
    const email = `fode-${browserName}@gmail.com`
    const name = `Fode ${browserName}`

    await page.fill('input[name="email"]', email);
    await page.fill('input[name="name"]', name);
    await page.fill('input[name="password"]', "SecurePassword1234");
    await page.fill('input[name="password-confirm"]', "AnotherSecurePassword123");

    await page.click('button[type="submit"]');

    const errorMessage = page.locator("#register-form-error-box");
    await expect(errorMessage).toBeVisible();
    await expect(errorMessage).toContainText("Password and confirmation must match");


    await page.fill('input[name="email"]', email);
    await page.fill('input[name="name"]', name);
    await page.fill('input[name="password"]', "SecurePassword1234");
    await page.fill('input[name="password-confirm"]', "SecurePassword1234");
    await page.click('button[type="submit"]');

    await page.click('button[type="submit"]');
    await expect(page).toHaveURL("/login");
});

test("[flaky] successful login should redirect to /", async ({ page, browserName }) => {
    const email = `bill-${browserName}@gmail.com`;
    const name = `Bill ${browserName}`;
    const password = "SecurePassword1234@!!";
    await createAccount(page, email, name, password);

    await page.fill('input[name="email"]', email);
    await page.fill('input[name="password"]', password);
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL("/");
});
