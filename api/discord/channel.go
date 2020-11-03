package discord

import (
	"errors"

	discord "github.com/bwmarrin/discordgo"
)

// Channel is the interface for operations on a discord channel.
type Channel interface {
	CreateMessage(content string) (messageID MessageID, err error)
	AddMessageReaction(messageID MessageID, emoji string) (err error)
	RemoveMessageReaction(messageID MessageID, emoji string) (err error)
	BulkDeleteMessages(messageIds []MessageID) (err error)
	DeleteMessage(messageID MessageID) (err error)
	VisitMessages(visitors []MessageBatchVisitor) (err error)
}

type discordChannel struct {
	session   Session
	channelID ChannelID
}

// InitDiscordChannel initializes a channel with the specified channel ID with a session.
func InitDiscordChannel(session Session, channelID ChannelID) (result Channel, err error) {
	ok, err := session.isValid()
	if !ok {
		return
	}

	result = discordChannel{
		session:   session,
		channelID: channelID,
	}

	return
}

func (channel discordChannel) CreateMessage(content string) (messageID MessageID, err error) {
	message, err := channel.session.session.ChannelMessageSend(string(channel.channelID), content)
	if err != nil {
		return
	}

	messageID = MessageID(message.ID)
	return
}

func (channel discordChannel) VisitMessages(visitors []MessageBatchVisitor) (err error) {
	batchSize := 5

	getMessagesBefore := func(beforeMessageId string) ([]*discord.Message, error) {
		return channel.session.session.ChannelMessages(string(channel.channelID), batchSize, beforeMessageId, "", "")
	}

	processBatch := func(apiMessageBatch []*discord.Message) {
		for _, visitor := range visitors {
			visitor.VisitMessageBatch(convertMessages(apiMessageBatch))
		}
	}

	startingMessageID, err := visitorMostRecentStartAfterMessageID(visitors)
	if err != nil {
		return
	}

	currentMessageBatch, err := getMessagesBefore(string(startingMessageID))
	if err != nil {
		return
	}

	// Keep looping until there is a batch that isn't full
	for len(currentMessageBatch) == batchSize {
		// Load the next batch before processing the current batch incase any messages in the current batch get deleted
		lastMessageIndex := batchSize - 1
		lastCurrentMessageID := currentMessageBatch[lastMessageIndex].ID
		nextMessageBatch, err := getMessagesBefore(lastCurrentMessageID)
		if err != nil {
			return err
		}

		// Process the current message batch and setup the next message batch
		processBatch(currentMessageBatch)
		currentMessageBatch = nextMessageBatch
	}

	// Process the final non-full message batch
	processBatch(currentMessageBatch)
	return
}

func (channel discordChannel) AddMessageReaction(messageID MessageID, emoji string) (err error) {
	return channel.session.session.MessageReactionAdd(string(channel.channelID), string(messageID), emoji)
}

func (channel discordChannel) RemoveMessageReaction(messageID MessageID, emoji string) (err error) {
	userID := "@me" // This is a value specified by the discord API to refer to the current user (i.e. the bot)
	return channel.session.session.MessageReactionRemove(string(channel.channelID), string(messageID), emoji, userID)
}

func (channel discordChannel) DeleteMessage(messageID MessageID) (err error) {
	return channel.session.session.ChannelMessageDelete(string(channel.channelID), string(messageID))
}

func (channel discordChannel) BulkDeleteMessages(messageIds []MessageID) (err error) {
	messageCount := len(messageIds)

	if messageCount > 100 {
		err = errors.New("At most 100 messages can be bulk deleted at once")
	} else if messageCount == 0 {
		// Do nothing
	} else if messageCount == 1 {
		err = channel.DeleteMessage(messageIds[0])
	} else {
		err = channel.session.session.ChannelMessagesBulkDelete(string(channel.channelID), messageIdsAsStrings(messageIds))
	}

	return
}

func messageIdsAsStrings(messageIDs []MessageID) (result []string) {
	for _, messageID := range messageIDs {
		result = append(result, string(messageID))
	}
	return
}

func visitorMostRecentStartAfterMessageID(visitors []MessageBatchVisitor) (MessageID, error) {
	var messageIds []MessageID
	for _, visitor := range visitors {
		messageIds = append(messageIds, visitor.StartAfterMessageWithID())
	}

	return mostRecentMessageID(messageIds)
}
