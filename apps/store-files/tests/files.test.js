// Tests of the helpers in ui/app.js; the Go test package_test.go runs them too.
"use strict";
const test = require("node:test");
const assert = require("node:assert/strict");
const path = require("node:path");
const { validName, join, kind, size, decode, b64, pickLang, TEXTS } = require(path.join(__dirname, "..", "ui", "app.js"));

test("names", () => {
  for (const n of ["a", "Fotos 2026", "gruß.txt", "ä-ö_ü.md", ".bashrc"]) assert.ok(validName(n), n);
  for (const n of ["", ".", "..", "a/b", "a\b", " a", "a ", "a\nb", "\u0000", ".upload-x", "x".repeat(256), "ä".repeat(128), null])
    assert.ok(!validName(n), String(n));
});

test("helpers", () => {
  assert.equal(join("", "a"), "a");
  assert.equal(join("a/b", "c"), "a/b/c");
  assert.equal(kind("Bild.JPG"), "image");
  assert.equal(kind("notiz.txt"), "text");
  assert.equal(kind("archiv.zip"), "");
  assert.equal(kind(".txt"), "");
  assert.equal(size(10), "10 B");
  assert.equal(size(1536), "1.5 KB");
  assert.deepEqual(decode(new TextEncoder().encode('{"entries":[]}')), { entries: [] });
  assert.throws(() => decode("[]"));
  assert.throws(() => decode("null"));
  assert.equal(b64(new Uint8Array([0, 255, 104, 105])), Buffer.from([0, 255, 104, 105]).toString("base64"));
  const big = new Uint8Array(200000).map((_, i) => i * 7);
  assert.equal(b64(big), Buffer.from(big).toString("base64"));
});

test("languages", () => {
  assert.equal(pickLang(["it-IT"]), "it");
  assert.equal(pickLang(["xx"]), "en");
  const keys = Object.keys(TEXTS.de).sort();
  assert.deepEqual(Object.keys(TEXTS).sort(), ["de", "en", "es", "fr", "it"]);
  for (const l of Object.keys(TEXTS)) assert.deepEqual(Object.keys(TEXTS[l]).sort(), keys, l);
  // Every error code of the service has a sentence (classes/files.go).
  for (const c of ["call", "path", "notfound", "exists", "notempty", "toolarge", "failed", "offline"]) assert.ok(TEXTS.de[c], c);
});

// Fuzz: CD_FUZZ runs, CD_FUZZ_SEED to replay. A name that passes never
// carries a separator or a control character, and never climbs.
const RUNS = Number(process.env.CD_FUZZ || 3000);
const SEED = Number(process.env.CD_FUZZ_SEED || Date.now() % 2 ** 31);
test(`fuzz: names (seed ${SEED})`, () => {
  let x = SEED || 1;
  const r = (n) => { x ^= x << 13; x >>>= 0; x ^= x >>> 17; x ^= x << 5; x >>>= 0; return x % n; };
  const ALPHA = "ab. /\\u0000\n\u007f\u0085äß😀�-_";
  for (let i = 0; i < RUNS; i++) {
    let s = "";
    for (let k = r(12); k > 0; k--) s += [...ALPHA][r([...ALPHA].length)];
    if (!validName(s)) continue;
    assert.ok(!/[/\\\u0000-\u001f\u007f-\u009f]/.test(s) && s !== "." && s !== ".." && s === s.trim(), JSON.stringify(s));
  }
});
