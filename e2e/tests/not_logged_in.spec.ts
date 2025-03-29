import { test, expect } from "@playwright/test";

test("index should redirect to /login", async ({ page }) => {
    await page.goto("/");
    await expect(page).toHaveURL("/login");
});
