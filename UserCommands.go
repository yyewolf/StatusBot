package main

import (
	"github.com/bwmarrin/discordgo"
)

func status(s *discordgo.Session, m *discordgo.MessageCreate) {
	timeStr := lastPingTime.Format(timeFormat)

	if embedToSend != nil {
		embedToSend.Footer = &discordgo.MessageEmbedFooter{
			Text: "Mis à jour a " + timeStr + ".",
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
