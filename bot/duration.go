package prune_history

import "time"

const (
	Seconds = time.Second
	Minutes = time.Minute
	Hours   = time.Hour
	Days    = 24 * Hours // This doesn't account for daylight savings time
)

func DurationOf(value int, unit time.Duration) time.Duration {
	return unit * time.Duration(value)
}

func LongestDuration(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}

	return b
}

func ShortestDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}

	return b
}
