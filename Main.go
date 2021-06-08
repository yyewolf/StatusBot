//go:generate goversioninfo
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var sess *discordgo.Session

func init() {
	ips = make(map[string]string)
	commands = make(map[string]func(*discordgo.Session, *discordgo.MessageCreate))
	lastPings = make(map[string]string)
	tempPings = make(map[string]string)

	defineCommands()
	loadData()
	pinger()
}

func main() {
	var err error
	sess, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Faled to start: ", err)
	}

	sess.AddHandler(botReady)
	sess.AddHandler(commandHandler)

	log.Println("Starting the shard manager")
	err = sess.Open()
	if err != nil {
		log.Fatal("Failed to start: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	sess.Close()
}

func botReady(s *discordgo.Session, evt *discordgo.Ready) {
	s.UpdateGameStatus(0, "s!status to check status")
}
