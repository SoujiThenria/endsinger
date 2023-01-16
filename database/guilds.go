package database

type Guild struct {
	Snowflake string
	Name      string
}

func (g *Guild) add() error {
	statement, err := db.Prepare(guildAdd)
	if err != nil {
		return err
	}
	_, err = statement.Exec(g.Snowflake, g.Name)
	return err
}

func (g *Guild) update() error {
	statement, err := db.Prepare(guildUpdate)
	if err != nil {
		return err
	}
	_, err = statement.Exec(g.Name, g.Snowflake)
	return err
}
