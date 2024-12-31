package utils

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// HasChannelAccess checks if the bot user has access to the specified channel.
func HasChannelAccess(userID string, channel *discordgo.Channel, guild *discordgo.Guild) bool {
	if channel.Type != discordgo.ChannelTypeGuildText && channel.Type != discordgo.ChannelTypeGuildNews {
		return false
	}

	permissions, err := guildMemberPermissions(userID, channel, guild)
	if err != nil {
		return false
	}

	canReadChannel := permissions&discordgo.PermissionReadMessageHistory != 0
	canAccessChannel := permissions&discordgo.PermissionSendMessages != 0
	canViewChannel := permissions&discordgo.PermissionViewChannel != 0

	return canReadChannel && canAccessChannel && canViewChannel
}

// MentionsUser checks if the user is mentioned in the message.
func MentionsUser(message *discordgo.Message, userID string, guild *discordgo.Guild) bool {
	// Check direct mentions
	for _, user := range message.Mentions {
		if user.ID == userID {
			return true
		}
	}

	// Check role mentions
	if guild != nil {
		botMember, err := guildMember(guild, userID)
		if err == nil {
			for _, roleID := range message.MentionRoles {
				for _, botRole := range botMember.Roles {
					if botRole == roleID {
						return true
					}
				}
			}
		}
	}

	return false
}

// guildMemberPermissions calculates the permissions for a member in a channel.
func guildMemberPermissions(userID string, channel *discordgo.Channel, guild *discordgo.Guild) (int64, error) {
	member, err := guildMember(guild, userID)
	if err != nil {
		return 0, err
	}

	// Default permissions based on guild
	permissions := guild.Permissions
	for _, roleID := range member.Roles {
		for _, role := range guild.Roles {
			if role.ID == roleID {
				permissions |= role.Permissions
			}
		}
	}

	// Apply channel-specific permissions
	for _, overwrite := range channel.PermissionOverwrites {
		if overwrite.Type == discordgo.PermissionOverwriteTypeRole {
			for _, roleID := range member.Roles {
				if overwrite.ID == roleID {
					permissions &= ^overwrite.Deny
					permissions |= overwrite.Allow
				}
			}
		} else if overwrite.Type == discordgo.PermissionOverwriteTypeMember && overwrite.ID == userID {
			permissions &= ^overwrite.Deny
			permissions |= overwrite.Allow
		}
	}

	return permissions, nil
}

// guildMember retrieves a member from a guild by user ID.
func guildMember(guild *discordgo.Guild, userID string) (*discordgo.Member, error) {
	for _, member := range guild.Members {
		if member.User.ID == userID {
			return member, nil
		}
	}
	return nil, errors.New("member not found")
}
