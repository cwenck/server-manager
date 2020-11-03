package discord

import (
	"errors"
	"strconv"
	"time"
)

const discordEpoch uint64 = 1420070400000

// MessageID is the snowflake ID for a discord message
type MessageID string

// ChannelID is the snowflake ID for a discord channel
type ChannelID string

// UserID is the snowflake ID for a discord user
type UserID string

func messageIDForTimestamp(timestamp time.Time) (messageID MessageID, err error) {
	unixEpochSec := uint64(timestamp.Unix())
	unixEpochMs := uint64(unixEpochSec) * 1000

	if unixEpochMs < discordEpoch {
		err = errors.New("Timestamp must be after the discord epoch")
		return
	}

	// Convert to the discord snowflake as specified in the Discord API documentation
	messageIDAsNumber := (unixEpochMs - discordEpoch) << 22
	messageIDAsString := strconv.FormatUint(messageIDAsNumber, 10)
	messageID = MessageID(messageIDAsString)
	return
}

func messageIDForMessageAge(messageAge time.Duration) (messageID MessageID, err error) {
	if messageAge < 0 {
		err = errors.New("Message age must be positive")
		return
	}

	messageCreateTime := time.Now().Add(-messageAge)
	return messageIDForTimestamp(messageCreateTime)
}

// compareMessageIDs compares the the two message IDs.
// Returns 0 if a = b.
// Returns a value less than 0 if a < b.
// Returns a value more than 0 if a > b.
func compareMessageIDs(a, b MessageID) (result int, err error) {
	if a == b {
		result = 0
		return
	}

	parsedA, err := strconv.ParseUint(string(a), 10, 64)
	if err != nil {
		return
	}

	parsedB, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return
	}

	if parsedA < parsedB {
		result = -1
	} else {
		result = 1
	}

	return
}

// mostRecentMessageID picks the message ID that is the most recent from the slice.
// If the slice is empty, then an empty message ID is returned with no error.
func mostRecentMessageID(messageIDs []MessageID) (result MessageID, _ error) {
	for _, messageID := range messageIDs {
		if result == "" {
			result = messageID
		} else {
			comparison, err := compareMessageIDs(messageID, result)
			if err != nil {
				return "", err
			} else if comparison > 0 {
				// The new message ID is greater than the current result,
				// so the new message ID is more recent
				result = messageID
			}
		}
	}
	return
}
