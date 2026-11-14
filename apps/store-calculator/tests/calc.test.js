// Tests of the parser in ui/app.js; run with `node --test apps/store-calculator/tests/`.
// The Go test beside manifest.json runs them too, so `go test ./...` covers them.
"use strict";
const test = require("node:test");
const assert = require("node:assert/strict");
const path = require("node:path");
const { evaluate, format, pickLang, TEXTS, MAX_LEN, MAX_DEPTH } = require(path.join(__dirname, "..", "ui", "app.js"));

const ok = (s, v) => assert.deepEqual(evaluate(s), { ok: true, value: v }, s);
const bad = (s, e) => assert.deepEqual(evaluate(s), { ok: false, error: e }, s);

test("point before line and brackets", () => {
  ok("1+2*3", 7);
  ok("(1+2)*3", 9);
  ok("2*(3+(4-1))/3", 4);
  ok("((((7))))", 7);
  ok("10-4-3", 3);
  ok("8/4/2", 1);
  ok("-(2+3)", -5);
  ok("2*-3", -6);
  ok(" 1 , 5 ".replace(/ /g, "") , 1.5);
  ok("3×(2−1)÷2", 1.5);
  ok("  7  ", 7);
});

test("errors", () => {
  bad("", "empty");
  bad("   ", "empty");
  bad("(1+2", "paren");
  bad("1+2)", "paren");
  bad(")(", "paren");
  bad("1/0", "divzero");
  bad("1/(2-2)", "divzero");
  bad("1+", "syntax");
  bad("()", "syntax");
  bad("1 2", "syntax");
  bad("1.", "syntax");
  bad("2**3", "syntax");
  bad("abc", "syntax");
  bad("1".repeat(MAX_LEN + 1), "long");
  bad("(".repeat(MAX_DEPTH + 1) + "1" + ")".repeat(MAX_DEPTH + 1), "long");
  // "range" cannot be reached within MAX_LEN: a product never has more digits than its input.
  bad(null, "syntax");
});

test("format and languages", () => {
  assert.equal(format(0.1 + 0.2), "0.3");
  assert.equal(format(-0), "0");
  assert.equal(pickLang(["fr-CH"]), "fr");
  assert.equal(pickLang(["xx", "es"]), "es");
  assert.equal(pickLang([]), "en");
  const keys = Object.keys(TEXTS.de).sort();
  assert.deepEqual(Object.keys(TEXTS).sort(), ["de", "en", "es", "fr", "it"]);
  for (const l of Object.keys(TEXTS)) assert.deepEqual(Object.keys(TEXTS[l]).sort(), keys, l);
});

// Fuzz: NGC_FUZZ runs (default 3000), seed from NGC_FUZZ_SEED so a failure can be replayed.
const RUNS = Number(process.env.NGC_FUZZ || 3000);
const SEED = Number(process.env.NGC_FUZZ_SEED || Date.now() % 2 ** 31);
function rng(seed) {
  let x = seed || 1;
  return (n) => {
    x ^= x << 13; x >>>= 0; x ^= x >>> 17; x ^= x << 5; x >>>= 0;
    return x % n;
  };
}

// A random tree, printed with brackets where needed and evaluated by its own
// rules: a differential oracle that shares no code with the parser.
function tree(r, d) {
  if (d === 0 || r(3) === 0) {
    const v = r(1000);
    return { s: String(v), v };
  }
  if (r(6) === 0) {
    const a = tree(r, d - 1);
    return { s: "-(" + a.s + ")", v: a.v === null ? null : -a.v };
  }
  const a = tree(r, d - 1), b = tree(r, d - 1);
  const op = "+-*/"[r(4)];
  let v = null;
  if (a.v !== null && b.v !== null) {
    if (op === "+") v = a.v + b.v;
    else if (op === "-") v = a.v - b.v;
    else if (op === "*") v = a.v * b.v;
    else v = b.v === 0 ? "divzero" : a.v / b.v;
  }
  if (a.v === "divzero" || b.v === "divzero") v = "divzero";
  return { s: "(" + a.s + op + b.s + ")", v };
}

test(`fuzz: trees agree with the oracle (seed ${SEED})`, () => {
  const r = rng(SEED);
  for (let n = 0; n < RUNS; n++) {
    const t = tree(r, 1 + r(5));
    if (t.s.length > MAX_LEN) continue;
    const got = evaluate(t.s);
    if (t.v === "divzero") assert.deepEqual(got, { ok: false, error: "divzero" }, t.s);
    else if (Number.isFinite(t.v)) {
      assert.ok(got.ok, t.s);
      assert.ok(Math.abs(got.value - t.v) <= 1e-9 * Math.max(1, Math.abs(t.v)), t.s);
    }
  }
});

test(`fuzz: any input gives an answer, never a throw (seed ${SEED})`, () => {
  const r = rng(SEED ^ 0x5bd1e995);
  const ALPHA = "0123456789+-*/()., ×÷−x\u0000é(((";
  const ERRORS = new Set(["empty", "paren", "divzero", "syntax", "long", "range"]);
  for (let n = 0; n < RUNS; n++) {
    let s = "";
    for (let k = r(MAX_LEN + 20); k > 0; k--) s += ALPHA[r(ALPHA.length)];
    const got = evaluate(s);
    if (got.ok) assert.ok(Number.isFinite(got.value), s);
    else assert.ok(ERRORS.has(got.error), s);
  }
});
