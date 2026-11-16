import { test, expect, signIn, expectDesktop, skipWithoutLogin } from "../fixtures/page.js";

// S12, the rest of the shell: the window's taskbar entry, the age of a screen,
// the sentence of a contract code and the setup's heading are the browser's
// language. The last three are asked through the functions that draw them,
// because reaching their screens takes a failure or a fresh installation.
for (const [locale, win, ago3, notFound, setupHead] of [
  ["de-DE", "Fenster", "vor 3 Minuten", "Nicht gefunden.", "CirroDesk einrichten"],
  ["en-US", "Window", "3 minutes ago", "Not found.", "Set up CirroDesk"],
  ["es-ES", "Ventana", "hace 3 minutos", "No encontrado.", "Configurar CirroDesk"],
  ["fr-FR", "Fenêtre", "il y a 3 minutes", "Introuvable.", "Configurer CirroDesk"],
  ["it-IT", "Finestra", "3 minuti fa", "Non trovato.", "Configura CirroDesk"],
]) {
  test.describe(`shell in ${locale}`, () => {
    test.use({ locale });
    test("the shell's own texts are in this language", async ({ page }) => {
      await page.goto("/");
      await skipWithoutLogin(page);
      await signIn(page);
      await expectDesktop(page);
      await page.locator('#dock .dapp[data-app="cd.users"]').click();
      await expect(page.locator("#tasks .task").first()).toHaveAttribute("aria-label", new RegExp(`^${win} `));
      const got = await page.evaluate(async () => {
        const { ago } = await import("/assets/js/shell/stale.js");
        const { t } = await import("/assets/js/shell/lang.js");
        const { ChannelError } = await import("/assets/js/net/channel.js");
        return { ago3: ago(3 * 60 * 1000), notFound: new ChannelError("E_NOTFOUND").message, setupHead: t("setup.ui.head") };
      });
      expect(got).toEqual({ ago3, notFound, setupHead });
    });
  });
}
