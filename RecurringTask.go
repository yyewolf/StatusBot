package main

import (
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"github.com/go-ping/ping"
)

func singleIPCheck(currentIP, currentName string) {
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

	//Prepares Embed
	var statusTxt string

	statusTxt += "```"
	keys := make([]string, 0)
	longest := 0
	for k := range lastPings {
		keys = append(keys, k)
		if utf8.RuneCountInString(k) > longest {
			longest = utf8.RuneCountInString(k)
		}
	}
	sort.Strings(keys)
	for _, key := range keys {
		pingTime := lastPings[key]
		var spaces string
		for i := 0; i < longest-utf8.RuneCountInString(key); i++ {
			spaces += " "
		}
		statusTxt += key + spaces + " : " + pingTime + "\n"
	}
	statusTxt += "```"

	if len(lastPings) == 0 {
		statusTxt = "Aucune IP n'est enregistré."
	}

	embedToSend = &discordgo.MessageEmbed{
		Title:       "Latences :",
		Description: statusTxt,
		Color:       0xFFDD00,
	}
}

func allIPCheck() {
	for ip, name := range ips {
		go singleIPCheck(ip, name)
	}
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
					isPinging = false
				}
			}
		}
	}()

	//first occurence
	if !isPinging {
		isPinging = true
		allIPCheck()
		isPinging = false
	}
}
