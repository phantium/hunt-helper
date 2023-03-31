package handlers

import (
	"crypto/md5"
	"discordbot/internal/common"
	"discordbot/internal/orm"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

func ConfigureBrowseChannelID(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!set channel browse" && common.MemberHasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionAdministrator) {
		orm.UpdateGuildConfig(m.GuildID, map[string]interface{}{"channel_browse": m.ChannelID})
		s.ChannelMessageSend(m.ChannelID, "Browse Channel ID set to: "+m.ChannelID+"!")
		log.Infof("Browse Channel ID set to: %s on Guild ID: %s", m.ChannelID, m.GuildID)
	} else if m.Content == "!set channel browse" && !common.MemberHasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionAdministrator) {
		s.ChannelMessageSend(m.ChannelID, "Sorry, you are not a server administrator!")
	}
}

func FindAGameStats_Time() string {
	return fmt.Sprintf("<t:%d:R>", time.Now().Unix())
}

func FindAGameStats_Message() string {
	dragon_runs, dragon_carries := orm.GetFindAGameType("dragon")
	kraken_runs, kraken_carries := orm.GetFindAGameType("kraken")
	yeti_runs, yeti_carries := orm.GetFindAGameType("yeti")
	maze_runs, maze_carries := orm.GetFindAGameType("maze")
	abyssal_runs, abyssal_carries := orm.GetFindAGameType("abyssal")
	members_count := orm.GetMembersCount()
	message_template := "**Requests Board**\n" +
		"Current open requests count for dungeon runs and carries:\n" +
		"**Players registered:** " + string(members_count) +
		"\n\n" +
		"**Dragon:** " + fmt.Sprintf("Runs: %d, Carries: %d", dragon_runs, dragon_carries) +
		"\n**Kraken:** " + fmt.Sprintf("Runs: %d, Carries: %d", kraken_runs, kraken_carries) +
		"\n**Yeti:** " + fmt.Sprintf("Runs: %d, Carries: %d", yeti_runs, yeti_carries) +
		"\n**Maze:** " + fmt.Sprintf("Runs: %d, Carries: %d", maze_runs, maze_carries) +
		"\n**Abyssal:** " + fmt.Sprintf("Runs: %d, Carries: %d", abyssal_runs, abyssal_carries)
	return message_template
}

func FindAGameStats_EmbedMessage() []*discordgo.MessageEmbed {
	dragon_runs, dragon_carries := orm.GetFindAGameType("dragon")
	kraken_runs, kraken_carries := orm.GetFindAGameType("kraken")
	yeti_runs, yeti_carries := orm.GetFindAGameType("yeti")
	maze_runs, maze_carries := orm.GetFindAGameType("maze")
	abyssal_runs, abyssal_carries := orm.GetFindAGameType("abyssal")
	members_count := orm.GetMembersCount()
	message_embed := []*discordgo.MessageEmbed{
		{
			Title: "Requests Board",
			Description: "Current open requests count for dungeon runs and carries:" +
				"\n\n**Dragon:** " + fmt.Sprintf("Runs: %d, Carries: %d", dragon_runs, dragon_carries) +
				"\n**Kraken:** " + fmt.Sprintf("Runs: %d, Carries: %d", kraken_runs, kraken_carries) +
				"\n**Yeti:** " + fmt.Sprintf("Runs: %d, Carries: %d", yeti_runs, yeti_carries) +
				"\n**Maze:** " + fmt.Sprintf("Runs: %d, Carries: %d", maze_runs, maze_carries) +
				"\n**Abyssal:** " + fmt.Sprintf("Runs: %d, Carries: %d", abyssal_runs, abyssal_carries),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Players registered",
					Value: members_count,
				},
				{},
			},
			Color: 0x500497,
		},
	}
	return message_embed
}

func FindAGameStats(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!dfstats" && common.MemberHasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionAdministrator) {
		var board_message *discordgo.Message
		var board_md5 [16]byte

		embed_message := FindAGameStats_EmbedMessage()
		embed_message[0].Fields[1] = &discordgo.MessageEmbedField{
			Name:  "Last updated:",
			Value: string(FindAGameStats_Time()),
		}

		board_message, _ = s.ChannelMessageSendEmbeds(
			m.ChannelID,
			embed_message,
		)

		s.ChannelMessageDelete(m.ChannelID, m.ID)

		board_md5 = md5.Sum([]byte(FindAGameStats_Message()))

		go func() {
			for {
				time.Sleep(60 * time.Second)

				// check if our message still exists
				_, err := s.ChannelMessage(m.ChannelID, board_message.ID)
				if err != nil {
					// break out of the for loop
					break
				}

				if board_md5 != md5.Sum([]byte(FindAGameStats_Message())) {
					embed_message := FindAGameStats_EmbedMessage()
					embed_message[0].Fields[1] = &discordgo.MessageEmbedField{
						Name:  "Last updated:",
						Value: string(FindAGameStats_Time()),
					}
					board_message, err = s.ChannelMessageEditEmbeds(
						m.ChannelID,
						board_message.ID,
						embed_message,
					)
					if err != nil {
						break
					}
				}
				board_md5 = md5.Sum([]byte(FindAGameStats_Message()))
			}
		}()
	}
}

var game_types = map[string]string{
	"🐉": "Dragon's Dungeon",
	"🐙": "Kraken's Ship",
	"⛄": "Yeti's Tundra",
	"🏰": "Maze",
	"😈": "Abyssal Maze",
	"💬": "Event",
}

var game_name = map[string]string{
	"🐉": "dragon",
	"🐙": "kraken",
	"⛄": "yeti",
	"🏰": "maze",
	"😈": "abyssal",
	"💬": "event",
}

func ReactToFindAGame(s *discordgo.Session, member_id string, member_id_poster string, guild *orm.Guilds, game_type string) {
	guild_config := orm.GetGuildConfig(guild.GuildID)
	player_id, _ := orm.GetMemberWithPlayerID(member_id)
	if player_id.PlayerID == "" {
		noid_msg, err := s.ChannelMessageSend(guild_config.ChannelBrowse, fmt.Sprintf("<@%s> Sorry, but you need to set your Player ID first! <#%s>", member_id, guild_config.ChannelPlayerID))
		if err != nil {
			return
		}
		timer := time.NewTimer(60 * time.Second)
		go func() {
			<-timer.C
			err := s.ChannelMessageDelete(guild_config.ChannelBrowse, noid_msg.ID)
			if err != nil {
				return
			}
		}()
		return
	}
	priv_chan, _ := s.UserChannelCreate(member_id_poster)
	s.ChannelMessageSend(priv_chan.ID, fmt.Sprintf("Message from server %s:\n<@%s> wants to play a Hunt Royale **%s** with you! :id: %s", guild.GuildName, member_id, game_type, player_id.PlayerID))
	// go func() {
	// 	time.Sleep(game_request_timeout * time.Minute)
	// 	s.ChannelMessageEdit(priv_chan.ID, react_msg.ID, fmt.Sprintf("Message from server %s:\nHunt Royale Dungeon Finder message has expired.", guild.GuildName))
	// }()
}

func FindAGameEmojiResponse(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	guild_info := orm.GetGuild(r.GuildID)
	// this means we are in the find a game channel, and we care what is being reacted
	// skip if the bot is reacting
	if r.UserID == s.State.User.ID {
		return
	}

	// check if we care about this emoji at all
	allowedemojis := []string{"🐉", "🐙", "⛄", "🏰", "😈", "💬"}
	if !slices.Contains(allowedemojis, r.Emoji.Name) {
		return
	}

	// check if user already responded to the message
	reaction := orm.GetFindAGameReaction(r.MessageID, r.UserID, r.GuildID)
	if !reaction.CreatedAt.IsZero() && !time.Now().After(reaction.CreatedAt.Add(30*time.Minute)) {
		msgref := &discordgo.MessageReference{
			ChannelID: r.ChannelID,
			MessageID: r.MessageID,
			GuildID:   r.GuildID,
		}
		react_msg, err := s.ChannelMessageSendReply(r.ChannelID, fmt.Sprintf("<@%s> you have already reacted to this game request!", r.UserID), msgref)
		if err != nil {
			return
		}
		timer := time.NewTimer(30 * time.Second)
		go func() {
			<-timer.C
			err := s.ChannelMessageDelete(r.ChannelID, react_msg.ID)
			if err != nil {
				return
			}
		}()
		return
	}

	// prevent user responding to their own request
	fag := orm.GetFindAGameByMsgID(r.MessageID, r.GuildID)
	if r.UserID == fag.UserID {
		msgref := &discordgo.MessageReference{
			ChannelID: r.ChannelID,
			MessageID: r.MessageID,
			GuildID:   r.GuildID,
		}
		react_msg, err := s.ChannelMessageSendReply(r.ChannelID, fmt.Sprintf("<@%s> you cannot respond to your own request!", r.UserID), msgref)
		if err != nil {
			return
		}
		timer := time.NewTimer(30 * time.Second)
		go func() {
			<-timer.C
			err := s.ChannelMessageDelete(r.ChannelID, react_msg.ID)
			if err != nil {
				return
			}
		}()
		return
	}
	ReactToFindAGame(s, r.UserID, fag.UserID, guild_info, game_types[r.Emoji.Name])
	orm.AddFindAGameReaction(r.MessageID, r.GuildID, fag.UserID, game_name[r.Emoji.Name])
}

// func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	if m.Author.ID == s.State.User.ID {
// 		return
// 	}

// 	// Do not show the notification on the #find-a-game channel
// 	if m.ChannelID == channel_findagame {
// 		return
// 	}

// 	if m.MentionRoles != nil {
// 		for _, role := range m.MentionRoles {
// 			// if slices.Contains()
// 			s.ChannelMessageSend(m.ChannelID, "You mentioned "+roles[role]+" but you're not in #find-a-game!")
// 		}
// 	}

// 	// fmt.Println(m.Message)
// }
