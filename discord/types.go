package discord

type StatusType int8

// Stuff needed to initilize the bot
type Data struct {
	Token  string
	Status *Status
}

// Keeps the status of the bot
type Status struct {
	String string
	// Is only used when 'Type' is 'StatusStreaming'
	URL  string
	Type StatusType
}
