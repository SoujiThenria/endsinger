package discord

import (
	"errors"
	"fmt"

	log "github.com/SoujiThenria/endsinger/logging"
	"github.com/bwmarrin/discordgo"
)

// Return the discordgo API error in a struct, so that the
// messgae and error code are seperate values
func simpleAPIError(em error) (response *discordgo.APIErrorMessage) {
	var restError *discordgo.RESTError
	// default response
	response = &discordgo.APIErrorMessage{
		Code:    0,
		Message: "There was nothing in this API error",
	}

	if errors.As(em, &restError) && restError.Message != nil {
		response = restError.Message
	}
	return
}

// Return the discordgo API error as a singe string
// in the format: string [code]
func oneStrAPIError(em error) (str string) {
	simpleResp := simpleAPIError(em)
	str = fmt.Sprintf("%s [%d]", simpleResp.Message, simpleResp.Code)
	return
}

// Function to replay to application commands
func applicationCommandReply(s *discordgo.Session, i *discordgo.InteractionCreate, responseMessage string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: responseMessage,
		},
	})
	if err != nil {
		log.Error("Cannot reply to the application command: %s - Error: %s", i.ApplicationCommandData().Name, err)
	}
}
