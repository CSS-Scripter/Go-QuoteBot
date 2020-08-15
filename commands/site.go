package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

//Site sends the website with all quotes from the server where the command was send from
func Site(session *discordgo.Session, message *discordgo.MessageCreate) {
	serverID := message.GuildID
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("https://quote.mylocalhost.app/%s", serverID))
}
