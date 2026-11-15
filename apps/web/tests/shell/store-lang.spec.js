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
