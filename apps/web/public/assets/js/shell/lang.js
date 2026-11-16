// The sentences the masks say, in the five languages of the house rule
// (de, en, fr, it, es). One file per language under /assets/lang/, the same
// keys in every one of them, and a new language is one more file and no line
// of code here - lang.js never holds a list of languages for its own sake,
// it holds the names of the files it may ask for.
//
// WHY THE FILES LIE IN public/assets/lang/ AND NOT IN apps/web/lang/, which
// is where the project standard puts translations: the browser has to fetch
// them, and only public/ is reachable from outside (apps/web/CLAUDE.md, and
// the one entry point of public/index.php). A lang/ outside the web root
// would need a second route through PHP, which is the one thing the README
// of this app promises there is not.
//
// WHAT IS HERE AND WHAT IS NOT. The keys are the sentences of the login mask
// and of the first run - every sentence a person reads when something went
// wrong. The REST OF THE SHELL IS STILL GERMAN IN ITS SOURCE, and moving it
// is its own task. How much there is of it, with the command, because a
// number without one is no number:
//
//   grep -rhoE '"[A-ZAOU][^"]{4,}"' auth/ shell/ apps/ | sort -u | wc -l
//
// answers 283 here (294 before the keys below left those files). It is an
// UPPER BOUND and not a count of sentences: it catches "Enter" and "Orange"
// too, which are a key name and a colour and no text for anybody. So: this
// is a mechanism with the error sentences in it, not a translated shell, and
// nothing here pretends otherwise.
//
// Fallback is per KEY and not per file, which is what the house rule asks
// for: a key the chosen language does not have is answered out of English,
// and a key English does not have either is answered out of BUILT_IN below.
// Only a key that is in no table at all comes back as "" - never the naked
// key, which is no sentence for anybody.
//
// WHY BUILT_IN EXISTS, and it is MEASURED, not foreseen. The sentences used
// to be constants in auth/login.js and auth/setup.js, so a mask always had
// them; with the files they are one fetch away, and a browser that cannot
// get the files has no sentence at all. Measured in chromium, rendered: with
// every /assets/lang/*.json refused and a wrong password typed, the alert of
// the login mask held "" and was not even laid out (offsetHeight 0). Before
// the files it read "Anmeldung fehlgeschlagen." A release that ships
// assets/js/ and forgets assets/lang/ is exactly that state - the release
// packer names what goes out file by file (tools/release/extra-web.txt, on
// the unmerged release-packer branch) - and so is a disk that fills up while
// the page is already open.
//
// It is NOT a second table that may drift: the English file and BUILT_IN are
// held key for key and value for value by a case of shell/busy.spec.js, so a
// sentence changed in one of them and not the other is a red case and not a
// surprise at somebody's screen.
import { decode, Budget, M, str } from "../contracts/json.js";
import { readBounded } from "../net/http.js";

// The files that may be asked for. A name from outside never reaches the
// address: what the browser offers is held against this list, so a tag like
// "../../etc/cd/web.json" cannot become a path. That is the whole reason
// this list exists - not to know which languages there are.
export const LANGS = Object.freeze(["de", "en", "fr", "it", "es"]);

// The language a file is written in when nothing else fits, and the one every
// missing key is answered out of.
export const FALLBACK = "en";

// MEASURED with `wc -c public/assets/lang/*.json`: the five files are 1 278
// to 1 482 bytes, the largest of them French. The bound is here so an answer
// that is NOT a language file is dropped while it arrives rather than after
// it is all in memory - the same reason net/config.js bounds /config.json.
export const MAX_LANG = 8 << 10;

// The English sentences, compiled in, so that a mask is never silent. It is a
// COPY of public/assets/lang/en.json and says so: the file is what somebody
// translates and what the browser reads, this is the floor under it for the
// one case where no file could be read at all. The case that holds the two
// equal is "the built-in sentences are the English file, word for word" in
// shell/busy.spec.js.
// THE BROWSERS THIS SHELL KNOWS CANNOT DO A PASSWORD LOGIN, as data and in one
// place, because two sentences have to agree about them.
//
// MEASURED (18.10.2026, Playwright's Firefox 141 on Windows): exportKey("jwk")
// of an Ed25519 key imported from a PKCS#8 seed answers OperationError, so
// contracts/proof.js publicFromSeed has no way to the public half and no
// password gets in. shell/main.js passwordLoginFits() measures it per engine at
// run time - this list is not a gate and must never become one. It is only
// here so the ADVICE the shell gives cannot recommend one of them: a case holds
// boot.tooold against this list in all five languages
// (apps/web/tests/shell/engines.spec.js).
//
// A name goes OUT of here the day an engine is measured to do it, and then the
// advice may name it again.
export const NO_PASSWORD_LOGIN = Object.freeze(["Firefox"]);

export const BUILT_IN = Object.freeze({
  __proto__: null,
  "boot.tooold": "This browser cannot create the keys CirroDesk needs for a login. Please use a current Chrome or Safari.",
  "boot.noroute.password": "This browser cannot do CirroDesk's password login: it does not hand out the public half of a password key. Please use a current Chrome or Safari. Firefox cannot do it at present.",
  "boot.noroute.passkeysoff": "No login is possible in this browser: it cannot do the password login, and passkeys are switched off here. Please switch passkeys on in the browser settings, or use a current Chrome or Safari.",
  "err.unreachable": "Service unavailable. Please try again later.",
  "err.busy": "The service is busy right now. Please try again in {seconds} seconds.",
  "login.failed": "Login failed.",
  "login.nopasswordkeys": "This browser cannot do a password login. Please log in with a passkey.",
  "login.locked": "Too many failed attempts. Please try again later.",
  "login.locked.left": "Too many failed attempts. {seconds} seconds left.",
  "login.poweredby": "running on CirroDesk {version}",
  "login.totp.ask": "Please enter the code from your authenticator app.",
  "login.head": "Sign in to CirroDesk",
  "login.name": "User name",
  "login.password": "Password",
  "login.totp": "Code from the authenticator app",
  "login.submit": "Sign in",
  "login.passkey": "Sign in with a passkey",
  "setup.code.incomplete": "The one-time code is not complete.",
  "setup.code.nomatch": "The one-time code does not fit this installation. Please check the code and try again.",
  "setup.silent": "The installation is not answering right now. The one-time code stays valid. Please try again in a moment.",
  "setup.refused": "The installation refused the setup. Please install CirroDesk again and use the new code.",
  "setup.already": "This installation is already set up.",
  "setup.broken": "The setup did not work. Please reload the page and try again.",
  "setup.name.bad": "A login name may hold only lower-case letters, digits and . _ - and must begin with a letter or a digit.",
  "setup.password.missing": "Please choose a password.",
  "setup.password.different": "The two passwords are not the same.",
  "setup.notready": "This installation is not ready to use yet. Its operator still has to finish the setup.",
  "setup.signin": "CirroDesk is set up. Please sign in now.",
  // The six accent colours of the start menu. They are in here and not only
  // in the files for the reason the comment above gives: without a sentence
  // the six swatches have no name at all, and a name is the only way to hit
  // them - there is no colour picker to fall back on.
  "look.accent": "Color {name}",
  "look.accent.blue": "Blue",
  "look.accent.teal": "Turquoise",
  "look.accent.pink": "Pink",
  "look.accent.orange": "Orange",
  "look.accent.violet": "Violet",
  "look.accent.green": "Green",
  // The band and the window commands. In here for the same reason as the rest:
  // without a sentence the mark at the left of the band has no name, the dock
  // has no name, and the menu of the app in front is four empty buttons.
  "bar.start": "Start",
  "bar.dock": "Dock",
  "win.minimize": "Minimise",
  "win.maximize": "Maximise",
  "win.restore": "Restore",
  "win.tile": "Arrange windows",
  "win.close": "Close",
  "desk.theme": "Switch design",
  "desk.search": "Search apps",
  "desk.allapps": "All apps",
  "desk.closeall": "Close all",
  "desk.noticeclose": "Close message",
  "desk.nohits": "No matches.",
  "desk.empty": "Your apps will appear here.",
  "desk.appsfailed": "The list of apps could not be loaded.",
  "desk.openstore": "Open store",
  "desk.adminactivates": "An administrator activates the apps of this installation.",
  "desk.nolicense": "No valid licence.",
  "desk.appfailed": "This app could not be opened.",
  "desk.appfailed.hint": "Please close the window and try again.",
  "look.head": "Appearance",
  "look.theme.system": "Like the system",
  "look.theme.light": "Light",
  "look.theme.dark": "Dark",
  "look.wallpaper": "Background {name}",
  "look.wallpaper.deep": "Space",
  "look.wallpaper.aurora": "Aurora",
  "look.wallpaper.sunset": "Sunset",
  "look.wallpaper.mint": "Mint",
  "look.wallpaper.plum": "Plum",
  "look.density.eng": "Compact",
  "look.density.normal": "Standard",
  "look.density.luftig": "Spacious",
  "look.glass.aus": "No glass",
  "look.glass.normal": "Glass",
  "look.glass.stark": "More glass",
  "store.tab.sources": "Sources",
  "store.tab.apps": "Apps",
  "store.tab.installed": "Installed",
  "store.tab.upload": "Upload",
  "store.tab.updates": "Updates",
  "store.need.network": "Access to the network",
  "store.need.root": "Full rights in its container",
  "store.part.web": "CirroDesk web interface",
  "store.part.api": "CirroDesk API",
  "store.part.controller": "CirroDesk services",
  "store.part.system": "Operating system",
  "store.upd.available": "New version available",
  "store.upd.downloaded": "Ready to install",
  "store.upd.installing": "Installing",
  "store.upd.installed": "Installed",
  "store.upd.rolledback": "Rolled back",
  "store.upd.failed": "Failed",
  "store.upd.rejected": "Rejected",
  "store.upd.needsgrants": "Needs new rights",
  "store.app.ok": "Running",
  "store.app.unreachable": "Not reachable",
  "store.app.unlicensed": "No valid licence",
  "store.mode.auto": "automatic",
  "store.mode.semi": "on click",
  "store.mode.notify": "notify only",
  "store.wait": "One moment …",
  "store.broken": "That did not work. Please close the window and open it again.",
  "store.src.lead": "This is where CirroDesk gets its apps.",
  "store.src.checked": "Last checked: {when}",
  "store.src.never": "Not checked yet",
  "store.src.silent": "This source does not answer. Please check the address.",
  "store.src.none": "No source entered yet.",
  "store.apps.nosource": "No source yet. Please enter one under “Sources”.",
  "store.apps.source": "Source",
  "store.apps.untrusted": "This store is not trusted yet. Please enter it first.",
  "store.apps.notfetched": "This store's catalogue has not been fetched yet. CirroDesk checks every hour; Sources shows when it last worked.",
  "store.apps.old": "This store sends an old catalogue. Please try again later.",
  "store.apps.fake": "This store's catalogue is not genuine.",
  "store.apps.install": "Install",
  "store.apps.install.for": "Install {name}",
  "store.apps.buy": "Buy",
  "store.apps.buy.for": "Buy {name}",
  "store.apps.none": "This store has no app yet.",
  "store.lic.nokey": "This installation's key did not arrive. Please try again.",
  "store.lic.asks": "The shop asks for these two details:",
  "store.lic.installation": "Installation",
  "store.lic.key": "Key",
  "store.lic.copy": "Copy",
  "store.lic.copy.both": "Copy installation and key",
  "store.lic.shop": "Shop: {shop}",
  "store.lic.pickfile": "Choose the licence file from the shop here:",
  "store.lic.pick": "Choose licence",
  "store.lic.toolarge": "This file is too large for a licence. Please choose the file from the shop.",
  "store.lic.valid": "The licence is valid.",
  "store.lic.taken": "The licence was accepted.",
  "store.lic.refused": "This licence was not accepted. Please choose the file from the shop.",
  "store.copied": "Copied.",
  "store.nocopy": "Copying does not work here. Please select the text by hand.",
  "store.inst.off": "Switched off",
  "store.inst.remove": "Remove",
  "store.inst.remove.for": "Remove {name}",
  "store.inst.none": "No app set up yet.",
  "store.inst.removed": "The app has been removed.",
  "store.up.lead": "Set up a package from disk.",
  "store.up.pick": "Choose package",
  "store.up.notcd": "This file is not a CirroDesk package. Please choose a file ending in .cd.",
  "store.up.app": "App",
  "store.up.version": "Version",
  "store.up.sum": "Checksum",
  "store.up.uisum": "UI checksum",
  "store.up.sending": "Sending …",
  "store.up.progress": "Sending … {kb} KB",
  "store.up.lost": "The package did not arrive. Please choose it again.",
  "store.up.swapped": "The installation received a different package from the one sent.",
  "store.act.rights": "This app wants to do more than others. Please allow each right on its own:",
  "store.act.nodesfailed": "The machines could not be loaded. Please reopen the window.",
  "store.act.unreadable": "A detail of this server cannot be read ({what}).",
  "store.act.nonode": "No machine for services. Please add a machine first.",
  "store.act.node": "Machine",
  "store.act.go": "Activate",
  "store.act.allrights": "Please allow all rights. Without them the app cannot work.",
  "store.act.done": "The app is set up.",
  "store.upd.lead": "What new versions there are.",
  "store.upd.for": "Updates for {name}",
  "store.upd.now": "Update now",
  "store.upd.now.for": "Update {name} now",
  "store.upd.withrights": "Activate the new version with new rights",
  "store.upd.none": "Everything is up to date.",
  "store.upd.started": "The new version is being installed.",
  "store.failed.again": "That did not work. Please try again.",
  "store.upd.notstore": "This app did not come from a store. Please set up the new version under “Upload”.",
  "store.upd.notincat": "This version is not in the catalogue. Please try again later.",
  "store.save.failed": "That could not be saved. Please try again.",
  "store.saved": "Saved.",
  "store.needlogin": "This needs a login. Please sign in again.",
  "store.grant.nokey": "The key for this approval is missing. Please enter this app's store first.",
  "store.nochange": "Nothing changed.",
  "store.grant.invalid": "This approval does not make a valid trust list. Please sign in again.",
  "store.grant.refused": "The approval was not accepted. Please try again.",
  "store.grant.reload": "The approval is stored. Please reload the page so CirroDesk is up to date here.",
  "store.nosign": "That could not be signed. Please sign in again.",
  "store.notrust": "The trust list is missing. Please sign in again.",
  "store.confirm.short": "The confirmation was not enough. Please try again.",
  "store.node.refused": "The machine did not accept that. Please try again.",
  "store.node.says": "The machine says: {text}",
  "users.bad.list": "The installation's answer was not in order. Nothing is shown.",
  "users.bad.offer": "The installation's answer was not in order. Nothing is offered.",
  "users.noreach": "The installation is not answering right now.",
  "users.nofp": "The fingerprint is still being calculated. Please try again.",
  "users.behind": "The server is entered. Please reload the page so CirroDesk is up to date here.",
  "users.root.none": "With system rights: nobody",
  "users.kind.login": "Sign-in",
  "users.kind.login-failed": "Sign-in failed",
  "users.kind.channel": "Channel",
  "users.kind.task": "Task",
  "users.kind.trust": "Trust list",
  "users.kind.rights": "Rights",
  "users.kind.app": "App",
  "users.kind.store": "Store",
  "users.kind.license": "Licence",
  "users.kind.update": "Update",
  "users.tab.users": "Users",
  "users.tab.rights": "Roles and rights",
  "users.tab.servers": "Servers",
  "users.tab.log": "Log",
  "users.title": "Users and rights",
  "users.root.legend": "Tasks with system rights",
  "users.root.hint": "Whoever is ticked here may set up and remove apps, delete data and open the terminal on this server.",
  "users.root.user.aria": "User {name} may run tasks with system rights",
  "users.root.people": "Individual people",
  "users.root.role.aria": "Role {name} may run tasks with system rights",
  "users.root.roles": "Whole roles",
  "users.root.later": "A role also applies to everyone who joins it later.",
  "users.role.named": "Role {name}",
  "users.root.who": "With system rights: {list}",
  "users.root.toomany": "Too many boxes are ticked. A server cannot hold more than {max} people and more than {max} roles.",
  "users.col.name": "Name",
  "users.col.roles": "Roles",
  "users.col.totp": "Two-step sign-in",
  "users.col.locked": "Locked",
  "users.col.last": "Last signed in",
  "users.users.head": "Users",
  "users.on": "on",
  "users.off": "off",
  "users.locked": "locked",
  "users.never": "never",
  "users.pk.legend": "Passkey",
  "users.pk.cannot": "No passkey can be entered right now.",
  "users.pk.nobrowser": "This browser cannot create passkeys.",
  "users.pk.add": "Add passkey",
  "users.pk.add.aria": "Add a passkey for me",
  "users.pk.hint": "A passkey on this device or on a stick, with which you will sign in without a password from now on.",
  "users.pk.notlisted": "Your user is not in the trust list.",
  "users.pk.toomany": "You already have {max} keys. The list holds no more.",
  "users.pk.wrongrp": "This address does not match this installation's passkeys.",
  "users.pk.ask": "Your browser will ask for the passkey in a moment. Then it is entered into the trust list; you confirm once more for that.",
  "users.next": "Continue",
  "users.pk.done": "The passkey is entered.",
  "users.pk.failed": "The passkey could not be entered.",
  "users.new.legend": "New user",
  "users.new.cannot": "No user can be created right now.",
  "users.new.login": "Login name",
  "users.new.pw1": "First password",
  "users.new.pw2": "Repeat first password",
  "users.new.go": "Create user",
  "users.new.go.aria": "Create a new user",
  "users.new.namerule": "Lower-case letters, digits and . _ - , the first character a letter or a digit, at most 64 characters.",
  "users.new.roles": "Roles",
  "users.new.norole": "Nobody can be created without a role, and an app is only seen by those whose role has the right “use” for it.",
  "users.new.badname": "“{name}” is not a valid login name.",
  "users.new.taken": "The login name “{name}” already exists.",
  "users.new.toomany": "This installation already has {max} users. The list holds no more.",
  "users.new.pickrole": "Choose at least one role. Nobody can be created without a role.",
  "users.new.hard": "Nobody can be created here for the role “{role}”: it demands a passkey, and this form sets up a password.",
  "users.new.nopw": "Enter a first password.",
  "users.new.pwdiffer": "The two passwords are not the same.",
  "users.new.title": "Create user “{name}”",
  "users.new.body": "“{name}” is entered into this installation's trust list and can then sign in with this password.",
  "users.new.bodyroles": "Roles: {list}.",
  "users.new.who": "You typed the first password yourself. Until “{name}” changes it, you know their key and could sign in their name. Hand it over in person and ask them to change it straight away.",
  "users.new.done": "“{name}” is created and can sign in.",
  "users.new.failed": "The user could not be created.",
  "users.rights.head": "Roles and rights",
  "users.rights.hint": "A tick means: everyone with this role may do this.",
  "users.role.legend": "New role",
  "users.role.cannot": "No role can be created right now.",
  "users.role.hint": "A role gathers rights. Create it, tick above what it may do, and then give it to a user.",
  "users.role.name": "Name of the role",
  "users.role.namerule": "Lower-case letters, digits and hyphens, the first letter not a digit, at most 32 characters.",
  "users.role.go": "Create role",
  "users.role.go.aria": "Create a new role",
  "users.role.badname": "“{name}” is not a valid role name.",
  "users.role.taken": "The role “{name}” already exists.",
  "users.role.toomany": "This installation already has {max} roles. The list holds no more.",
  "users.role.title": "Create role “{name}”",
  "users.role.body1": "The role is entered into this installation's trust list. At first it may do nothing.",
  "users.role.body2": "What it may do, you tick on this page afterwards. An app for which the role lacks the right “use” is not in the start menu of anyone with this role.",
  "users.role.done": "The role “{name}” is created. It may do nothing yet.",
  "users.role.failed": "The role could not be created.",
  "users.role.allapps": "This role may use every app that is set up. That is fixed and cannot be unticked.",
  "users.role.noapp": "No app is set up yet.",
  "users.right.aria": "Role {role}: {right}",
  "users.role.norights": "There is no app yet that there would be rights for.",
  "users.rights.save": "Save rights",
  "users.rights.save.aria": "Save the rights of role {role}",
  "users.rights.same": "Nothing has changed on the role “{role}”.",
  "users.rights.adds": "Everyone with the role “{role}” may then also:",
  "users.rights.gone": "Taken away:",
  "users.rights.title": "Change the rights of role “{role}”",
  "users.rights.go": "Grant rights",
  "users.rights.saved": "The rights of role “{role}” are saved.",
  "users.rights.failed": "The rights could not be saved.",
  "users.who.none": "Nobody has this role at the moment.",
  "users.who.one": "This applies at once to 1 person with this role: {names}.",
  "users.who.many": "This applies at once to {n} people with this role: {names}.",
  "users.srv.head": "Servers",
  "users.srv.pending": "Waiting servers",
  "users.srv.compare": "Compare this fingerprint with the one the script shows on the server.",
  "users.srv.pendapi": "Waiting API servers",
  "users.srv.known": "Servers of this installation",
  "users.srv.listed": "Entered servers",
  "users.srv.listed.hint": "These servers are in this installation's trust list.",
  "users.srv.paircode": "Pairing code",
  "users.srv.code.hint": "The code is only shown if it matches this installation.",
  "users.srv.joincode": "Join code",
  "users.srv.nowait": "No server is waiting at the moment.",
  "users.srv.take": "Take in",
  "users.srv.take.aria": "Take in server {name}",
  "users.srv.nowaitapi": "No API server is waiting at the moment.",
  "users.srv.nametaken": "This name is already taken in this installation. The server has to join under another name.",
  "users.nolist": "The trust list is not there yet.",
  "users.srv.none": "There is no server in the trust list yet.",
  "users.srv.kind.api": "API server",
  "users.srv.kind.ctl": "Server",
  "users.addr.change": "Change address",
  "users.addr.change.aria": "Change the address of {name}",
  "users.sys.change": "Change system rights",
  "users.sys.change.aria": "Change system rights on {name}",
  "users.retire": "Retire",
  "users.retire.aria": "Retire server {name}",
  "users.retire.lastapi": "The only API server of this installation cannot be retired. Take in another one first.",
  "users.addr.new.aria": "New address of {name}",
  "users.addr.title": "Change the address of “{name}”",
  "users.addr.body1": "This installation then reaches the server at the new address. The server's keys stay the same; only where it can be found changes.",
  "users.addr.body2": "The address looks like this: wss://server-name:7443/tunnel. An IPv4 address may stand instead of the name.",
  "users.addr.body3": "Only enter it if the server can be reached at exactly this address.",
  "users.addr.body4": "Important: also tell the server its new address yourself (on the server with ulh, entry cd-api) and restart it. Otherwise the web server turns it away because its certificate still names the old address.",
  "users.addr.same": "That is the address that is already entered. Nothing was changed.",
  "users.addr.bad": "That is not an address at which an API server can be reached. It has to look like this: wss://server-name:7443/tunnel",
  "users.addr.done": "The address of “{name}” is changed.",
  "users.addr.failed": "The address could not be changed.",
  "users.sys.title": "Change system rights on “{name}”",
  "users.sys.body1": "Whoever stands here may set up and remove apps, delete data and open the terminal on this server.",
  "users.sys.body2": "The server's keys stay the same. Only who may give it tasks with system rights changes.",
  "users.sys.same": "Nothing has changed there. Nothing was written.",
  "users.sys.done": "The system rights on “{name}” are changed.",
  "users.sys.failed": "The system rights could not be changed.",
  "users.retire.title": "Retire server “{name}”",
  "users.retire.body1": "The server is taken out of this installation's trust list. After that this installation no longer trusts it and makes no more connections to it.",
  "users.retire.body2": "The server itself keeps running; its services do not stop. To take it in again, it has to join anew.",
  "users.retire.go": "Retire server",
  "users.retire.onlyapi": "That is the only API server of this installation. It cannot be retired, because the installation could no longer be reached afterwards.",
  "users.retire.done": "The server “{name}” is retired.",
  "users.retire.failed": "The server could not be retired.",
  "users.col.server": "Server",
  "users.col.online": "Reachable",
  "users.col.tasks": "Tasks",
  "users.col.version": "Version",
  "users.col.arch": "Machine",
  "users.yes": "yes",
  "users.no": "no",
  "users.code.pair": "Show pairing code",
  "users.code.join": "Show join code",
  "users.code.nomatch": "This code does not match this installation.",
  "users.code.until": "Valid until {when}.",
  "users.cannot.server": "No server can be taken in right now.",
  "users.cannot.nolist": "The trust list is not there yet. {tail}",
  "users.cannot.nologin": "This needs a sign-in. {tail}",
  "users.join.title": "Take in server “{name}”",
  "users.join.body1": "The server is entered into this installation's trust list. After that it may run tasks of this installation.",
  "users.join.id": "In the list it is called: {id}",
  "users.join.only": "Only take it in if exactly this fingerprint stands on the server.",
  "users.join.tick": "Your box is ticked. Untick it if you yourself are not to set anything up on this server — then at first nobody can do so there. That can be changed later under “Entered servers”.",
  "users.join.go": "Take in server",
  "users.join.done.none": "The server “{name}” is taken in. Nobody may set up apps on this server yet.",
  "users.join.done.root": "The server “{name}” is taken in. Allowed to run tasks with system rights there: {list}.",
  "users.join.failed": "The server could not be taken in.",
  "users.joinapi.title": "Take in API server “{name}”",
  "users.joinapi.body": "The server is entered into this installation's trust list. After that this web interface connects to it and trusts it.",
  "users.joinapi.only": "Only take it in if exactly this fingerprint and exactly this address stand on the server.",
  "users.joinapi.done": "The server “{name}” is taken in.",
  "users.col.time": "Time",
  "users.col.user": "User",
  "users.col.kind": "Kind",
  "users.col.app": "App",
  "users.col.action": "Action",
  "users.col.result": "Result",
  "users.col.detail": "Details",
  "users.log.more": "Load more",
  "users.log.head": "Log",
  "users.log.empty": "The log is empty.",
  "users.log.start": "This is the start of the log.",
  "users.cancel": "Cancel",
  "users.nochange": "Nothing changed.",
  "users.unreadable": "A detail of this server cannot be read ({what}).",
  "users.addr.aria": "Address {addr}",
  "users.right.use": "use",
  // S12, shell: every English key of en-shell.json, word for word.
  "trust.alarm": "The trust list of this installation is not right. CirroDesk has signed you out to be safe.",
  "trust.notsetup": "This installation is not fully set up yet. Please pair the web server's controller.",
  "trust.nottrusted": "This server is not in the trust list.",
  "trust.alarm.head": "Something is wrong here",
  "trust.alarm.maybe": "Someone may have placed themselves between this device and your installation.",
  "trust.alarm.advice": "Please do not enter any more passwords here. Close this window and ask the person who set up CirroDesk for you.",
  "trust.alarm.foot": "This message cannot be dismissed, on purpose.",
  "frame.notgenuine": "This app's interface is not genuine. CirroDesk will not run it.",
  "frame.noentry": "This app's start page is missing. CirroDesk will not run it.",
  "frame.needsclassic": "This app needs classic scripts.",
  "frame.notloaded": "This app's interface could not be loaded.",
  "frame.hint.admin": "Please ask the administrator.",
  "frame.hint.retry": "Please close the window and try again.",
  "confirm.why.task": "This is an important step on this computer.",
  "confirm.why.list": "This standing permission will be entered in the trust list.",
  "confirm.head": "Please confirm",
  "confirm.line.passkey": "{why} Please confirm with your passkey.",
  "confirm.line.password": "{why} Please enter your password once more.",
  "confirm.password": "Password",
  "confirm.cancel": "Cancel",
  "confirm.yes.passkey": "Confirm with passkey",
  "confirm.yes": "Confirm",
  "boot.https": "CirroDesk needs a secure connection. Please open the page via https://.",
  "boot.notsetup": "CirroDesk is not set up.",
  "boot.loadfailed": "CirroDesk could not be loaded. Please reload the page. If that does not help, CirroDesk is not running at the moment.",
  "boot.passkeyhint": "Switching to passkeys recommended",
  "desk.unreachable": "{name} cannot be reached.",
  "win.tab": "Window {title}",
  "win.msgclose": "Close message",
  "stale.now": "Check now",
  "stale.now.title": "Fetch this state again now",
  "stale.fetching": "As of: fetching …",
  "stale.at": "As of: {ago}",
  "stale.justnow": "just now",
  "stale.min1": "1 minute ago",
  "stale.mins": "{n} minutes ago",
  "stale.hour1": "1 hour ago",
  "stale.hours": "{n} hours ago",
  "stale.day1": "1 day ago",
  "stale.days": "{n} days ago",
  // S12, net: every English key of en-net.json, word for word.
  "code.E_FORBIDDEN": "Permission denied.",
  "code.E_UNREACHABLE": "Service unavailable.",
  "code.E_EXIT": "The service ended with an error.",
  "code.E_TIMEOUT": "Time limit exceeded.",
  "code.E_LIMIT": "Too many open connections.",
  "code.E_PROTOCOL": "Invalid message.",
  "code.E_CONFIRM": "Please confirm again.",
  "code.E_LICENSE": "No valid licence.",
  "code.E_NOTFOUND": "Not found.",
  "code.E_VERIFY": "Invalid signature.",
  "code.E_INTERNAL": "Internal error.",
  "code.E_TOOLARGE": "Too much data.",
  "code.E_LOGIN": "Login failed.",
  "code.E_LOCKED": "Too many failed attempts. Please try again later.",
  "code.E_NOT_SET_UP": "CirroDesk is not set up yet.",
  "world.noaction": "CirroDesk does not know this request.",
  "world.wrongway": "This request cannot be made this way.",
  "world.nonode": "This app does not say which computer its service runs on.",
  "world.nosessionkey": "This sign-in cannot reach an app's own service. Please reload the page.",
  "world.notrustlist": "CirroDesk has not opened the trust list of this installation. Please reload the page.",
  "world.nosessiontask": "This sign-in cannot sign the task for this app. Please sign in again.",
  "world.nosignature": "The task for this app could not be signed. Please sign in again.",
  "world.notconfirmed": "Not confirmed. The app was not opened.",
  // S12, setup: every English key of en-setup.json, word for word.
  "setup.ui.head": "Set up CirroDesk",
  "setup.ui.code": "One-time code",
  "setup.ui.next": "Next",
  "setup.ui.compare": "Please compare: does this match what the server's screen shows? Then create your user account.",
  "setup.ui.site": "Installation",
  "setup.ui.id": "Identifier",
  "setup.ui.apis": "API servers",
  "setup.ui.node.name": "Name",
  "setup.ui.node.addr": "Address",
  "setup.ui.node.tls": "Fingerprint",
  "setup.ui.store": "App store",
  "setup.ui.store.use": "Allow installing apps from the CirroDesk store",
  "setup.ui.account": "Your user account",
  "setup.ui.account.lead": "This account may do everything. Remember the password well: without it, nobody can get into this installation any more.",
  "setup.ui.name": "User name",
  "setup.ui.password": "Password",
  "setup.ui.again": "Repeat password",
  "setup.ui.go": "Set up",
  // The appearance names, so theme.js keeps no German fallback of its own.
  "look.wallpaper.deep": "Space",
  "look.wallpaper.aurora": "Aurora",
  "look.wallpaper.sunset": "Sunset",
  "look.wallpaper.mint": "Mint",
  "look.wallpaper.plum": "Plum",
  "look.density.eng": "Compact",
  "look.density.normal": "Standard",
  "look.density.luftig": "Spacious",
  "look.glass.aus": "No glass",
  "look.glass.normal": "Glass",
  "look.glass.stark": "More glass",
});

// A map of strings, read through the house decoder with a budget: a language
// file is a file on a server, and a file on a server is exactly what somebody
// who got that far would edit. json.MAX_MEMBERS caps it at 64 keys, which is
// the growth this table is allowed before somebody has to think about it -
// MEASURED through this very shape rather than read off the constant: 64 keys
// decode, 65 leave as BadJSON "too many members".
const langShape = M(str);

let chosen = FALLBACK;
let table = Object.freeze({ __proto__: null });
let fallbackTable = Object.freeze({ __proto__: null });

// pick takes what the browser offers and answers with a name out of LANGS.
// It walks the offers in the browser's own order of preference and compares
// only the primary subtag: "de-CH" and "de" are the same file, and a region
// nobody has a file for is not a reason to fall back to English.
//
// It is exported because the test asks it directly: a function that is only
// ever reached through a fetch is a function measured through two things.
export function pick(offers) {
  const list = Array.isArray(offers) ? offers : [];
  for (const offer of list) {
    if (typeof offer !== "string") continue;
    const primary = offer.toLowerCase().split("-")[0];
    if (LANGS.includes(primary)) return primary;
  }
  return FALLBACK;
}

// The files of one language: the base file and one per area of the shell.
// WHY MORE THAN ONE: a file is one JSON object, and json.MAX_MEMBERS (64, a
// contract of shared/) caps an object's members. The base file holds the login
// and the first run; each area that is translated brings its own file
// (S12). A part that is missing is no error - its keys fall back per key like
// any other.
export const PARTS = Object.freeze(["", "desk", "store", "store2", "users", "users2", "users3", "users4", "shell", "net", "setup"]);

async function fetchLang(name) {
  const out = { __proto__: null };
  let any = false;
  for (const part of PARTS) {
    const got = await fetchPart(part === "" ? name : `${name}-${part}`);
    if (got === null) continue;
    any = true;
    Object.assign(out, got);
  }
  return any ? Object.freeze(out) : null;
}

async function fetchPart(name) {
  let res;
  try {
    res = await fetch(`/assets/lang/${name}.json`, { cache: "no-store" });
  } catch {
    return null;
  }
  if (!res.ok) return null;
  const raw = await readBounded(res, MAX_LANG);
  if (raw === null) return null;
  try {
    // MEASURED, not assumed: decode hands M(str) back as a PLAIN OBJECT and
    // not as a Map. The first version of this read it with
    // "for (const [k, v] of got)", which throws on an object, was caught by
    // the catch below and left every sentence empty - four cases red, and the
    // one that fetched the file itself green beside them, which is what made
    // it look like a loader that worked.
    //
    // It is copied into an object with NO PROTOTYPE, so t() for a key from
    // anywhere answers out of the file and never out of Object.prototype -
    // the same reason contracts/wire.js keeps its own table that way.
    const got = decode(raw, new Budget({ size: MAX_LANG }), langShape);
    const out = { __proto__: null };
    for (const [k, v] of Object.entries(got)) out[k] = v;
    return Object.freeze(out);
  } catch {
    return null;
  }
}

// load reads the language of this browser and English beside it, and says
// which language it settled on. It is called once, before the first mask is
// drawn; a second call simply reads again.
//
// A file that cannot be read is NOT a failure of the page: the shell goes on
// with what it has, and t() answers "" for a key it then does not hold. A
// login mask that refuses to appear because a sentence could not be fetched
// would be worse than a login mask with one empty line in it.
// ready settles once the first load() is done, file or no file. The shell
// awaits load() itself before it draws; a page that draws an app on its own
// (a test) waits here, so it never draws before the language is in.
let settleReady;
export const ready = new Promise((r) => { settleReady = r; });

export async function load(offers = navigator.languages ?? [navigator.language]) {
  chosen = pick(offers);
  fallbackTable = (await fetchLang(FALLBACK)) ?? Object.freeze({ __proto__: null });
  table = chosen === FALLBACK ? fallbackTable : (await fetchLang(chosen)) ?? fallbackTable;
  document.documentElement.lang = chosen;
  settleReady(chosen);
  return chosen;
}

export function language() {
  return chosen;
}

// t answers the sentence for a key, with {name} filled in from vars. The
// values are put in as they are and the result goes to textContent, so there
// is nothing here to escape; a placeholder nobody passed a value for stays
// as it is written, which shows up in a test rather than disappearing.
export function t(key, vars = null) {
  const found = typeof key === "string"
    ? (table[key] ?? fallbackTable[key] ?? BUILT_IN[key] ?? "")
    : "";
  if (vars === null || found === "") return found;
  let out = found;
  for (const [name, value] of Object.entries(vars)) {
    out = out.split(`{${name}}`).join(String(value));
  }
  return out;
}
