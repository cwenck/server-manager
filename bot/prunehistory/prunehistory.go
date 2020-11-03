package prunehistory

import (
	"log"
	"time"

	discord "github.com/cwenck/server-manager/api/discord"
	bot "github.com/cwenck/server-manager/bot"
)

// PruneHistory prunes the history of the specified channel.
// Messages will be deleted if they are older than the keep message duration.
func PruneHistory(channel discord.Channel, keepMessageDuration time.Duration) {
	visitors := []discord.MessageBatchVisitor{
		bulkDeleteVisitor(channel, keepMessageDuration),
		slowDeleteVisitor(channel, keepMessageDuration),
	}

	channel.VisitMessages(visitors)
}

func bulkDeleteVisitor(channel discord.Channel, keepMessageDuration time.Duration) discord.MessageBatchVisitor {
	maxMessageAge := bot.DurationOf(14, bot.Days)
	minMessageAge := keepMessageDuration

	return discord.NewMessageVisitorBuilder().
		StartAfterMessageAge(minMessageAge).
		Filter(
			func(message discord.Message, now time.Time) bool {
				messageAge := message.Age(now)
				return minMessageAge < messageAge && messageAge < maxMessageAge && !message.Pinned()
			},
		).
		BulkOperation(
			func(messages []discord.Message, now time.Time) discord.ProcessingStatus {
				var messageIds []discord.MessageID
				for _, message := range messages {
					messageIds = append(messageIds, message.MessageID())
				}

				err := channel.BulkDeleteMessages(messageIds)
				if err != nil {
					log.Printf("Failed to delete bulk delete messages: %s\n", err)
				}

				return keepProcessingBulkDelete(messages, now)
			},
		).
		Build()
}

func slowDeleteVisitor(channel discord.Channel, keepMessageDuration time.Duration) discord.MessageBatchVisitor {
	minMessageAge := bot.LongestDuration(keepMessageDuration, bot.DurationOf(14, bot.Days))

	return discord.NewMessageVisitorBuilder().
		StartAfterMessageAge(minMessageAge).
		Filter(
			func(message discord.Message, now time.Time) bool {
				messageAge := message.Age(now)
				return minMessageAge < messageAge && !message.Pinned()
			},
		).
		Operation(
			func(message discord.Message, now time.Time) discord.ProcessingStatus {
				err := channel.DeleteMessage(message.MessageID())
				if err != nil {
					log.Printf("Failed to delete message with ID %s: %s\n", message.MessageID(), err)
				}

				return discord.ContinueProcessing
			},
		).
		Build()
}

func keepProcessingBulkDelete(messages []discord.Message, now time.Time) discord.ProcessingStatus {
	for _, message := range messages {
		if message.Age(now) > bot.DurationOf(14, bot.Days) {
			return discord.StopProcessing
		}
	}

	return discord.ContinueProcessing
}
