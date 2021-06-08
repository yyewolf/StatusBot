package main

import (
	"fmt"
	"math"
	"time"

	"github.com/bwmarrin/discordgo"
)

func status(s *discordgo.Session, m *discordgo.MessageCreate) {
	var secondsAgo string
	i, f := math.Modf(time.Now().Sub(lastPingTime).Seconds())
	if f < 0 {
		f = -f
	}
	secondsAgo = fmt.Sprint(i)

	if embedToSend != nil {
		embedToSend.Footer = &discordgo.MessageEmbedFooter{
			Text: "Mis à jour il y a " + secondsAgo + " secondes.",
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embedToSend)
		return
	}
	var statusTxt string
	if len(lastPings) == 0 {
		statusTxt = "Aucune IP n'est enregistré."
	}

	embedToSend = &discordgo.MessageEmbed{
		Title:       "Latences :",
		Description: statusTxt,
		Color:       0xFFDD00,
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embedToSend)
}
