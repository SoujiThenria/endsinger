package discord

import (
	"github.com/SoujiThenria/endsinger/database"
	log "github.com/SoujiThenria/endsinger/logging"
	"github.com/bwmarrin/discordgo"
)

func initHandlers() {
	session.AddHandler(onReady)
	session.AddHandler(onCommand)

	session.AddHandler(channelUpdate)
	session.AddHandler(channelDelete)

	session.AddHandler(guildUpdate)
	session.AddHandler(guildDelete)
}

// Print the name and discriminator of the bot user
func onReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Info("Connected to discord successfully as %s#%s", r.User.Username, r.User.Discriminator)
}

// Application command handler
func onCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}
}

// When a channel gets modified, only changes in the name are interesting
func channelUpdate(s *discordgo.Session, c *discordgo.ChannelUpdate) {
	channel, err := database.ChannelGet(c.ID)
	// Channel does not exists or name did not changed
	if err == database.ErrorNowRows || c.Name == channel.Name {
		return
	}
	// If a "real" error occured.
	if err != nil {
		log.Error("Cannot update the channel name in the database to: %s [%s] - Error: %s", c.Name, c.ID, err)
		return
	}

	// Update the channel name
	channel.Name = c.Name
	err = database.Update(channel)
	if err != nil {
		log.Error("Cannot update the channel name in the database to: %s [%s] - Error: %s", channel.Name, channel.Snowflake, err)
		return
	}
}

// Delete a channel from the database, if the channel was deleted in the guild
func channelDelete(s *discordgo.Session, c *discordgo.ChannelDelete) {
	channel, err := database.ChannelGet(c.ID)
	if err == database.ErrorNowRows {
		return
	}
	// If a "real" error occured.
	if err != nil {
		log.Error("Cannot remove the channel from the database: %s [%s] - Error: %s", c.Name, c.ID, err)
		return
	}
	if err := database.ChannelDelete(channel.Snowflake); err != nil {
		log.Error("Cannot remove the channel from the database: %s [%s] - Error: %s", channel.Name, channel.Snowflake, err)
		return
	}
}

// Update the guild name if the name was chaned
func guildUpdate(s *discordgo.Session, g *discordgo.GuildUpdate) {
	guild, err := database.GuildGet(g.ID)
	// Guild does not exists or the name did not changed
	if err == database.ErrorNowRows || guild.Name == g.Name {
		return
	}

	// If a "real" error occured.
	if err != nil {
		log.Error("Cannot update the guild name in the database to: %s [%s] - Error: %s", g.Name, g.ID, err)
		return
	}

	guild.Name = g.Name
	err = database.Update(guild)
	if err != nil {
		log.Error("Cannot update the guild name in the database to: %s [%s] - Error: %s", guild.Name, guild.Snowflake, err)
		return
	}
}

// Dete a guild from the database
func guildDelete(s *discordgo.Session, g *discordgo.GuildDelete) {
	guild, err := database.GuildGet(g.ID)
	if err == database.ErrorNowRows {
		return
	}

	// If a "real" error occured.
	if err != nil {
		log.Error("Cannot remove the guild from the database: %s [%s] - Error: %s", g.Name, g.ID, err)
		return
	}
	if err := database.GuildDelete(guild.Snowflake); err != nil {
		log.Error("Cannot remove the guild from the database: %s [%s] - Error: %s", guild.Name, guild.Snowflake, err)
		return
	}
}
