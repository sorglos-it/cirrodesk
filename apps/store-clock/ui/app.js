// org.sorglos.clock - time and date, analog and digital, in the browser's
// own time zone. A classic script (see store-calculator); under node the
// pure helper is exported for the tests and no DOM is touched.
(function () {
  "use strict";

  // hands returns the angles of the three hands in degrees, 0 at twelve,
  // clockwise. The hour and minute hands move smoothly, the second hand jumps.
  function hands(h, m, s) {
    return { h: ((h % 12) + m / 60 + s / 3600) * 30, m: (m + s / 60) * 6, s: s * 6 };
  }

  const TITLES = { de: "Uhr", en: "Clock", fr: "Horloge", it: "Orologio", es: "Reloj" };

  function pickLang(tags) {
    for (const t of tags) {
      const l = String(t || "").slice(0, 2).toLowerCase();
      if (Object.prototype.hasOwnProperty.call(TITLES, l)) return l;
    }
    return "en";
  }

  if (typeof document === "undefined") {
    module.exports = { hands, pickLang, TITLES };
    return;
  }

  const lang = pickLang([window.cd && cd.lang, document.documentElement.lang]);
  document.documentElement.lang = lang;
  document.title = TITLES[lang];

  const NS = "http://www.w3.org/2000/svg";
  const svg = (tag, attrs) => {
    const e = document.createElementNS(NS, tag);
    for (const k in attrs) e.setAttribute(k, attrs[k]);
    return e;
  };
  const box = document.createElement("div");
  box.className = "clock";
  const dial = svg("svg", { viewBox: "-100 -100 200 200", "aria-hidden": "true" });
  dial.append(svg("circle", { class: "face", r: 96 }));
  for (let i = 0; i < 12; i++) {
    const t = svg("line", { class: "tick", x1: 0, y1: -88, x2: 0, y2: i % 3 ? -82 : -76 });
    t.setAttribute("transform", "rotate(" + i * 30 + ")");
    dial.append(t);
  }
  const hand = (cls, len) => {
    const l = svg("line", { class: cls, x1: 0, y1: 12, x2: 0, y2: -len });
    dial.append(l);
    return l;
  };
  const H = hand("h", 50), M = hand("m", 74), S = hand("s", 82);
  const time = document.createElement("time");
  time.className = "time";
  const date = document.createElement("div");
  date.className = "date";
  box.append(dial, time, date);

  const fTime = new Intl.DateTimeFormat(lang, { hour: "2-digit", minute: "2-digit", second: "2-digit" });
  const fDate = new Intl.DateTimeFormat(lang, { weekday: "long", year: "numeric", month: "long", day: "numeric" });

  function tick() {
    const now = new Date();
    const a = hands(now.getHours(), now.getMinutes(), now.getSeconds());
    H.setAttribute("transform", "rotate(" + a.h + ")");
    M.setAttribute("transform", "rotate(" + a.m + ")");
    S.setAttribute("transform", "rotate(" + a.s + ")");
    time.textContent = fTime.format(now);
    time.dateTime = now.toISOString();
    date.textContent = fDate.format(now);
    // To the next full second, so the second hand does not drift against it.
    setTimeout(tick, 1000 - now.getMilliseconds());
  }

  const paint = (theme, accent) => {
    document.documentElement.dataset.theme = theme === "light" ? "light" : "dark";
    if (/^#[0-9a-fA-F]{3,8}$/.test(String(accent))) document.documentElement.style.setProperty("--accent", accent);
  };
  if (typeof cd === "object" && cd) { paint(cd.theme, cd.accent); cd.onTheme(paint); }

  document.body.append(box);
  tick();
})();
