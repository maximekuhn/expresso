import { expect, Page } from "@playwright/test";

export async function createAccount(page: Page, email: string, name: string, password: string) {
    await page.fill('input[name="email"]', email);
    await page.fill('input[name="name"]', name);
    await page.fill('input[name="password"]', password);
    await page.fill('input[name="password-confirm"]', password);

    await page.click('button[type="submit"]');

    await expect(page).toHaveURL("/login")
}
