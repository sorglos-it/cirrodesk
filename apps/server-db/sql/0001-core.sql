-- 0001 creates the whole schema of CirroDesk. Every statement is one step: the
-- runner records each one, so a crash in the middle of a migration picks up
-- where it stopped. Only portable SQL, nothing but what MariaDB, SQLite,
-- MySQL, MSSQL and Oracle all understand: no engine keyword, no IF NOT
-- EXISTS, no AUTO_INCREMENT, no string literal, no database function. The
-- {{...}} are column types; the runner fills them from the fixed table of the
-- dialect (classes/sqlstore/names.go).
--
-- Every column is binary or a number, and never NULL. Binary means that both
-- databases sort and compare exactly alike, the way Go does: a log ID is
-- ordered byte for byte, a name is never matched in a different case, and no
-- collation can change any of that later.

CREATE TABLE cd_blob (
	kind {{key}} NOT NULL,
	id {{key}} NOT NULL,
	seq {{int}} NOT NULL,
	data {{data}} NOT NULL,
	updated {{int}} NOT NULL,
	PRIMARY KEY (kind, id)
){{opts}};

CREATE INDEX cd_blob_seq ON cd_blob (kind, seq, id);

CREATE TABLE cd_app (
	id {{key}} NOT NULL,
	version {{key}} NOT NULL,
	store {{key}} NOT NULL,
	node {{key}} NOT NULL,
	active {{int}} NOT NULL,
	manifest {{doc}} NOT NULL,
	sha256 {{key}} NOT NULL,
	grants {{doc}} NOT NULL,
	kept_data {{int}} NOT NULL,
	updated {{int}} NOT NULL,
	PRIMARY KEY (id)
){{opts}};

CREATE TABLE cd_once (
	kind {{key}} NOT NULL,
	token {{key}} NOT NULL,
	data {{key}} NOT NULL,
	expires {{int}} NOT NULL,
	used {{int}} NOT NULL,
	PRIMARY KEY (kind, token)
){{opts}};

CREATE INDEX cd_once_expires ON cd_once (expires);

CREATE TABLE cd_user (
	id {{key}} NOT NULL,
	name {{key}} NOT NULL,
	roles {{doc}} NOT NULL,
	totp {{doc}} NOT NULL,
	locked {{int}} NOT NULL,
	failures {{int}} NOT NULL,
	last_login {{int}} NOT NULL,
	PRIMARY KEY (id)
){{opts}};

CREATE UNIQUE INDEX cd_user_name ON cd_user (name);

CREATE TABLE cd_session (
	token {{key}} NOT NULL,
	user_id {{key}} NOT NULL,
	cert {{doc}} NOT NULL,
	proof {{doc}} NOT NULL,
	api_key {{doc}} NOT NULL,
	api_cert {{doc}} NOT NULL,
	created {{int}} NOT NULL,
	expires {{int}} NOT NULL,
	seen {{int}} NOT NULL,
	PRIMARY KEY (token)
){{opts}};

CREATE INDEX cd_session_expires ON cd_session (expires);

CREATE TABLE cd_pending (
	pubkey {{key}} NOT NULL,
	name {{key}} NOT NULL,
	since {{int}} NOT NULL,
	PRIMARY KEY (pubkey)
){{opts}};

CREATE INDEX cd_pending_since ON cd_pending (since, name);

CREATE TABLE cd_log (
	id {{key}} NOT NULL,
	ts {{int}} NOT NULL,
	user_id {{key}} NOT NULL,
	kind {{key}} NOT NULL,
	app {{key}} NOT NULL,
	action {{key}} NOT NULL,
	node {{key}} NOT NULL,
	result_code {{key}} NOT NULL,
	detail {{text}} NOT NULL,
	PRIMARY KEY (id)
){{opts}};

CREATE INDEX cd_log_ts ON cd_log (ts, id);

CREATE INDEX cd_log_user ON cd_log (user_id, ts, id);
