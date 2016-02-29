package migrations

// Even though sqlite has no notion of truncated text columns,
// let's provide them to provide clarity around expectations and
// guidance for when other database support is added.
var create_tables string = `
	CREATE TABLE IF NOT EXISTS camera (
		id INTEGER PRIMARY KEY,
		name VARCHAR(64),
		url TEXT,
		username VARCHAR(64),
		password VARCHAR(64),
		threshold INTEGER,
		min_change INTEGER
	);

	CREATE UNIQUE INDEX IF NOT EXISTS uniq_cam_name ON camera (name);

	CREATE TABLE IF NOT EXISTS event (
		id INTEGER PRIMARY KEY,
		camera_id INTEGER,
		type VARCHAR(32),
		start_time DATETIME,
		end_time DATETIME,
		duration INTEGER,
		filepath VARCHAR(1024),
		first_frame VARCHAR(1024),
		last_frame VARCHAR(1024),
		frame_count INTEGER
	);

	CREATE INDEX IF NOT EXISTS cam_start_end_duration ON event(camera_id, start_time, end_time, duration);
	CREATE INDEX IF NOT EXISTS filepath_idx ON event(filepath);

	CREATE TABLE IF NOT EXISTS setting (
		id INTEGER PRIMARY KEY,
		key VARCHAR(255),
		value VARCHAR(1024)
	);

	CREATE UNIQUE INDEX IF NOT EXISTS key_value ON setting(key, value);

	CREATE TABLE IF NOT EXISTS camuser (
		id INTEGER PRIMARY KEY,
		name VARCHAR(64),
		password_hash VARCHAR(255)
	);

	CREATE UNIQUE INDEX IF NOT EXISTS user_name ON user(name);
`

func GetSQLiteMigrations() []string {
	return []string{
		create_tables,
	}
}
