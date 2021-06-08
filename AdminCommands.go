package main

import (
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
	if !hasPermission(member, s, m.GuildID, PERM_BAN_MEMBERS) {
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

	ips[ipToAdd] = nameOfIP
	saveData()

	singleIPCheck(ipToAdd, nameOfIP)
	lastPings = tempPings
	lastPingTime = time.Now()
}

func remIP(s *discordgo.Session, m *discordgo.MessageCreate) {
	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		return
	}
	//Check that user has ban perms
	if !hasPermission(member, s, m.GuildID, PERM_BAN_MEMBERS) {
		return
	}

	args := strings.Split(m.Content, " ")
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Cette commande s'utilise comme cela : `s!remip 127.0.0.1`")
		return
	}
	ipToRemove := args[1]

	if val, ok := ips[ipToRemove]; ok {
		delete(lastPings, val)
		delete(tempPings, val)
		delete(ips, ipToRemove)
		s.ChannelMessageSend(m.ChannelID, "J'ai supprimé : `"+ipToRemove+"`")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Cette IP n'est pas enregistré : `"+ipToRemove+"`")
	}

	saveData()
}
