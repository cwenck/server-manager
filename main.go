package main

import (
	"log"
	"os"
	"time"

	discord "github.com/cwenck/server-manager/api/discord"
	bot "github.com/cwenck/server-manager/bot"
	// hist "github.com/cwenck/server-manager/bot/prune_history"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	log.Printf("BOT_TOKEN=%s\n", token)

	session, err := discord.NewDiscordSession(token)
	if err != nil {
		log.Fatalf("Failed to create a discord session: %s\n", err)
	}
	defer session.Close()

	channel, err := discord.InitDiscordChannel(session, discord.ChannelID("756727763980910692"))
	if err != nil {
		log.Fatalln("Failed to initialize discord channel")
	}

	firstVisitor := discord.NewMessageVisitorBuilder().
		Operation(
			func(message discord.Message, now time.Time) discord.ProcessingStatus {
				// channel.AddMessageReaction(message.MessageID(), "🎃")
				channel.RemoveMessageReaction(message.MessageID(), "🎃")

				// var messageIds []discord.MessageId
				// for _, message := range messages {
				// 	messageIds = append(messageIds, message.MessageId())
				// }

				// channel.BulkDeleteMessages(messageIds)

				return discord.ContinueProcessing
			},
		).
		StartAfterMessageAge(bot.DurationOf(42, bot.Days)).
		// Filter(
		// 	func(message discord.Message) bool {
		// 		// const hoursInDay = 24
		// 		// messageAge := time.Since(message.CreateTime())
		// 		// messageDaysOld := messageAge.Hours() / hoursInDay

		// 		// messageHasCorrectAge := messageDaysOld > 7 && messageDaysOld < 14
		// 		// return messageHasCorrectAge && messageHasCorrectAuthor
		// 		return messageHasCorrectAuthor
		// 	},
		// )
		Build()

	err = channel.VisitMessages([]discord.MessageBatchVisitor{firstVisitor})
	if err != nil {
		log.Fatalf("Error visiting messages: %s", err)
	}
}
