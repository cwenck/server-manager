package discord

import (
	"errors"
	"strings"

	discord "github.com/bwmarrin/discordgo"
)

// Session holds discord session information.
type Session struct {
	session *discord.Session
	open    bool
}

// NewDiscordSession creates a new session from the specified token.
func NewDiscordSession(token string) (result Session, err error) {
	normalizedToken := strings.TrimSpace(token)
	if normalizedToken == "" {
		err = errors.New("Missing token")
		return
	}
	session, err := discord.New("Bot " + normalizedToken)
	if err != nil {
		return
	}

	err = session.Open()
	if err != nil {
		return
	}

	result = Session{
		session: session,
		open:    true,
	}
	return
}

// Close closes the session.
func (session Session) Close() {
	ok, _ := session.isValid()
	if ok && session.open {
		session.session.Close()
	}
}

func (session Session) isValid() (ok bool, err error) {
	ok = session.session != nil
	if !ok {
		err = errors.New("Invalid discord session")
	}

	return
}
