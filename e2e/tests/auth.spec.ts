import test, { expect } from "@playwright/test";

test("successful account creation should redirect to /login", async ({ page, browserName }) => {
    await page.goto("/register");

    const email = `jeff-${browserName}@gmail.com`
    const name = `Jeff ${browserName}`


    await page.fill('input[name="email"]', email);
    await page.fill('input[name="name"]', name);
    await page.fill('input[name="password"]', "SecurePassword123");
    await page.fill('input[name="password-confirm"]', "SecurePassword123");

    await page.click('button[type="submit"]');

    await expect(page).toHaveURL("/login");
});
