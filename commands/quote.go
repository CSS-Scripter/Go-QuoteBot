package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"app/config"
	"app/utils"
)

type quoteStruct struct {
	Message string `json:"message"`
	By      string `json:"by"`
	Year    string `json:"year"`
}

// Quote extracts the quote and sender from the message, and sends it to the quote site.
func Quote(session *discordgo.Session, message *discordgo.MessageCreate) {
	content := message.ContentWithMentionsReplaced()
	messageContent := strings.TrimPrefix(content, fmt.Sprintf("%squote ", config.Configuration.Prefix))
	quoteAndBy := strings.SplitAfter(messageContent, "-")

	if len(quoteAndBy) <= 1 {
		session.ChannelMessageSend(message.ChannelID, "Looks like you forgot to quote someone!")
		session.ChannelMessageSend(message.ChannelID, "Usage: !quote [quote] - [quoted by]")
		return
	}

	quote := strings.Join(quoteAndBy[:len(quoteAndBy)-1], "")
	quote = strings.TrimSuffix(quote, "-")
	quote = strings.TrimSpace(quote)

	by := quoteAndBy[len(quoteAndBy)-1]
	by = strings.TrimSpace(by)
	by = strings.TrimPrefix(by, "@")

	requestBody := createQuoteRequestBody(by, quote)
	targetAddress := fmt.Sprintf("http://%s:%d/%s/new", config.Configuration.Hostname, config.Configuration.Port, message.GuildID)

	resp, err := http.Post(
		targetAddress,
		"application/json",
		bytes.NewBuffer(requestBody))

	utils.Check(err)

	if resp.Status == "200 OK" {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s - %s", quote, by))
	} else {
		session.ChannelMessageSend(message.ChannelID, "Looks like something went wrong! Try again later!")
	}
}

func createQuoteRequestBody(by string, message string) []byte {
	year, _, _ := time.Now().Date()
	quote := quoteStruct{
		By:      by,
		Message: message,
		Year:    fmt.Sprintf("%d", year),
	}

	requestBody, _ := json.Marshal(quote)
	return requestBody
}
