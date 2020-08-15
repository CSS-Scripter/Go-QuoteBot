package commands

import "github.com/bwmarrin/discordgo"

// Command is a structure used for commands
type Command struct {
	Activation string
	Action     func(session *discordgo.Session, message *discordgo.MessageCreate)
}

// GetCommands retrieves all registered commands
func GetCommands() []Command {
	return []Command{
		{
			Activation: "quote",
			Action:     Quote,
		},
	}
}
