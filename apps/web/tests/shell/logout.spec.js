import { test, expect, signIn, skipWithoutLogin } from "../fixtures/page.js";

// Signing out. Until this case there was no way: the start menu offered
// "Fenster anordnen" and "Alle schließen" and nothing that ended the session,
// and net/http.js logout() had one caller - the trust alarm in
// shell/trustlist.js - so a person on a shared machine stayed signed in until
// the cookie expired.
//
// The skip asks for the capability, never for an engine's name:
// skipWithoutLogin is "this engine finishes a password login", which is what
// standing on the desktop needs (fixtures/page.js).
test.beforeEach(async ({ shell }) => {
  await skipWithoutLogin(shell);
  await signIn(shell);
  await expect(shell.locator("#taskbar")).toBeVisible();
});

// There IS a way out, where a person looks for it, named out of the language
// file of this browser (de-DE in this suite), and the keyboard reaches it.
test("the start menu offers a way out, named in the language of the browser", async ({ shell }) => {
  await shell.getByRole("button", { name: "Start" }).click();
  await expect(shell.getByRole("button", { name: "Abmelden" })).toBeVisible();
  // The menu's own trapTab() walks these; a button outside them would be one
  // the keyboard cannot reach.
  const stops = await shell.evaluate(() => [...document.getElementById("startmenu")
    .querySelectorAll('input, button, [role="menuitem"][tabindex="0"]')].map((e) => e.textContent));
  expect(stops, "the button is a tab stop of the menu").toContain("Abmelden");
});

// What the button does: the API is asked once, and the login mask is back with
// nothing of the desktop left - taskbar, windows and tiles all live in #root,
// which the login replaces. A window is open first, so the teardown has
// something to take.
test("signing out asks the API and lands back at the login mask", async ({ shell }) => {
  let asked = 0;
  await shell.route("**/api/v1/logout", async (r) => { asked++; await r.fallback(); });
  await shell.getByRole("button", { name: "Start" }).click();
  await shell.getByRole("menuitem", { name: "Rechner", exact: true }).click();
  await expect(shell.locator(".win")).toHaveCount(1);

  await shell.getByRole("button", { name: "Start" }).click();
  await shell.getByRole("button", { name: "Abmelden" }).click();

  await expect(shell.getByRole("heading", { name: "Bei CirroDesk anmelden" })).toBeVisible();
  expect(asked, "the API was told").toBe(1);
  await expect(shell.locator("#taskbar"), "no taskbar").toHaveCount(0);
  await expect(shell.locator(".win"), "no window").toHaveCount(0);
  await expect(shell.locator(".smitem"), "no tile").toHaveCount(0);
});

// POINT 50: a screen that says "signed out" to somebody who is not is worse
// than one that says the attempt failed. The mock's fault door makes the next
// request answer E_BUSY, which is what a store under load answers.
test("a logout the server refuses leaves the desktop standing and says so", async ({ shell }) => {
  await shell.request.post("/api/v1/mock?fault=busy");
  await shell.getByRole("button", { name: "Start" }).click();
  await shell.getByRole("button", { name: "Abmelden" }).click();

  await expect(shell.locator("#taskbar")).toBeVisible();
  await expect(shell.getByRole("heading", { name: "Bei CirroDesk anmelden" })).toHaveCount(0);
  // The load-bearing half of the sentence, not its wording: which side of
  // point 50 the person is on. And as an error, which waits for the person.
  const status = shell.getByRole("status");
  await expect(status).toHaveText(/noch angemeldet/);
  await expect(status).toHaveClass(/error/);
  // A second press after the fault is spent gets through: no latch was left.
  await shell.getByRole("button", { name: "Start" }).click();
  await shell.getByRole("button", { name: "Abmelden" }).click();
  await expect(shell.getByRole("heading", { name: "Bei CirroDesk anmelden" })).toBeVisible();
});
