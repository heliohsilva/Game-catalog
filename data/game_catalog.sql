PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE games (
		id TEXT PRIMARY KEY NOT NULL,
		title TEXT NOT NULL,
		platform TEXT NOT NULL,
		subcategory TEXT,
		genre TEXT NOT NULL DEFAULT 'Gaming',
		time_to_beat TEXT,
		time_to_beat_main REAL,
		added_at TEXT NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	);
INSERT INTO games VALUES('game-8f3d26e5-1859-401b-93dd-5378f331b855','Steep','PC','Steam','Action','12h',12.0,'2026-10-04T19:37:58Z','2026-10-04T19:37:58Z','2026-10-04T19:37:58Z');
INSERT INTO games VALUES('game-ae5b6c72-e8ce-46d8-aff0-e3a6fa50e641','Hades II','Nintendo Switch','Switch','RogueLike','32h',32.0,'2026-10-04T19:38:14Z','2026-10-04T19:38:14Z','2026-10-04T19:38:14Z');
CREATE TABLE platforms (
		id TEXT PRIMARY KEY NOT NULL,
		name TEXT UNIQUE NOT NULL COLLATE NOCASE,
		created_at TEXT NOT NULL
	);
INSERT INTO platforms VALUES('plat-default-1','PC','2026-10-04T18:44:59Z');
INSERT INTO platforms VALUES('plat-default-2','PlayStation','2026-10-04T18:44:59Z');
INSERT INTO platforms VALUES('plat-default-3','Nintendo Switch','2026-10-04T18:44:59Z');
INSERT INTO platforms VALUES('plat-default-4','Xbox','2026-10-04T18:44:59Z');
INSERT INTO platforms VALUES('plat-default-5','Retro / Emulation','2026-10-04T18:44:59Z');
INSERT INTO platforms VALUES('plat-28220a31-fa17-4e87-b7a6-8f4736df4706','DS','2026-10-04T19:37:35Z');
CREATE TABLE platform_subcategories (
		id TEXT PRIMARY KEY NOT NULL,
		platform_name TEXT NOT NULL COLLATE NOCASE,
		name TEXT NOT NULL COLLATE NOCASE,
		created_at TEXT NOT NULL,
		UNIQUE(platform_name, name),
		FOREIGN KEY (platform_name) REFERENCES platforms(name) ON DELETE CASCADE ON UPDATE CASCADE
	);
INSERT INTO platform_subcategories VALUES('sub-default-1-1','PC','Steam','2026-10-04T18:44:59Z');
INSERT INTO platform_subcategories VALUES('sub-default-1-2','PC','GOG','2026-10-04T18:44:59Z');
INSERT INTO platform_subcategories VALUES('sub-default-1-3','PC','Epic','2026-10-04T18:44:59Z');
INSERT INTO platform_subcategories VALUES('sub-default-2-1','PlayStation','PS2','2026-10-04T18:44:59Z');
INSERT INTO platform_subcategories VALUES('sub-default-2-2','PlayStation','PS3','2026-10-04T18:44:59Z');
INSERT INTO platform_subcategories VALUES('sub-default-2-3','PlayStation','PS4','2026-10-04T18:44:59Z');
INSERT INTO platform_subcategories VALUES('sub-default-2-4','PlayStation','PS5','2026-10-04T18:44:59Z');
INSERT INTO platform_subcategories VALUES('sub-default-3-1','Nintendo Switch','Switch','2026-10-04T18:44:59Z');
INSERT INTO platform_subcategories VALUES('sub-default-4-1','Xbox','Xbox Series X/S','2026-10-04T18:44:59Z');
INSERT INTO platform_subcategories VALUES('sub-default-4-2','Xbox','Xbox One','2026-10-04T18:44:59Z');
INSERT INTO platform_subcategories VALUES('sub-default-4-3','Xbox','Xbox 360','2026-10-04T18:44:59Z');
INSERT INTO platform_subcategories VALUES('sub-default-5-1','Retro / Emulation','PS1','2026-10-04T18:44:59Z');
INSERT INTO platform_subcategories VALUES('sub-default-5-2','Retro / Emulation','N64','2026-10-04T18:44:59Z');
INSERT INTO platform_subcategories VALUES('sub-default-5-3','Retro / Emulation','SNES','2026-10-04T18:44:59Z');
INSERT INTO platform_subcategories VALUES('sub-default-5-4','Retro / Emulation','Genesis','2026-10-04T18:44:59Z');
INSERT INTO platform_subcategories VALUES('sub-default-5-5','Retro / Emulation','NES','2026-10-04T18:44:59Z');
INSERT INTO platform_subcategories VALUES('sub-default-5-6','Retro / Emulation','Master System','2026-10-04T18:44:59Z');
INSERT INTO platform_subcategories VALUES('sub-default-5-7','Retro / Emulation','Neo Geo','2026-10-04T18:44:59Z');
INSERT INTO platform_subcategories VALUES('sub-18af165a-07cf-4dc3-8e13-c24dea8b57a5','PC','Ubisoft Connect','2026-10-04T18:45:30Z');
CREATE INDEX idx_games_platform ON games (platform);
CREATE INDEX idx_games_subcategory ON games (subcategory);
CREATE INDEX idx_games_genre ON games (genre);
CREATE INDEX idx_games_added_at ON games (added_at);
CREATE INDEX idx_games_title ON games (title);
CREATE UNIQUE INDEX idx_games_unique_entry 
	ON games (LOWER(TRIM(title)), platform, COALESCE(subcategory, ''));
CREATE INDEX idx_platform_subcategories_platform ON platform_subcategories (platform_name);
COMMIT;
