// Saving a file an app hands over (cd.save). The app's frame has no
// allow-downloads and keeps it that way: an app that could start downloads on
// its own could drop any file it likes onto the person's disk. The shell does
// it instead, and only for what passes the checks below - one plain file name,
// a bounded size, a type that is a type - and always as a download, never as a
// page: the file is not opened, so nothing in it runs under any origin.

// 64 MiB, agreed with the apps part. A file bigger than this belongs to a
// stream the app reads from its own service, not to one message.
export const MAX_SAVE = 64 << 20;
export const MAX_SAVE_NAME = 255;

const TYPE = /^[a-z0-9.+-]+\/[a-z0-9.+-]+$/;
const enc = new TextEncoder();

// validSaveName is one file name and nothing more: no path, no control
// character, not "." or "..", no space at either end, at most 255 bytes of
// UTF-8. A name that fails is refused, never mended: the app is told, and a
// name the person did not see coming never lands on the disk.
export function validSaveName(n) {
  if (typeof n !== "string" || n === "" || n === "." || n === "..") return false;
  if (n !== n.trim()) return false;
  if (enc.encode(n).length > MAX_SAVE_NAME) return false;
  if (n.isWellFormed && !n.isWellFormed()) return false;
  // / and \ are paths, the rest are control characters (C0, DEL, C1).
  return !/[/\\\u0000-\u001f\u007f-\u009f]/.test(n);
}

// saveError says why a request is refused, "" when it may be saved. The codes
// are the ones the app sees: "name", "toolarge", "unavailable".
export function saveError(m) {
  if (!validSaveName(m.name)) return "name";
  const d = m.data;
  if (!(d instanceof Uint8Array) && !(d instanceof ArrayBuffer)) return "unavailable";
  if (d.byteLength > MAX_SAVE) return "toolarge";
  return "";
}

export const saveType = (t) => (typeof t === "string" && t.length <= 255 && TYPE.test(t))
  ? t : "application/octet-stream";

// saveFile hands the bytes to the browser's own download. The browser asks
// where to put it if the person has told it to; that is the question before
// saving, and it is theirs to switch on, not the app's to switch off.
export function saveFile(m) {
  // A copy: the app's buffer stays the app's, and the URL lives only for this
  // one click.
  const bytes = new Uint8Array(m.data instanceof ArrayBuffer ? m.data : m.data.slice().buffer);
  const url = URL.createObjectURL(new Blob([bytes], { type: saveType(m.type) }));
  const a = document.createElement("a");
  a.href = url;
  a.download = m.name;
  a.rel = "noopener";
  document.body.append(a);
  try {
    a.click();
  } finally {
    a.remove();
    // Not at once: some engines start reading the blob after click returns.
    setTimeout(() => URL.revokeObjectURL(url), 60_000);
  }
}
