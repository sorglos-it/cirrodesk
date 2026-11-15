// org.sorglos.files - the user's own files: list, folders, upload, rename,
// delete, edit text and view pictures in the same window (Thomas, F1 = A).
// The service is ngc-files (../classes); its answers are JSON for list,
// mkdir, rename and remove, raw bytes on the channels get and put.
//
// A classic script, not a module (see store-calculator). Under node (the
// tests) the pure helpers are exported and no DOM is touched.
(function () {
  "use strict";

  // The same rule as ValidPath of the service, for one name: the UI refuses
  // early, the service decides.
  function validName(s) {
    if (typeof s !== "string" || s === "" || s === "." || s === ".." || s.startsWith(".upload-")) return false;
    if (new TextEncoder().encode(s).length > 255 || s !== s.trim()) return false;
    return !/[\/\\\u0000-\u001f\u007f-\u009f�]/.test(s);
  }

  const join = (dir, name) => (dir === "" ? name : dir + "/" + name);

  // Only what the browser shows safely inline: text as text, pictures as
  // data: images (the frame's policy allows img-src data: and nothing else).
  const IMAGES = { __proto__: null, png: "image/png", jpg: "image/jpeg", jpeg: "image/jpeg", gif: "image/gif", webp: "image/webp", bmp: "image/bmp" };
  const TEXTS_EXT = new Set(["txt", "md", "csv", "json", "log", "ini", "conf", "yaml", "yml", "xml", "html", "css", "js", "sh", "go", "py"]);
  const ext = (n) => { const i = n.lastIndexOf("."); return i > 0 ? n.slice(i + 1).toLowerCase() : ""; };
  const kind = (n) => (IMAGES[ext(n)] ? "image" : TEXTS_EXT.has(ext(n)) ? "text" : "");

  // The largest file opened in the window; larger ones are only downloaded.
  const MAX_OPEN = 8 << 20;

  function size(n) {
    if (n < 1024) return n + " B";
    const u = ["KB", "MB", "GB"];
    let i = -1;
    do { n /= 1024; i++; } while (n >= 1024 && i < u.length - 1);
    return n.toFixed(n < 10 ? 1 : 0) + " " + u[i];
  }

  // decode reads an answer of the service: bytes (or a string) holding JSON.
  function decode(v) {
    const s = typeof v === "string" ? v : new TextDecoder().decode(v);
    const a = JSON.parse(s);
    if (!a || typeof a !== "object" || Array.isArray(a)) throw new Error("answer");
    return a;
  }

  function b64(bytes) {
    let s = "";
    for (let i = 0; i < bytes.length; i += 0x8000) s += String.fromCharCode.apply(null, bytes.subarray(i, i + 0x8000));
    return btoa(s);
  }

  const TEXTS = {
    de: { title: "Dateien", home: "Meine Dateien", up: "Hoch", newdir: "Neuer Ordner", upload: "Hochladen",
      rename: "Umbenennen", remove: "Löschen", download: "Herunterladen", save: "Speichern", close: "Schließen",
      empty: "Dieser Ordner ist leer.", cut: "Nicht alle Einträge werden gezeigt.", name: "Name",
      askName: "Name:", askRemove: "Wirklich löschen?", noDownload: "Download nicht verfügbar.",
      tooBig: "Die Datei ist zu groß zum Öffnen.", saved: "Gespeichert.", busy: "Bitte warten …",
      call: "Die Anfrage wurde abgelehnt.", path: "Dieser Name ist nicht erlaubt.", notfound: "Nicht gefunden.",
      exists: "Dieser Name ist schon vergeben.", notempty: "Der Ordner ist nicht leer.",
      toolarge: "Die Datei ist zu groß.", failed: "Das hat nicht geklappt.", offline: "Der Dienst antwortet nicht." },
    en: { title: "Files", home: "My files", up: "Up", newdir: "New folder", upload: "Upload",
      rename: "Rename", remove: "Delete", download: "Download", save: "Save", close: "Close",
      empty: "This folder is empty.", cut: "Not all entries are shown.", name: "Name",
      askName: "Name:", askRemove: "Really delete?", noDownload: "Download not available.",
      tooBig: "The file is too large to open.", saved: "Saved.", busy: "Please wait …",
      call: "The request was refused.", path: "This name is not allowed.", notfound: "Not found.",
      exists: "This name is already taken.", notempty: "The folder is not empty.",
      toolarge: "The file is too large.", failed: "That did not work.", offline: "The service does not answer." },
    fr: { title: "Fichiers", home: "Mes fichiers", up: "Remonter", newdir: "Nouveau dossier", upload: "Téléverser",
      rename: "Renommer", remove: "Supprimer", download: "Télécharger", save: "Enregistrer", close: "Fermer",
      empty: "Ce dossier est vide.", cut: "Toutes les entrées ne sont pas affichées.", name: "Nom",
      askName: "Nom :", askRemove: "Vraiment supprimer ?", noDownload: "Téléchargement indisponible.",
      tooBig: "Le fichier est trop grand pour être ouvert.", saved: "Enregistré.", busy: "Veuillez patienter …",
      call: "La demande a été refusée.", path: "Ce nom n'est pas autorisé.", notfound: "Introuvable.",
      exists: "Ce nom est déjà pris.", notempty: "Le dossier n'est pas vide.",
      toolarge: "Le fichier est trop grand.", failed: "Cela n'a pas fonctionné.", offline: "Le service ne répond pas." },
    it: { title: "File", home: "I miei file", up: "Su", newdir: "Nuova cartella", upload: "Carica",
      rename: "Rinomina", remove: "Elimina", download: "Scarica", save: "Salva", close: "Chiudi",
      empty: "Questa cartella è vuota.", cut: "Non tutte le voci sono mostrate.", name: "Nome",
      askName: "Nome:", askRemove: "Eliminare davvero?", noDownload: "Download non disponibile.",
      tooBig: "Il file è troppo grande da aprire.", saved: "Salvato.", busy: "Attendere …",
      call: "La richiesta è stata rifiutata.", path: "Questo nome non è consentito.", notfound: "Non trovato.",
      exists: "Questo nome è già in uso.", notempty: "La cartella non è vuota.",
      toolarge: "Il file è troppo grande.", failed: "Non ha funzionato.", offline: "Il servizio non risponde." },
    es: { title: "Archivos", home: "Mis archivos", up: "Subir nivel", newdir: "Nueva carpeta", upload: "Subir",
      rename: "Renombrar", remove: "Eliminar", download: "Descargar", save: "Guardar", close: "Cerrar",
      empty: "Esta carpeta está vacía.", cut: "No se muestran todas las entradas.", name: "Nombre",
      askName: "Nombre:", askRemove: "¿Eliminar de verdad?", noDownload: "Descarga no disponible.",
      tooBig: "El archivo es demasiado grande para abrirlo.", saved: "Guardado.", busy: "Espera …",
      call: "La solicitud fue rechazada.", path: "Este nombre no está permitido.", notfound: "No encontrado.",
      exists: "Este nombre ya existe.", notempty: "La carpeta no está vacía.",
      toolarge: "El archivo es demasiado grande.", failed: "No ha funcionado.", offline: "El servicio no responde." },
  };

  function pickLang(tags) {
    for (const t of tags) {
      const l = String(t || "").slice(0, 2).toLowerCase();
      if (Object.prototype.hasOwnProperty.call(TEXTS, l)) return l;
    }
    return "en";
  }

  if (typeof document === "undefined") {
    module.exports = { validName, join, kind, size, decode, b64, pickLang, TEXTS, MAX_OPEN };
    return;
  }

  const lang = pickLang([window.ngc && ngc.lang, document.documentElement.lang]);
  const T = TEXTS[lang];
  document.documentElement.lang = lang;
  document.title = T.title;

  function el(tag, cls, txt) {
    const e = document.createElement(tag);
    if (cls) e.className = cls;
    if (txt !== undefined) e.textContent = String(txt);
    return e;
  }
  function button(label, fn, cls) {
    const b = el("button", cls || "", label);
    b.type = "button";
    b.addEventListener("click", fn);
    return b;
  }

  const api = typeof ngc === "object" && ngc ? ngc : null;

  // ask runs one call and returns its answer; the error code becomes a thrown
  // Error whose message is that code.
  async function ask(action, params) {
    if (!api) throw new Error("offline");
    let a;
    try {
      a = decode(await api.call("files", action, params));
    } catch {
      throw new Error("offline");
    }
    if (a.error) throw new Error(a.error);
    return a;
  }

  // fetchFile reads a file over the get channel, at most limit bytes.
  function fetchFile(path, limit) {
    return new Promise((resolve, reject) => {
      if (!api) return reject(new Error("offline"));
      api.open("files", "get", { path }).then((s) => {
        const parts = [];
        let n = 0;
        s.onData = (d) => {
          n += d.length;
          if (n > limit) { reject(new Error("toolarge")); try { s.close && s.close(); } catch {} return; }
          parts.push(d);
        };
        s.onClose = (status) => {
          if (status !== "ok") return reject(new Error("failed"));
          const out = new Uint8Array(n);
          let o = 0;
          for (const p of parts) { out.set(p, o); o += p.length; }
          resolve(out);
        };
      }, () => reject(new Error("offline")));
    });
  }

  // storeFile writes bytes over the put channel in pieces.
  function storeFile(path, bytes) {
    return new Promise((resolve, reject) => {
      if (!api) return reject(new Error("offline"));
      api.open("files", "put", { path }).then(async (s) => {
        s.onClose = (status) => (status === "ok" ? resolve() : reject(new Error("failed")));
        try {
          for (let i = 0; i < bytes.length; i += 64 << 10) await s.write(bytes.subarray(i, i + (64 << 10)));
          await s.end();
        } catch {
          reject(new Error("offline"));
        }
      }, () => reject(new Error("offline")));
    });
  }

  // ---- the page ----
  const box = el("div", "files");
  const bar = el("div", "bar");
  const crumbs = el("nav", "crumbs");
  const status = el("p", "status");
  status.setAttribute("role", "status");
  const list = el("ul", "list");
  const view = el("section", "view");
  view.hidden = true;
  const picker = el("input");
  picker.type = "file";
  picker.multiple = true;
  picker.hidden = true;
  box.append(bar, crumbs, list, view, status, picker);

  let dir = "";
  let selected = null;

  // The frame has no allow-modals, so prompt and confirm answer nothing: the
  // question is asked in a row of the page instead.
  const question = el("form", "question");
  question.hidden = true;
  box.insertBefore(question, crumbs);
  function inline(label, value) {
    return new Promise((resolve) => {
      question.replaceChildren();
      const l = el("label", "", label);
      const field = el("input");
      field.type = "text";
      field.value = value === null ? "" : value;
      field.hidden = value === null;
      l.append(field);
      const ok = el("button", "", "OK");
      ok.type = "submit";
      const no = button(T.close, () => done(null));
      const done = (v) => { question.hidden = true; question.onsubmit = null; resolve(v); };
      question.onsubmit = (ev) => { ev.preventDefault(); done(value === null ? "" : field.value); };
      question.append(l, ok, no);
      question.hidden = false;
      (value === null ? ok : field).focus();
    });
  }

  const say = (msg, bad) => {
    status.className = bad ? "status bad" : "status";
    status.textContent = msg || "";
  };
  const fail = (e) => say(T[e && e.message] || T.failed, true);

  const bUp = button(T.up, () => { dir = dir.includes("/") ? dir.slice(0, dir.lastIndexOf("/")) : ""; load(); });
  const bNew = button(T.newdir, async () => {
    const n = await inline(T.askName, "");
    if (n === null) return;
    if (!validName(n)) return say(T.path, true);
    try { await ask("mkdir", { path: join(dir, n) }); load(); } catch (e) { fail(e); }
  });
  const bUpload = button(T.upload, () => picker.click());
  const bRename = button(T.rename, async () => {
    if (!selected) return;
    const n = await inline(T.askName, selected.name);
    if (n === null || n === selected.name) return;
    if (!validName(n)) return say(T.path, true);
    try { await ask("rename", { path: join(dir, selected.name), to: join(dir, n) }); load(); } catch (e) { fail(e); }
  });
  const bRemove = button(T.remove, async () => {
    if (!selected || (await inline(T.askRemove + " " + selected.name, null)) === null) return;
    try { await ask("remove", { path: join(dir, selected.name) }); load(); } catch (e) { fail(e); }
  });
  const bDownload = button(T.download, async () => {
    if (!selected || selected.dir) return;
    // The frame may not save files itself (sandbox without allow-downloads):
    // the shell does it (Thomas, F2 = A). Without that part there is no way.
    if (!api || typeof api.save !== "function") return say(T.noDownload, true);
    say(T.busy);
    try {
      const bytes = await fetchFile(join(dir, selected.name), 64 << 20);
      await api.save(selected.name, bytes, IMAGES[ext(selected.name)] || "application/octet-stream");
      say("");
    } catch (e) {
      // The shell's codes for ngc.save, agreed with Web: refused, toolarge,
      // name, unavailable. Only the browser may ask before saving, so
      // "refused" comes only if the shell later asks itself.
      const m = { refused: "", name: "path", unavailable: "noDownload" }[e && e.message];
      if (m === "") say("");
      else if (m) say(T[m], true);
      else fail(e);
    }
  });
  bar.append(bUp, bNew, bUpload, bRename, bRemove, bDownload);

  const refresh = () => {
    bUp.disabled = dir === "";
    bRename.disabled = bRemove.disabled = !selected;
    bDownload.disabled = !selected || selected.dir;
  };

  picker.addEventListener("change", async () => {
    const files = [...picker.files];
    picker.value = "";
    for (const f of files) {
      if (!validName(f.name)) { say(T.path + " " + f.name, true); continue; }
      say(T.busy + " " + f.name);
      try { await storeFile(join(dir, f.name), new Uint8Array(await f.arrayBuffer())); say(""); } catch (e) { fail(e); }
    }
    load();
  });

  function showCrumbs() {
    crumbs.replaceChildren();
    const parts = dir === "" ? [] : dir.split("/");
    const go = (d) => () => { dir = d; load(); };
    crumbs.append(button(T.home, go(""), "link"));
    parts.forEach((p, i) => {
      crumbs.append(el("span", "sep", "/"));
      crumbs.append(button(p, go(parts.slice(0, i + 1).join("/")), "link"));
    });
  }

  async function load() {
    selected = null;
    closeView();
    showCrumbs();
    refresh();
    let a;
    try { a = await ask("list", { path: dir }); } catch (e) { list.replaceChildren(); return fail(e); }
    list.replaceChildren();
    const entries = Array.isArray(a.entries) ? a.entries : [];
    if (entries.length === 0) list.append(el("li", "none", T.empty));
    for (const e of entries) {
      if (!e || typeof e.name !== "string") continue;
      const li = el("li", e.dir ? "dir" : "file");
      const b = button("", () => {
        for (const o of list.querySelectorAll("[aria-pressed]")) o.setAttribute("aria-pressed", "false");
        b.setAttribute("aria-pressed", "true");
        selected = e;
        refresh();
      });
      b.setAttribute("aria-pressed", "false");
      b.append(el("span", "icon", e.dir ? "📁" : kind(e.name) === "image" ? "🖼" : "📄"), el("span", "name", e.name),
        el("span", "size", e.dir ? "" : size(Number(e.size) || 0)));
      b.addEventListener("dblclick", () => (e.dir ? ((dir = join(dir, e.name)), load()) : open(e)));
      b.addEventListener("keydown", (k) => { if (k.key === "Enter") b.dispatchEvent(new MouseEvent("dblclick")); });
      li.append(b);
      list.append(li);
    }
    say(a.cut ? T.cut : "");
  }

  function closeView() {
    view.hidden = true;
    list.hidden = false;
    view.replaceChildren();
  }

  async function open(e) {
    const k = kind(e.name);
    if (!k) return;
    if (Number(e.size) > MAX_OPEN) return say(T.tooBig, true);
    const path = join(dir, e.name);
    say(T.busy);
    let bytes;
    try { bytes = await fetchFile(path, MAX_OPEN); } catch (err) { return fail(err); }
    say("");
    view.replaceChildren();
    const head = el("div", "bar");
    head.append(el("strong", "", e.name));
    view.append(head);
    if (k === "image") {
      const img = el("img");
      img.alt = e.name;
      img.src = "data:" + IMAGES[ext(e.name)] + ";base64," + b64(bytes);
      view.append(img);
    } else {
      const ta = el("textarea");
      ta.spellcheck = false;
      ta.value = new TextDecoder().decode(bytes);
      ta.setAttribute("aria-label", e.name);
      head.append(button(T.save, async () => {
        say(T.busy);
        try { await storeFile(path, new TextEncoder().encode(ta.value)); say(T.saved); } catch (err) { fail(err); }
      }));
      view.append(ta);
    }
    head.append(button(T.close, () => load()));
    list.hidden = true;
    view.hidden = false;
  }

  const paint = (theme, accent) => {
    document.documentElement.dataset.theme = theme === "light" ? "light" : "dark";
    if (/^#[0-9a-fA-F]{3,8}$/.test(String(accent))) document.documentElement.style.setProperty("--accent", accent);
  };
  if (api) { paint(api.theme, api.accent); api.onTheme(paint); }

  document.body.append(box);
  load();
})();
