package discord

import (
	"log"
	"time"

	discord "github.com/bwmarrin/discordgo"
)

// Message is the interface for accessing discord message metadata.
type Message interface {
	MessageID() MessageID
	ChannelID() ChannelID
	Content() string
	Pinned() bool
	Author() User
	CreateTime() time.Time
	Age(time.Time) time.Duration
}

type discordMessage struct {
	messageID  MessageID
	channelID  ChannelID
	content    string
	pinned     bool
	author     discordUser
	createTime time.Time
}

func newMessage(apiMessage *discord.Message) discordMessage {
	createTime, err := discord.SnowflakeTimestamp(apiMessage.ID)
	if err != nil {
		log.Fatalf("Failed to parse timestamp from snowflake ID: %s", err)
	}

	return discordMessage{
		messageID:  MessageID(apiMessage.ID),
		channelID:  ChannelID(apiMessage.ChannelID),
		content:    apiMessage.Content,
		pinned:     apiMessage.Pinned,
		author:     newDiscordUser(apiMessage.Author),
		createTime: createTime,
	}
}

func (message discordMessage) MessageID() MessageID {
	return message.messageID
}

func (message discordMessage) ChannelID() ChannelID {
	return message.channelID
}

func (message discordMessage) Content() string {
	return message.content
}

func (message discordMessage) Pinned() bool {
	return message.pinned
}

func (message discordMessage) Author() User {
	return message.author
}

func (message discordMessage) CreateTime() time.Time {
	return message.createTime
}

func (message discordMessage) Age(now time.Time) time.Duration {
	return now.Sub(message.createTime)
}

func convertMessages(apiMessages []*discord.Message) (result []Message) {
	for _, apiMessage := range apiMessages {
		result = append(result, newMessage(apiMessage))
	}
	return
}
