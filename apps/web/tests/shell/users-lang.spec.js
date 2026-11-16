import { test, expect, signIn, expectDesktop, skipWithoutLogin } from "../fixtures/page.js";

// S12, Benutzer und Rechte: the tabs and the users table are the browser's
// language. The dock button is found by the app's id, not its name, because
// the name is the server's own and does not change with the browser.
for (const [locale, tabs, cols] of [
  ["de-DE", ["Benutzer", "Rollen und Rechte", "Server", "Protokoll"],
    ["Name", "Rollen", "Anmeldung in zwei Schritten", "Gesperrt", "Zuletzt angemeldet"]],
  ["en-US", ["Users", "Roles and rights", "Servers", "Log"],
    ["Name", "Roles", "Two-step sign-in", "Locked", "Last signed in"]],
  ["es-ES", ["Usuarios", "Roles y derechos", "Servidores", "Registro"],
    ["Nombre", "Roles", "Inicio de sesión en dos pasos", "Bloqueado", "Último inicio de sesión"]],
  ["fr-FR", ["Utilisateurs", "Rôles et droits", "Serveurs", "Journal"],
    ["Nom", "Rôles", "Connexion en deux étapes", "Bloqué", "Dernière connexion"]],
  ["it-IT", ["Utenti", "Ruoli e diritti", "Server", "Registro"],
    ["Nome", "Ruoli", "Accesso in due passaggi", "Bloccato", "Ultimo accesso"]],
]) {
  test.describe(`users in ${locale}`, () => {
    test.use({ locale });
    test("the users app's tabs and table are in this language", async ({ page }) => {
      await page.goto("/");
      await skipWithoutLogin(page);
      await signIn(page);
      await expectDesktop(page);
      await page.locator('#dock .dapp[data-app="cd.users"]').click();
      const app = page.locator(".usersapp");
      await expect(app).toBeVisible();
      await expect(app.getByRole("tab")).toHaveText(tabs);
      await expect(app.locator(".panel.users th")).toHaveText(cols);
    });
  });
}
