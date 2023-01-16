package database

// Create tables
const (
	tableCreateGuild = `CREATE TABLE IF NOT EXISTS guilds (
		"ID" 		INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"Snowflake" 	TEXT NOT NULL UNIQUE,
		"Name" 		TEXT NOT NULL
	);`
	tableCreateChannel = `CREATE TABLE IF NOT EXISTS channels (
		"ID" 		INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"Snowflake" 	TEXT NOT NULL UNIQUE,
		"Name" 		TEXT NOT NULL,
		"Days" 		INTEGER NOT NULL DEFAULT 5,
		"Active" 	BOOLEAN DEFAULT TRUE,
		"GuildID" 	INTEGER,

		CONSTRAINT fk_guilds
			FOREIGN KEY (GuildID)
			REFERENCES guilds(ID)
	);`
)

// All channel-related queries.
const (
	channelAdd = `INSERT INTO channels(Snowflake, Name, Days, GuildID) VALUES(?, ?, ?,
 		(SELECT ID FROM guilds WHERE Snowflake = ?)
 	);`
	channelDeactivate       = `UPDATE channels SET Active = FALSE WHERE Snowflake = ?;`
	channelActivate         = `UPDATE channels SET Active = TRUE WHERE Snowflake = ?`
	channelDelete           = `DELETE FROM channels WHERE Snowflake = ?;`
	channelGet              = `SELECT Name, Active, Days FROM channels WHERE Snowflake = ?;`
	channelGetActive        = `SELECT Name, Snowflake, Days FROM channels WHERE Active = TRUE AND GuildID = (SELECT ID FROM guilds WHERE Snowflake = ?);`
	channelGetActiveAll     = `SELECT Name, Snowflake, Days FROM channels WHERE Active = TRUE;`
	channelUpdate           = `UPDATE channels SET Name = ?, Active = ?, Days = ? WHERE Snowflake = ?;`
	chanenlDeleteWhereGuild = `DELETE FROM channels WHERE GuildID = (SELECT ID FROM guilds WHERE Snowflake = ?);`
)

// All guild-related queries.
const (
	guildAdd    = `INSERT INTO guilds(Snowflake, Name) VALUES(?, ?);`
	guildDelete = `DELETE FROM guilds WHERE Snowflake = ?;`
	guildGet    = `SELECT Name FROM guilds WHERE Snowflake = ?;`
	guildUpdate = `UPDATE guilds SET Name = ? WHERE Snowflake = ?;`
)
