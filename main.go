package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"./commands"
	"./config"
	"./utils"
)

func main() {
	config.InitConfig()

	discord, err := discordgo.New("Bot " + config.Configuration.AuthToken)
	utils.Check(err)

	discord.AddHandler(onMessage)

	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = discord.Open()
	utils.Check(err)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

func onMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	for _, command := range commands.GetCommands() {
		if checkMessageForCommand(message.Content, command.Activation) {
			command.Action(session, message)
		}
	}
}

func checkMessageForCommand(message string, command string) bool {
	return strings.HasPrefix(message, fmt.Sprintf("%s%s", config.Configuration.Prefix, command))
}
