package discord

// import (
// 	"errors"

// 	discord "github.com/bwmarrin/discordgo"
// )

// type channel struct {
// 	session   Session
// 	channelId ChannelId
// }

// func InitDiscordChannel(session Session, channelId ChannelId) (result Channel, err error) {
// 	ok, err := session.isValid()
// 	if !ok {
// 		return
// 	}

// 	result = channel{
// 		session:   session,
// 		channelId: channelId,
// 	}

// 	return
// }

// func (self channel) CreateMessage(content string) (messageId MessageId, err error) {
// 	message, err := self.session.session.ChannelMessageSend(string(self.channelId), content)
// 	if err != nil {
// 		return
// 	}

// 	messageId = MessageId(message.ID)
// 	return
// }

// // func (self channel) VisitMessages(operation func(message Message)) (err error) {
// // 	return self.VisitFilteredMessages(operation, func(message Message) bool { return true })
// // }

// // func (self channel) VisitFilteredMessages(operation func(message Message), filter func(message Message) bool) (err error) {
// // 	batchOperation := func(messages []Message) {
// // 		for _, message := range messages {
// // 			operation(message)
// // 		}
// // 	}

// // 	return self.VisitFilteredMessagesInBatches(100, batchOperation, filter)
// // }

// // func (self channel) VisitMessagesInBatches(batchSize int, operation func(messages []Message)) (err error) {
// // 	return self.VisitFilteredMessagesInBatches(batchSize, operation, func(message Message) bool { return true })
// // }

// // func (self channel) VisitFilteredMessagesInBatches(batchSize int, operation func(messages []Message), filter func(message Message) bool) (err error) {
// // 	if batchSize < 1 || batchSize > 100 {
// // 		msg := fmt.Sprintf("Expected batch size must be between 1 and 100, but got %d", batchSize)
// // 		return errors.New(msg)
// // 	}

// // 	processBatch := func(messageBatch []*discord.Message) {
// // 		// Loop through each message in the current message batch
// // 		var filteredMessages []Message
// // 		for _, rawMessage := range messageBatch {
// // 			message := newDiscordMessage(rawMessage)
// // 			if filter(message) {
// // 				filteredMessages = append(filteredMessages, message)
// // 			}
// // 		}

// // 		// Run the operation of the slice of filtered messages
// // 		operation(filteredMessages)
// // 	}

// // 	currentMessageBatch, err := self.getMessagesBefore(batchSize, "")
// // 	if err != nil {
// // 		return
// // 	}

// // 	// Keep looping until there is a batch that isn't full
// // 	for len(currentMessageBatch) == batchSize {
// // 		// Load the next batch before processing the current batch incase any messages in the current batch get deleted
// // 		lastMessageIndex := batchSize - 1
// // 		lastCurrentMessageId := currentMessageBatch[lastMessageIndex].ID
// // 		nextMessageBatch, err := self.getMessagesBefore(batchSize, lastCurrentMessageId)
// // 		if err != nil {
// // 			return err
// // 		}

// // 		// Process the current message batch and setup the next message batch
// // 		processBatch(currentMessageBatch)
// // 		currentMessageBatch = nextMessageBatch
// // 	}

// // 	// Process the final non-full message batch
// // 	processBatch(currentMessageBatch)
// // 	return
// // }

// func (self channel) VisitMessages(visitor MessageBatchVisitor) (err error) {
// 	batchSize := 100
// 	getMessagesBefore := func(beforeMessageId string) ([]*discord.Message, error) {
// 		return self.session.session.ChannelMessages(string(self.channelId), batchSize, beforeMessageId, "", "")
// 	}

// 	processBatch := func(messageBatch []*discord.Message) {
// 		var convertedMessages []Message
// 		for _, message := range messageBatch {
// 			convertedMessages = append(convertedMessages, newDiscordMessage(message))
// 		}

// 		visitor.VisitBatch(convertedMessages)
// 	}

// 	currentMessageBatch, err := getMessagesBefore("")
// 	if err != nil {
// 		return
// 	}

// 	// Keep looping until there is a batch that isn't full
// 	for len(currentMessageBatch) == batchSize {
// 		// Load the next batch before processing the current batch incase any messages in the current batch get deleted
// 		lastMessageIndex := batchSize - 1
// 		lastCurrentMessageId := currentMessageBatch[lastMessageIndex].ID
// 		nextMessageBatch, err := getMessagesBefore(lastCurrentMessageId)
// 		if err != nil {
// 			return err
// 		}

// 		// Process the current message batch and setup the next message batch
// 		processBatch(currentMessageBatch)
// 		currentMessageBatch = nextMessageBatch
// 	}

// 	// Process the final non-full message batch
// 	processBatch(currentMessageBatch)
// 	return
// }

// func (self channel) AddMessageReaction(messageId MessageId, emoji string) (err error) {
// 	return self.session.session.MessageReactionAdd(string(self.channelId), string(messageId), emoji)
// }

// func (self channel) RemoveMessageReaction(messageId MessageId, emoji string) (err error) {
// 	userId := "@me" // This is a value specified by the discord API to refer to the current user (i.e. the bot)
// 	return self.session.session.MessageReactionRemove(string(self.channelId), string(messageId), emoji, userId)
// }

// func (self channel) DeleteMessage(messageId MessageId) (err error) {
// 	return self.session.session.ChannelMessageDelete(string(self.channelId), string(messageId))
// }

// func (self channel) BulkDeleteMessages(messageIds []MessageId) (err error) {
// 	messageCount := len(messageIds)

// 	if messageCount > 100 {
// 		return errors.New("At most 100 messages can be bulk deleted at once")
// 	}

// 	if messageCount == 0 {
// 		// Do nothing
// 		return
// 	} else if messageCount == 1 {
// 		return self.DeleteMessage(messageIds[0])
// 	} else {
// 		return self.session.session.ChannelMessagesBulkDelete(string(self.channelId), messageIdsAsStrings(messageIds))
// 	}
// }

// func messageIdsAsStrings(messageIds []MessageId) (result []string) {
// 	for _, messageId := range messageIds {
// 		result = append(result, string(messageId))
// 	}
// 	return
// }
