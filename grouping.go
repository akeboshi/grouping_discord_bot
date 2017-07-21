package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v2"
)

type BotConfig struct {
	Token           string
	GroupingTrigger string "grouping_trigger"
}

var config BotConfig

func init() {
	configText, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configText, &config)
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
}

func main() {
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running, Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == config.GroupingTrigger {
		channel, _ := s.State.Channel(m.ChannelID)
		guild, _ := s.State.Guild(channel.GuildID)

		_, channelID := getJoinedChannelID(m.Author.ID, guild.VoiceStates)

		userIDs := getChannelJoinedUserIDs(channelID, guild.VoiceStates)

		// get user nick or name by userID
		userNames := []string{}
		for _, uid := range userIDs {
			userNames = append(userNames, getDisplayUserName(uid, guild.Members))
		}

		shuffle(userNames)

		size := len(userNames)/2 + len(userNames)%2
		s.ChannelMessageSend(m.ChannelID, "red team:  "+strings.Join(userNames[:size], ","))
		s.ChannelMessageSend(m.ChannelID, "blue team: "+strings.Join(userNames[size:], ","))
	}
}

func getDisplayUserName(uid string, members []*discordgo.Member) string {
	for _, member := range members {
		if member.User.ID == uid {
			name := member.User.Username
			if member.Nick != "" {
				name = member.Nick
			}
			return name
		}
	}
	return ""
}

// get useIDs joined voice channel with author.
func getChannelJoinedUserIDs(channelID string, states []*discordgo.VoiceState) []string {
	userIDs := []string{}
	if channelID != "" {
		for _, v := range states {
			if v.ChannelID == channelID {
				userIDs = append(userIDs, v.UserID)
			}
		}
	}
	return userIDs
}

// get author joined voice channel
func getJoinedChannelID(uid string, states []*discordgo.VoiceState) (error, string) {
	for _, v := range states {
		if v.UserID == uid {
			return nil, v.ChannelID
		}
	}

	return fmt.Errorf("author don't join channel"), ""
}

func shuffle(data []string) {
	n := len(data)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		data[i], data[j] = data[j], data[i]
	}
}

// unused method
func getChannelByChannelID(id string, channels []*discordgo.Channel) (error, *discordgo.Channel) {
	for _, c := range channels {
		if c.ID == id {
			return nil, c
		}
	}
	return fmt.Errorf("channel is not found"), nil
}

// unused method
func getVoiceChannels(guildID string, s *discordgo.Session) []*discordgo.Channel {
	return getChannelsOfKind(guildID, s, "voice")
}

// unused method
// should specify "voice" or "text" for kind
func getChannelsOfKind(guildID string, s *discordgo.Session, kind string) []*discordgo.Channel {
	guild, err := s.State.Guild(guildID)
	channels := []*discordgo.Channel{}
	if err != nil {
		fmt.Println(err)
		return channels
	}
	for _, c := range guild.Channels {
		if c.Type == kind {
			channels = append(channels, c)
		}
	}
	return channels
}
