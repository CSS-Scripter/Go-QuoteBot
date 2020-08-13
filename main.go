package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tkanos/gonfig"
	"fmt"
	"os"
	"syscall"
	"os/signal"
	"strings"
	"net/http"
	"bytes"
	"time"
)

type Configuration struct {
	Prefix string
	AuthToken string
	Hostname string
	Port int
}

var config = Configuration{}

func main() {
	initConfig()

	discord, err := discordgo.New("Bot " + config.AuthToken)
	check(err)

	discord.AddHandler(onMessage)
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	
	err = discord.Open()
	check(err)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

func initConfig() {
	err := gonfig.GetConf("./config.json", &config)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func onMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	content := message.ContentWithMentionsReplaced()

	if (checkMessageForCommand(content, "quote")) {
		session.ChannelMessageSend(message.ChannelID, "I'm currently fucking around with the bot, it'l be back later on or something, I dunno")
		return
	}

	if (checkMessageForCommand(content, "easy")) {
		messageContent := strings.TrimPrefix(content, fmt.Sprintf("%seasy ", config.Prefix))
		quoteAndBy := strings.SplitAfter(messageContent, "-")
		
		if (len(quoteAndBy) <= 1) {
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
		targetAddress := fmt.Sprintf("http://%s:%d/%s/new", config.Hostname, config.Port, message.GuildID)

		resp, err := http.Post(
			targetAddress, 
			"application/json", 
			bytes.NewBuffer([]byte(requestBody)))

		check(err)

		if (resp.Status == "200 OK") {
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s | %s", quote, by))
		} else {
			session.ChannelMessageSend(message.ChannelID, "Looks like something went wrong! Try again later!")
		}
	}
}

func checkMessageForCommand(message string, command string) bool {
	return strings.HasPrefix(message, fmt.Sprintf("%s%s", config.Prefix, command))
}

func createQuoteRequestBody(by string, quote string) string {
	year, _, _ := time.Now().Date()
	return fmt.Sprintf(
		`{
			"message":"%s",
			"by":"%s",
			"year":"%d"
		}`, quote, by, year)
}
