import { test, expect } from "../fixtures/page.js";

// S12, the login mask: its heading, fields and buttons are the browser's
// language too. Read before signing in, as the browser names them.
for (const [locale, head, user, submit] of [
  ["de-DE", "Bei CirroDesk anmelden", "Benutzername", "Anmelden"],
  ["en-US", "Sign in to CirroDesk", "User name", "Sign in"],
  ["es-ES", "Iniciar sesión en CirroDesk", "Nombre de usuario", "Iniciar sesión"],
  ["fr-FR", "Se connecter à CirroDesk", "Nom d'utilisateur", "Se connecter"],
  ["it-IT", "Accedi a CirroDesk", "Nome utente", "Accedi"],
]) {
  test.describe(`login mask in ${locale}`, () => {
    test.use({ locale });
    test("the mask is in this language", async ({ page }) => {
      await page.goto("/");
      await expect(page.getByRole("heading", { name: head, exact: true })).toBeVisible();
      await expect(page.getByLabel(user, { exact: true })).toBeVisible();
      await expect(page.getByRole("button", { name: submit, exact: true })).toBeVisible();
    });
  });
}
