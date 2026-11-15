// org.sorglos.calc - a calculator that understands brackets.
//
// No eval and no Function: the frame's policy forbids them, and a calculator
// that evaluates a string runs whatever was typed. The expression is read by a
// small recursive-descent parser instead:
//
//   expr   = term   { ("+" | "-") term }
//   term   = factor { ("*" | "/") factor }
//   factor = ("+" | "-") factor | "(" expr ")" | number
//   number = digits [ ("." | ",") digits ]
//
// A classic script, not a module: the frame carries every file as a data:
// address and a module could not import anything. Under node (the tests) the
// parser is exported and no DOM is touched.
(function () {
  "use strict";

  // Limits, because the input comes from outside: the length bounds the work,
  // the depth bounds the recursion, so no input can exhaust the stack.
  const MAX_LEN = 256;
  const MAX_DEPTH = 64;

  // Two digits dropped before display: 0.1 + 0.2 is 0.30000000000000004.
  const ROUND = 1e10;

  class CalcError extends Error {
    constructor(code) {
      super(code);
      this.code = code;
    }
  }

  // evaluate returns {ok: true, value} or {ok: false, error}, error one of
  // "empty", "paren", "divzero", "syntax", "long", "range". It never throws.
  function evaluate(src) {
    if (typeof src !== "string") return { ok: false, error: "syntax" };
    if (src.length > MAX_LEN) return { ok: false, error: "long" };
    const s = src.replace(/×/g, "*").replace(/÷/g, "/").replace(/−/g, "-");
    // Brackets are counted before parsing, so an unpaired one is always
    // "paren" and never a syntax error further along.
    let open = 0;
    for (const c of s) {
      if (c === "(") open++;
      else if (c === ")" && --open < 0) return { ok: false, error: "paren" };
    }
    if (open !== 0) return { ok: false, error: "paren" };
    let i = 0;
    let depth = 0;

    const skip = () => {
      while (i < s.length && (s[i] === " " || s[i] === "\t")) i++;
    };
    const peek = () => {
      skip();
      return s[i];
    };

    function expr() {
      let v = term();
      for (let c = peek(); c === "+" || c === "-"; c = peek()) {
        i++;
        const r = term();
        v = c === "+" ? v + r : v - r;
      }
      return v;
    }

    function term() {
      let v = factor();
      for (let c = peek(); c === "*" || c === "/"; c = peek()) {
        i++;
        const r = factor();
        if (c === "/") {
          if (r === 0) throw new CalcError("divzero");
          v /= r;
        } else v *= r;
      }
      return v;
    }

    function factor() {
      if (++depth > MAX_DEPTH) throw new CalcError("long");
      try {
        const c = peek();
        if (c === "+" || c === "-") {
          i++;
          const v = factor();
          return c === "-" ? -v : v;
        }
        if (c === "(") {
          i++;
          const v = expr();
          if (peek() !== ")") throw new CalcError("syntax");
          i++;
          return v;
        }
        return number();
      } finally {
        depth--;
      }
    }

    function number() {
      const start = i;
      while (i < s.length && s[i] >= "0" && s[i] <= "9") i++;
      if (i === start) throw new CalcError("syntax");
      if (s[i] === "." || s[i] === ",") {
        i++;
        const frac = i;
        while (i < s.length && s[i] >= "0" && s[i] <= "9") i++;
        if (i === frac) throw new CalcError("syntax");
      }
      return Number(s.slice(start, i).replace(",", "."));
    }

    if (peek() === undefined) return { ok: false, error: "empty" };
    try {
      const v = expr();
      if (peek() !== undefined) return { ok: false, error: "syntax" };
      if (!Number.isFinite(v)) return { ok: false, error: "range" };
      return { ok: true, value: v };
    } catch (e) {
      if (e instanceof CalcError) return { ok: false, error: e.code };
      throw e;
    }
  }

  // format rounds for display; -0 shows as 0.
  function format(v) {
    const r = Math.round(v * ROUND) / ROUND;
    return String(Object.is(r, -0) ? 0 : Number.isFinite(r) ? r : v);
  }

  // The words of the page in the five languages of the shell (shell/lang.js).
  const TEXTS = {
    de: { title: "Taschenrechner", input: "Rechnung", result: "Ergebnis", clear: "Löschen", back: "Zurück",
      empty: "Bitte eine Rechnung eingeben.", paren: "Die Klammern passen nicht zusammen.",
      divzero: "Durch null kann man nicht teilen.", syntax: "Diese Rechnung verstehe ich nicht.",
      long: "Die Rechnung ist zu lang.", range: "Das Ergebnis ist zu groß." },
    en: { title: "Calculator", input: "Calculation", result: "Result", clear: "Clear", back: "Back",
      empty: "Please enter a calculation.", paren: "The brackets do not match.",
      divzero: "You cannot divide by zero.", syntax: "This calculation cannot be read.",
      long: "The calculation is too long.", range: "The result is too large." },
    fr: { title: "Calculatrice", input: "Calcul", result: "Résultat", clear: "Effacer", back: "Retour",
      empty: "Veuillez saisir un calcul.", paren: "Les parenthèses ne correspondent pas.",
      divzero: "Impossible de diviser par zéro.", syntax: "Ce calcul est illisible.",
      long: "Le calcul est trop long.", range: "Le résultat est trop grand." },
    it: { title: "Calcolatrice", input: "Calcolo", result: "Risultato", clear: "Cancella", back: "Indietro",
      empty: "Inserisci un calcolo.", paren: "Le parentesi non corrispondono.",
      divzero: "Non si può dividere per zero.", syntax: "Questo calcolo non è leggibile.",
      long: "Il calcolo è troppo lungo.", range: "Il risultato è troppo grande." },
    es: { title: "Calculadora", input: "Cálculo", result: "Resultado", clear: "Borrar", back: "Atrás",
      empty: "Introduce un cálculo.", paren: "Los paréntesis no coinciden.",
      divzero: "No se puede dividir entre cero.", syntax: "Este cálculo no se puede leer.",
      long: "El cálculo es demasiado largo.", range: "El resultado es demasiado grande." },
  };

  // The shell's fallback is English (shell/lang.js FALLBACK).
  function pickLang(tags) {
    for (const t of tags) {
      const l = String(t || "").slice(0, 2).toLowerCase();
      if (Object.prototype.hasOwnProperty.call(TEXTS, l)) return l;
    }
    return "en";
  }

  if (typeof document === "undefined") {
    module.exports = { evaluate, format, pickLang, TEXTS, MAX_LEN, MAX_DEPTH };
    return;
  }

  const lang = pickLang([window.cd && cd.lang, document.documentElement.lang]);
  const T = TEXTS[lang];
  document.documentElement.lang = lang;
  document.title = T.title;

  function el(tag, cls, txt) {
    const e = document.createElement(tag);
    if (cls) e.className = cls;
    if (txt !== undefined) e.textContent = String(txt);
    return e;
  }

  const box = el("div", "calc");
  const input = el("input");
  input.type = "text";
  input.inputMode = "decimal";
  input.autocomplete = "off";
  input.spellcheck = false;
  input.maxLength = MAX_LEN;
  input.setAttribute("aria-label", T.input);
  const out = el("output", "", "0");
  out.setAttribute("aria-label", T.result);
  out.setAttribute("aria-live", "polite");
  const keys = el("div", "keys");
  box.append(input, out, keys);

  const equals = () => {
    const r = evaluate(input.value);
    out.className = r.ok ? "" : "err";
    out.textContent = r.ok ? format(r.value) : T[r.error];
  };
  const insert = (t) => {
    const a = input.selectionStart ?? input.value.length;
    const b = input.selectionEnd ?? a;
    input.setRangeText(t, a, b, "end");
    input.focus();
  };

  const KEYS = [
    ["C", "op", T.clear], ["(", "op"], [")", "op"], ["÷", "op"],
    ["7"], ["8"], ["9"], ["×", "op"],
    ["4"], ["5"], ["6"], ["−", "op"],
    ["1"], ["2"], ["3"], ["+", "op"],
    ["0"], [","], ["⌫", "op", T.back], ["=", "eq"],
  ];
  for (const [label, cls, aria] of KEYS) {
    const b = el("button", cls || "", label);
    b.type = "button";
    if (aria) b.setAttribute("aria-label", aria);
    b.addEventListener("click", () => {
      if (label === "C") { input.value = ""; out.className = ""; out.textContent = "0"; input.focus(); }
      else if (label === "⌫") { input.value = input.value.slice(0, -1); input.focus(); }
      else if (label === "=") equals();
      else insert(label);
    });
    keys.append(b);
  }
  input.addEventListener("keydown", (e) => {
    if (e.key === "Enter" || e.key === "=") { e.preventDefault(); equals(); }
    else if (e.key === "Escape") { e.preventDefault(); input.value = ""; }
  });

  // The shell's look; cd.js may be missing when the page is opened alone.
  const paint = (theme, accent) => {
    document.documentElement.dataset.theme = theme === "light" ? "light" : "dark";
    if (/^#[0-9a-fA-F]{3,8}$/.test(String(accent))) {
      document.documentElement.style.setProperty("--accent", accent);
    }
  };
  if (typeof cd === "object" && cd) {
    paint(cd.theme, cd.accent);
    cd.onTheme(paint);
  }

  document.body.append(box);
  input.focus();
})();
