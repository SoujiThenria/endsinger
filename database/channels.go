package database

type Channel struct {
	Snowflake string
	Name      string
	Active    bool
	Days      int8
	Guild     *Guild
}

func (c *Channel) add() error {
	statement, err := db.Prepare(channelAdd)
	if err != nil {
		return err
	}
	_, err = statement.Exec(c.Snowflake, c.Name, c.Days, c.Guild.Snowflake)
	return err
}

func (c *Channel) update() error {
	statement, err := db.Prepare(channelUpdate)
	if err != nil {
		return err
	}
	_, err = statement.Exec(c.Name, c.Active, c.Days, c.Snowflake)
	return nil
}
