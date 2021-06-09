package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"github.com/go-ping/ping"
)

func makeEmbed() {
	//Prepares Embed
	var statusTxt string

	statusTxt += "```"
	keys := make([]string, 0)
	longest := 0
	for k := range tempPings {
		keys = append(keys, k)
		if utf8.RuneCountInString(k) > longest {
			longest = utf8.RuneCountInString(k)
		}
	}
	sort.Strings(keys)
	for _, key := range keys {
		pingTime := tempPings[key]
		var spaces string
		for i := 0; i < longest-utf8.RuneCountInString(key); i++ {
			spaces += " "
		}
		statusTxt += key + spaces + " : " + pingTime + "\n"
	}
	statusTxt += "```"

	if len(config.Ips) == 0 {
		statusTxt = "Aucune IP n'est enregistré."
	}

	embedToSend = &discordgo.MessageEmbed{
		Title:       "Latences :",
		Description: statusTxt,
		Color:       0xFFDD00,
	}
}

func singleIPCheck(wg *sync.WaitGroup, currentIP, currentName string) {
	if wg != nil {
		defer wg.Done()
	}
	currentName = strings.TrimSpace(currentName)

	pinger, err := ping.NewPinger(currentIP)
	pinger.Timeout = 5 * time.Second
	if err != nil {
		return
	}
	pinger.SetPrivileged(true)
	pinger.Count = 3
	err = pinger.Run()
	if err != nil {
		return
	}
	stats := pinger.Statistics()

	pingTime := strconv.FormatInt(stats.AvgRtt.Milliseconds(), 10)

	tempMutex.Lock()
	tempPings[currentName] = pingTime + "ms."

	if stats.PacketLoss == 100 {
		tempPings[currentName] = "Pas de réponse."
	}
	tempMutex.Unlock()

	makeEmbed()
}

func allIPCheck() {
	var wg sync.WaitGroup
	for ip, name := range config.Ips {
		wg.Add(1)
		go singleIPCheck(&wg, ip, name)
	}
	wg.Wait()
	lastPings = tempPings
	lastPingTime = time.Now()
}

func pinger() {
	ticker := time.NewTicker(2 * time.Minute)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				//Do all IPs
				if !isPinging {
					isPinging = true
					allIPCheck()
					updateStatus()
					isPinging = false
				}
			}
		}
	}()

	//first occurence
	if !isPinging {
		isPinging = true
		allIPCheck()
		updateStatus()
		isPinging = false
	}
}

func updateStatus() {
	if config.UpdateMessage.ID == "" {
		return
	}
	s := sess
	timeStr := lastPingTime.Format(timeFormat)

	if embedToSend != nil {
		embedToSend.Footer = &discordgo.MessageEmbedFooter{
			Text: "Mis à jour a " + timeStr + ".",
		}
		msg, err := s.ChannelMessageEditEmbed(config.UpdateMessage.ChannelID, config.UpdateMessage.ID, embedToSend)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(msg.ID)
		config.UpdateMessage.ID = msg.ID
		config.UpdateMessage.ChannelID = msg.ChannelID
		saveData()
		return
	}
	var statusTxt string
	if len(config.Ips) == 0 {
		statusTxt = "Aucune IP n'est enregistré."
	}

	embedToSend = &discordgo.MessageEmbed{
		Title:       "Latences :",
		Description: statusTxt,
		Color:       0xFFDD00,
	}
	_, err := s.ChannelMessageEditEmbed(config.UpdateMessage.ChannelID, config.UpdateMessage.ID, embedToSend)
	if err != nil {
		msg, err := s.ChannelMessageSendEmbed(config.UpdateMessage.ChannelID, embedToSend)
		if err != nil {
			config.UpdateMessage.ID = ""
			saveData()
			return
		}

		config.UpdateMessage.ID = msg.ID
		config.UpdateMessage.ChannelID = msg.ChannelID
		saveData()
		return
	}
}
