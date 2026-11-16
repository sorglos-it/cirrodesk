import { test, expect, signIn, expectDesktop, skipWithoutLogin } from "../fixtures/page.js";

// S12, App-Store: the tabs are the browser's language. The dock button is
// found by its place, not its name, because the app's name is the store's own.
for (const [locale, tabs] of [
  ["de-DE", ["Quellen", "Apps", "Installiert", "Hochladen", "Updates"]],
  ["en-US", ["Sources", "Apps", "Installed", "Upload", "Updates"]],
  ["es-ES", ["Fuentes", "Apps", "Instaladas", "Subir", "Actualizaciones"]],
  ["fr-FR", ["Sources", "Apps", "Installées", "Téléverser", "Mises à jour"]],
  ["it-IT", ["Fonti", "App", "Installate", "Carica", "Aggiornamenti"]],
]) {
  test.describe(`store in ${locale}`, () => {
    test.use({ locale });
    test("the store's tabs are in this language", async ({ page }) => {
      await page.goto("/");
      await skipWithoutLogin(page);
      await signIn(page);
      await expectDesktop(page);
      await page.locator("#dock .dapp").first().click();
      await expect(page.locator(".store")).toBeVisible();
      await expect(page.getByRole("tab")).toHaveText(tabs);
    });
  });
}

// S12, App-Store: a source's "checked" time is written in the shell's
// language too, not in a German nailed into store.js.
test.describe("store dates in en-US", () => {
  test.use({ locale: "en-US" });
  test("the source's check time has the American order", async ({ page }) => {
    await page.goto("/");
    await skipWithoutLogin(page);
    await signIn(page);
    await expectDesktop(page);
    await page.locator("#dock .dapp").first().click();
    await page.getByRole("tab", { name: "Sources", exact: true }).click();
    await expect(page.locator(".store .rowside").filter({ hasText: "2026" }).first())
      .toHaveText(/\d\d\/\d\d\/2026, \d\d:\d\d [AP]M/);
  });
});

// S12, App-Store: a price is written the way this language writes prices -
// in en-US the sign before the amount and a point before the cents.
test.describe("store prices in en-US", () => {
  test.use({ locale: "en-US" });
  test("a paid app's price has the American form", async ({ page }) => {
    await page.goto("/");
    await skipWithoutLogin(page);
    await signIn(page);
    await expectDesktop(page);
    await page.locator("#dock .dapp").first().click();
    await page.getByRole("tab", { name: "Apps", exact: true }).click();
    await expect(page.locator(".approw").filter({ hasText: "Tresor" }))
      .toContainText(/[€$£]\d+\.\d\d/);
  });
});
