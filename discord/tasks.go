package discord

import (
	"time"

	"github.com/SoujiThenria/endsinger/database"
	log "github.com/SoujiThenria/endsinger/logging"
	"github.com/bwmarrin/discordgo"
)

// The stop channel can cause the goroutine to stop, but do not have to.
// But honestly it does not ready matter if the go routine just stops when the programm is terminated.
func channelCleanup(stop chan struct{}) {
	log.Debug("Goroutin channelCleanup started.")
	tNow := time.Now()
	tExe := time.Date(tNow.Year(), tNow.Month(), tNow.Day(), 23, 50, 0, 0, tNow.Location())
	tWait := tExe.Sub(tNow)
	if tWait < 0 {
		tWait = tExe.Sub(tExe.Add(24 * time.Hour))
	}
	timer := time.NewTimer(tWait)

	// Wait for timer to finish or to receive shutdown signal
	select {
	case <-timer.C:
	case <-stop:
		return
	}

	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			ticker.Reset(24 * time.Hour)
			channels, err := database.ChannelListActiveAll()
			if err != nil {
				log.Error("Cannot get channels: %s", err)
				break
			}
			log.Debug("channelCleanup: Number of active channels: %d", len(channels))

			for _, v := range channels {
				log.Debug("Channel: %s", v.Name)
				messages, err := getMessages(v)
				if err != nil {
					log.Error("Cannot receive messages: %s [%s] - Error: %s", v.Name, v.Snowflake, oneStrAPIError(err))
					// Skip this iteration
					continue
				}

				log.Debug("%s: Message count: %d", v.Name, len(messages))
				if len(messages) > 0 {
					log.Info("Messages to delete: Channel: %s - %d", v.Name, len(messages))
				}
				if err := deleteMessages(v, messages); err != nil {
					log.Error("Cannot delete messages: %s [%s] - Error: %s", v.Name, v.Snowflake, err)
				}
			}
		case <-stop:
			return
		}
	}

}

func getMessages(c *database.Channel) (messageIDs []string, err error) {
	t := time.Now().Add((time.Duration(c.Days) * (-24)) * time.Hour)
	lastMessageID := ""
	for {
		msg, err := session.ChannelMessages(c.Snowflake, 100, lastMessageID, "", "")
		if err != nil {
			return nil, err
		}

		for _, m := range msg {
			snowTime, err := discordgo.SnowflakeTimestamp(m.ID)
			if err != nil {
				return nil, err
			}
			if t.After(snowTime) {
				messageIDs = append(messageIDs, m.ID)
			}
		}
		if len(msg) < 100 {
			return messageIDs, err
		}
		lastMessageID = msg[len(msg)-1].ID
	}
}

func deleteMessages(c *database.Channel, messages []string) (err error) {
	messageCount := len(messages)

	for i := 0; i < messageCount; i += 100 {
		var messageSlice []string
		if (i + 100) < messageCount {
			messageSlice = messages[i : i+100]
		} else {
			messageSlice = messages[i:]
		}
		log.Debug("%s: Bulk delete messages: %d", c.Name, len(messageSlice))
		err := session.ChannelMessagesBulkDelete(c.Snowflake, messageSlice)
		if err != nil {
			log.Debug("Error while executing bulk delete: %s", err)
			for _, m := range messageSlice {
				if err := session.ChannelMessageDelete(c.Snowflake, m); err != nil {
					return err
				}
			}
		}
	}
	return
}
