// Tests of ui/app.js; the Go test package_test.go runs them too.
"use strict";
const test = require("node:test");
const assert = require("node:assert/strict");
const path = require("node:path");
const { hands, pickLang, TITLES } = require(path.join(__dirname, "..", "ui", "app.js"));

test("hands", () => {
  assert.deepEqual(hands(0, 0, 0), { h: 0, m: 0, s: 0 });
  assert.deepEqual(hands(3, 0, 0), { h: 90, m: 0, s: 0 });
  assert.deepEqual(hands(15, 30, 0), { h: 105, m: 180, s: 0 });
  assert.equal(hands(12, 0, 45).s, 270);
  // Every second of a day: all hands stay in [0, 360).
  for (let t = 0; t < 86400; t++) {
    const a = hands(Math.floor(t / 3600), Math.floor(t / 60) % 60, t % 60);
    for (const k of ["h", "m", "s"]) assert.ok(a[k] >= 0 && a[k] < 360, `${t} ${k}`);
  }
});

test("languages", () => {
  assert.equal(pickLang(["es-MX"]), "es");
  assert.equal(pickLang([]), "en");
  assert.deepEqual(Object.keys(TITLES).sort(), ["de", "en", "es", "fr", "it"]);
});
