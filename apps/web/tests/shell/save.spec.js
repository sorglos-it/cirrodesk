import { test, expect, signIn, expectDesktop, skipWithoutLogin } from "../fixtures/page.js";

// cd.save: an app hands a file to the shell, the shell checks it and saves it
// as a download. Driven from inside the terminal's frame, because that is
// where an app's code runs; the terminal itself plays no part.

test.beforeEach(async ({ shell }) => {
  await skipWithoutLogin(shell);
  await signIn(shell);
  await expectDesktop(shell);
});

async function frame(shell) {
  await shell.getByRole("button", { name: "Start" }).click();
  await shell.getByRole("menuitem", { name: "Terminal", exact: true }).click();
  const f = shell.frameLocator(".appframe");
  await expect(f.getByRole("textbox", { name: "Eingabe" })).toBeVisible();
  return f.locator("html");
}

// save runs cd.save in the app and hands back what the promise said.
const save = (html, name, data, type) => html.evaluate(
  (_, [n, d, t]) => window.cd.save(n, d, t).then((v) => ({ value: v }), (e) => ({ error: e.message })),
  [name, data, type]);

test("a file from the app is saved as a download, byte for byte", async ({ shell }) => {
  const html = await frame(shell);
  const download = shell.waitForEvent("download");
  const got = await html.evaluate(() =>
    window.cd.save("notiz.txt", new TextEncoder().encode("Hallo ÄÖÜ"), "text/plain").then((v) => ({ value: v })));
  expect(got).toEqual({ value: true });
  const d = await download;
  expect(d.suggestedFilename()).toBe("notiz.txt");
  const chunks = [];
  for await (const c of await d.createReadStream()) chunks.push(c);
  expect(Buffer.concat(chunks).toString("utf8")).toBe("Hallo ÄÖÜ");
});

test("a name with a path or a control character is refused, not mended", async ({ shell }) => {
  const html = await frame(shell);
  let downloads = 0;
  shell.on("download", () => downloads++);
  for (const name of ["../x", "a/b", "a\\b", "..", ".", "", " x", "x ", "a\u0000b", "a\nb",
    "a\u0085b", "x".repeat(256), "ä".repeat(128)]) {
    expect(await save(html, name, new Uint8Array([1]), "text/plain"), name).toEqual({ error: "name" });
  }
  // 255 bytes exactly is still a name.
  const d = shell.waitForEvent("download");
  expect(await save(html, "x".repeat(255), new Uint8Array([1]), "text/plain")).toEqual({ value: true });
  await d;
  expect(downloads).toBe(1);
});

test("more than 64 MiB is refused, and data that is not bytes is too", async ({ shell }) => {
  const html = await frame(shell);
  expect(await html.evaluate(() => window.cd.save("big.bin", new Uint8Array((64 << 20) + 1))
    .then(() => "saved", (e) => e.message))).toBe("toolarge");
  expect(await save(html, "x.txt", "text, not bytes", "text/plain")).toEqual({ error: "unavailable" });
});

test("the type is checked as well, and the shell's own check holds whatever the app says", async ({ shell }) => {
  const types = await shell.evaluate(async () => {
    const s = await import("/assets/js/shell/save.js");
    return ["text/plain", "application/vnd.ms-excel", "text/html; charset=utf-8", "<script>", "", null]
      .map(s.saveType);
  });
  expect(types).toEqual(["text/plain", "application/vnd.ms-excel", "application/octet-stream",
    "application/octet-stream", "application/octet-stream", "application/octet-stream"]);
});
