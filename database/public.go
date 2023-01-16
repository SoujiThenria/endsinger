package database

import "database/sql"

type Adder interface {
	add() error
}

type Updater interface {
	update() error
}

var (
	ErrorNowRows = sql.ErrNoRows
)

// Add a channel or a guild to the database.
func Add(a Adder) (err error) {
	err = a.add()
	return
}

// Update a channel or a guild in the database.
func Update(u Updater) (err error) {
	err = u.update()
	return
}

// Delete a guild in the database and all associated channels.
func GuildDelete(snowflake string) (err error) {
	// Remove channels.
	statement, err := db.Prepare(chanenlDeleteWhereGuild)
	if err != nil {
		return
	}
	_, err = statement.Exec(snowflake)
	if err != nil {
		return
	}

	// Remove the guild.
	statement, err = db.Prepare(guildDelete)
	if err != nil {
		return
	}
	_, err = statement.Exec(snowflake)
	return
}

// Get a guild from the database.
func GuildGet(snowflake string) (guild *Guild, err error) {
	guild = &Guild{
		Snowflake: snowflake,
	}
	row := db.QueryRow(guildGet, snowflake)
	err = row.Scan(&guild.Name)
	return
}

// Delete a channel from the database.
func ChannelDelete(snowflake string) (err error) {
	statement, err := db.Prepare(channelDelete)
	if err != nil {
		return
	}
	_, err = statement.Exec(snowflake)
	return
}

// Get a channel from the database.
func ChannelGet(snowflake string) (channel *Channel, err error) {
	channel = &Channel{
		Snowflake: snowflake,
	}
	row := db.QueryRow(channelGet, snowflake)
	err = row.Scan(&channel.Name, &channel.Active, &channel.Days)
	return
}

// Returns a slice with all active channels for a specific guild.
func ChannelListActive(guildID string) (channels []*Channel, err error) {
	channels = []*Channel{}
	rows, err := db.Query(channelGetActive, guildID)
	if err != nil {
		return
	}
	for rows.Next() {
		c := &Channel{
			Active: true,
		}
		err = rows.Scan(&c.Name, &c.Snowflake, &c.Days)
		if err != nil {
			return
		}
		channels = append(channels, c)
	}
	return
}

// Returns a slice with all active channels.
func ChannelListActiveAll() (channels []*Channel, err error) {
	channels = []*Channel{}
	rows, err := db.Query(channelGetActiveAll)
	if err != nil {
		return
	}
	for rows.Next() {
		c := &Channel{
			Active: true,
		}
		err = rows.Scan(&c.Name, &c.Snowflake, &c.Days)
		if err != nil {
			return
		}
		channels = append(channels, c)
	}
	return
}
