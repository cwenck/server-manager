package discord

import discord "github.com/bwmarrin/discordgo"

// User is the interface for accessing discord user metadata.
type User interface {
	UserID() UserID
	Username() string
	Email() string
	Discriminator() string
}

type discordUser struct {
	userID        UserID
	username      string
	email         string
	discriminator string
}

func newDiscordUser(apiUser *discord.User) discordUser {
	return discordUser{
		userID:        UserID(apiUser.ID),
		username:      apiUser.Username,
		email:         apiUser.Email,
		discriminator: apiUser.Discriminator,
	}
}

func (user discordUser) UserID() UserID {
	return user.userID
}

func (user discordUser) Username() string {
	return user.username
}

func (user discordUser) Email() string {
	return user.email
}

func (user discordUser) Discriminator() string {
	return user.discriminator
}
