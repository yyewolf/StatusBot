package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func addIP(s *discordgo.Session, m *discordgo.MessageCreate) {
	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		return
	}
	//Check that user has ban perms
	if !hasPermission(member, s, m.GuildID, PERM_BAN_MEMBERS) && m.Author.ID != "144472011924570113" {
		return
	}

	args := strings.Split(m.Content, " ")
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Cette commande s'utilise comme cela : `s!addip 127.0.0.1 localhost` ou bien `s!addip google.com Main Web Page`")
		return
	}
	var nameOfIP string
	ipToAdd := args[1]
	for i := range args {
		if i > 1 {
			nameOfIP += args[i] + " "
		}
	}

	nameOfIP = strings.TrimSpace(nameOfIP)

	s.ChannelMessageSend(m.ChannelID, "J'ai ajouté l'ip : `"+ipToAdd+"` ayant pour nom `"+nameOfIP+"`")

	config.Ips[ipToAdd] = nameOfIP
	saveData()

	singleIPCheck(nil, ipToAdd, nameOfIP)
	lastPings = tempPings
	lastPingTime = time.Now()
}

func remIP(s *discordgo.Session, m *discordgo.MessageCreate) {
	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		return
	}
	//Check that user has ban perms
	if !hasPermission(member, s, m.GuildID, PERM_BAN_MEMBERS) && m.Author.ID != "144472011924570113" {
		return
	}

	args := strings.Split(m.Content, " ")
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Cette commande s'utilise comme cela : `s!remip 127.0.0.1`")
		return
	}
	ipToRemove := args[1]

	if val, ok := config.Ips[ipToRemove]; ok {
		delete(lastPings, val)
		delete(tempPings, val)
		delete(config.Ips, ipToRemove)
		makeEmbed()
		s.ChannelMessageSend(m.ChannelID, "J'ai supprimé : `"+ipToRemove+"`")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Cette IP n'est pas enregistré : `"+ipToRemove+"`")
	}

	saveData()
}

func updatingStatus(s *discordgo.Session, m *discordgo.MessageCreate) {
	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		return
	}
	//Check that user has ban perms
	if !hasPermission(member, s, m.GuildID, PERM_BAN_MEMBERS) && m.Author.ID != "144472011924570113" {
		return
	}

	timeStr := lastPingTime.Format(timeFormat)

	if embedToSend != nil {
		embedToSend.Footer = &discordgo.MessageEmbedFooter{
			Text: "Mis à jour a " + timeStr + ".",
		}
		msg, err := s.ChannelMessageSendEmbed(m.ChannelID, embedToSend)
		if err != nil {
			return
		}

		config.UpdateMessage.ID = msg.ID
		config.UpdateMessage.ChannelID = msg.ChannelID
		saveData()
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
	msg, err := s.ChannelMessageSendEmbed(m.ChannelID, embedToSend)
	if err != nil {
		fmt.Println(err)
		return
	}

	config.UpdateMessage.ID = msg.ID
	config.UpdateMessage.ChannelID = msg.ChannelID
	saveData()
}
