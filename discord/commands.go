package discord

import (
	"bytes"
	"sort"
	"strconv"

	"github.com/SoujiThenria/endsinger/database"
	log "github.com/SoujiThenria/endsinger/logging"
	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
)

var (
	adminPermission int64   = 8
	minDays         float64 = 1
	// The maximum is 7 days because that's the
	// maximum age of a message to be valid for bulk operations.
	maxDays float64 = 7

	commands = []*discordgo.ApplicationCommand{
		{
			Name:                     "add",
			Description:              "Add a channel to the database.",
			DefaultMemberPermissions: &adminPermission,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "The channel to add to the database.",
					Required:    false,
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "days",
					Description: "Message age in days before it will be deleted.",
					Required:    false,
					MaxValue:    maxDays,
					MinValue:    &minDays,
				},
			},
		},
		{
			Name:                     "remove",
			Description:              "Remove a channel from the database.",
			DefaultMemberPermissions: &adminPermission,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "The channel to remove from the database.",
					Required:    false,
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,
					},
				},
			},
		},
		{
			Name:                     "update",
			Description:              "Update the days for a channel in the database.",
			DefaultMemberPermissions: &adminPermission,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "days",
					Description: "Message age in days before it will be deleted.",
					Required:    true,
					MaxValue:    maxDays,
					MinValue:    &minDays,
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "The channel to update.",
					Required:    false,
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,
					},
				},
			},
		},
		{
			Name:                     "list",
			Description:              "List all channels in the database.",
			DefaultMemberPermissions: &adminPermission,
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"add": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			applicationCommandReply(s, i, commandAdd(i))
		},
		"remove": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			applicationCommandReply(s, i, commandRemove(i))
		},
		"list": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			applicationCommandReply(s, i, commandList(i))
		},
		"update": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			applicationCommandReply(s, i, commandUpdate(i))
		},
	}
)

// Register application commands
func registerCommands() {
	for i, v := range commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, "", v)
		if err != nil {
			log.Error("Cannot create application command: %s - Error: %s", v.Name, oneStrAPIError(err))
		} else {
			commands[i].ID = cmd.ID
		}
	}
}

// Remove application commands
func removeCommands() {
	for _, v := range commands {
		if v.ID != "" {
			err := session.ApplicationCommandDelete(session.State.User.ID, "", v.ID)
			if err != nil {
				log.Error("Cannot delete application command: %s - Error: %s", v.Name, oneStrAPIError(err))
			}
		}
	}
}

// Handle the /add command.
func commandAdd(i *discordgo.InteractionCreate) (msg string) {
	// Default response message
	msg = "Something went wrong, and I have no clue what..."

	var dChannelID string = i.ChannelID
	var Days int8 = 5

	// Create option map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(i.ApplicationCommandData().Options))
	for _, opt := range i.ApplicationCommandData().Options {
		optionMap[opt.Name] = opt
	}
	if option, ok := optionMap["channel"]; ok {
		dChannelID = option.ChannelValue(nil).ID
	}
	if option, ok := optionMap["days"]; ok {
		Days = int8(option.IntValue())
	}

	// Search the channel in the database.
	channel, err := database.ChannelGet(dChannelID)
	// The channel exists, and so must the guild.
	if err == nil {
		if channel.Active {
			msg = "This channel is already in the database."
			return
		}
		channel.Active = true
		channel.Days = Days
		err = database.Update(channel)
		if err != nil {
			log.Error("Cannot update (activate) the channel: %s [%s] - Error: %s", channel.Name, channel.Snowflake, err)
			msg = "Failed to add the channel to the database."
			return
		}
		msg = "Channel successfully added to the database."
		return
	}

	// The channel does not exists. Test if the guild exists.
	guild, err := database.GuildGet(i.GuildID)
	// The guild does not exists -> add it
	if err == database.ErrorNowRows {
		dg, err := session.Guild(i.GuildID)
		if err != nil {
			log.Error("Cannot resolve a guild data: %s", err)
			msg = "Failed to add the channel to the database. Please try it again later."
			return
		}
		guild.Snowflake = dg.ID
		guild.Name = dg.Name
		err = database.Add(guild)
		if err != nil {
			log.Error("Cannot add a guild to the database: %s [%s] - Error: %s", guild.Name, guild.Snowflake, err)
			msg = "Failed to add the channel to the database."
			return
		}
	} else if err != nil {
		msg = "Something went wrong. Please see the logs."
		log.Error("Cannot get guild: %s", err)
		return
	}

	// Resolve the channel ID to all channel-related data.
	dChannel, err := session.Channel(dChannelID)
	if err != nil {
		log.Error("Cannot resolve channel data: %s", oneStrAPIError(err))
		msg = "Failed to add the channel to the database: " + oneStrAPIError(err)
		return
	}
	// Add the channel to the database
	err = database.Add(&database.Channel{
		Snowflake: dChannel.ID,
		Name:      dChannel.Name,
		Active:    true,
		Days:      Days,
		Guild:     guild,
	})
	if err != nil {
		log.Error("Cannot add a channel to the database: %s [%s] - Error: %s", channel.Name, channel.Snowflake, err)
		msg = "Failed to add the channel to the database."
		return
	}

	msg = "Channel successfully added to the database."
	return
}

// Handle the /remove command.
func commandRemove(i *discordgo.InteractionCreate) (msg string) {
	// The dafault response message.
	msg = "Something went wrong, I don't know what..."

	// This variable changes if a channel was as command-option specified.
	var snowflake string = i.ChannelID

	// change the snowflake if the /remove command was called with a channel option.
	if i.ApplicationCommandData().Options != nil {
		channel := i.ApplicationCommandData().Options[0].ChannelValue(session)
		snowflake = channel.ID
	}

	channel, err := database.ChannelGet(snowflake)
	if err == database.ErrorNowRows || !channel.Active {
		msg = "This channel isn't in the database."
		return
	}
	if err != nil {
		log.Error("Cannot get channel from the database: %s", err)
		return
	}

	if channel.Active {
		channel.Active = false
		err = database.Update(channel)
		if err != nil {
			log.Error("Cannot update (activate) the channel: %s [%s] - Error: %s", channel.Name, channel.Snowflake, err)
			msg = "Failed to remove the channel from the database."
			return
		}
		msg = "Channel was successfully removed from the database."
	}
	return
}

// Handle the /update command.
func commandUpdate(i *discordgo.InteractionCreate) (msg string) {
	// The default response message.
	msg = "Something went wrong, I don't know what..."

	var channelID string = i.ChannelID
	var days int8 = 5

	// Create option map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(i.ApplicationCommandData().Options))
	for _, opt := range i.ApplicationCommandData().Options {
		optionMap[opt.Name] = opt
	}
	if option, ok := optionMap["channel"]; ok {
		channelID = option.ChannelValue(nil).ID
	}
	if option, ok := optionMap["days"]; ok {
		days = int8(option.IntValue())
	}

	channel, err := database.ChannelGet(channelID)
	if err == database.ErrorNowRows || !channel.Active {
		msg = "This channel isn't in the database."
		return
	}
	if err != nil {
		log.Error("Cannot get channel from the database: %s", err)
		return
	}

	channel.Days = days
	err = database.Update(channel)
	if err != nil {
		log.Error("Cannot update the channel: %s [%s] - Error: %s", channel.Name, channel.Snowflake, err)
		msg = "Failed to update the channel."
		return
	}

	msg = "The channel was successfully updated."
	return
}

// Handle the /list command.
func commandList(i *discordgo.InteractionCreate) (msg string) {
	// The dafault response message.
	msg = "Something went wrong, I don't know what..."

	// Get all active channels for the guild the command was executed from.
	channels, err := database.ChannelListActive(i.GuildID)
	if err != nil {
		log.Error("Failed to fetch all active marked channels for the guild [%s] - Error: %s", i.GuildID, err)
		msg = "Cannot list channels. There was an error. For more information, please see the logs."
		return
	}

	// There where no active channels.
	if len(channels) < 1 {
		msg = "There are no channels currently in the database."
		return
	}

	sort.SliceStable(channels, func(i, j int) bool {
		return channels[i].Name < channels[j].Name
	})

	buf := new(bytes.Buffer)
	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"# Channel", "Snowflake", "Days"})
	for _, v := range channels {
		table.Append([]string{"# " + v.Name, v.Snowflake, strconv.Itoa(int(v.Days))})
	}
	table.Render()
	return "```\n" + buf.String() + "\n```"
}
