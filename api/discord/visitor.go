package discord

import (
	"log"
	"time"
)

// ProcessingStatus is the result of processing one or more items.
// It indicates what actions should be taken on any items that have yet to be processed.
type ProcessingStatus bool

const (
	// ContinueProcessing if there are more items to process.
	ContinueProcessing ProcessingStatus = true
	// StopProcessing even if there are more items to process.
	StopProcessing = false
)

// MessageBatchVisitor is the interface for visiting messages in batches
type MessageBatchVisitor interface {
	StartAfterMessageWithID() MessageID
	VisitMessageBatch(messages []Message) ProcessingStatus
}

// MessagePredicate defines a predicate that takes a message as input
type MessagePredicate func(message Message, now time.Time) bool

// MessageVisitorOperation defines an operation on a single message to execute when visiting the message
type MessageVisitorOperation func(message Message, now time.Time) ProcessingStatus

// MessageVisitorBulkOperation defines an operation on a batch of messages to execute when visiting the messages
type MessageVisitorBulkOperation func(message []Message, now time.Time) ProcessingStatus

type messageVisitor struct {
	filter              MessagePredicate
	bulkOperation       MessageVisitorBulkOperation
	startAfterMessageID MessageID
	now                 time.Time
}

// MessageVisitorBuilder a builder for a message visitor
type MessageVisitorBuilder messageVisitor

// NewMessageVisitorBuilder creates a builder for an implementation of a MessageBatchVisitor.
func NewMessageVisitorBuilder() MessageVisitorBuilder {
	return MessageVisitorBuilder{
		filter:        func(message Message, now time.Time) bool { return true },
		bulkOperation: func(message []Message, now time.Time) ProcessingStatus { return ContinueProcessing },
	}
}

func (visitorBuilder MessageVisitorBuilder) Operation(operation MessageVisitorOperation) MessageVisitorBuilder {
	visitorBuilder.bulkOperation = func(messages []Message, now time.Time) ProcessingStatus {
		for _, message := range messages {
			if operation(message, now) == StopProcessing {
				return StopProcessing
			}
		}
		return ContinueProcessing
	}

	return visitorBuilder
}

func (visitorBuilder MessageVisitorBuilder) BulkOperation(bulkOperation MessageVisitorBulkOperation) MessageVisitorBuilder {
	visitorBuilder.bulkOperation = bulkOperation
	return visitorBuilder
}

func (visitorBuilder MessageVisitorBuilder) Filter(filter MessagePredicate) MessageVisitorBuilder {
	visitorBuilder.filter = filter
	return visitorBuilder
}

func (visitorBuilder MessageVisitorBuilder) StartAfterMessageAge(messageAge time.Duration) MessageVisitorBuilder {
	messageID, err := messageIDForMessageAge(messageAge)
	if err != nil {
		log.Fatalf("Invalid message age: %s", err)
	}

	visitorBuilder.startAfterMessageID = messageID
	return visitorBuilder
}

func (visitorBuilder MessageVisitorBuilder) Build() messageVisitor {
	return messageVisitor(visitorBuilder)
}

////////////////////////////////////////
// MessageBatchVisitor Implementation //
////////////////////////////////////////

// VisitMessageBatch provides an implementation for the MessageBatchVisitor interface.
func (visitor messageVisitor) VisitMessageBatch(messages []Message) ProcessingStatus {
	now := time.Now()

	var filteredMessages []Message
	for _, message := range messages {
		if visitor.filter(message, now) {
			filteredMessages = append(filteredMessages, message)
		}
	}

	return visitor.bulkOperation(filteredMessages, now)
}

func (visitor messageVisitor) StartAfterMessageWithID() MessageID {
	return visitor.startAfterMessageID
}

////////////////////////////////////////
