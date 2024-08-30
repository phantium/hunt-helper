package handlers

import (
	"crypto/md5"
	"discordbot/internal/common"
	"discordbot/internal/orm"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	dungeons := []string{"dragon", "kraken", "yeti", "maze", "abyssal"}
	var stats []string

	caser := cases.Title(language.English)
	for _, dungeon := range dungeons {
		runs, carries, carry_offers := orm.GetFindAGameType(dungeon)
		stats = append(stats, fmt.Sprintf("**%s:** Runs: %d, Carry Requests: %d, Carry Offers: %d",
			caser.String(dungeon), runs, carries, carry_offers))
	}

	members_count := orm.GetMembersCount()

	message_template := fmt.Sprintf(`**Requests Board**
Current open requests count for dungeon runs and carries:
**Players registered:** %s

%s`, members_count, strings.Join(stats, "\n"))

	return message_template
}

func FindAGameStats_EmbedMessage() []*discordgo.MessageEmbed {
	dragon_runs, dragon_carries, dragon_carry_offers := orm.GetFindAGameType("dragon")
	kraken_runs, kraken_carries, kraken_carry_offers := orm.GetFindAGameType("kraken")
	yeti_runs, yeti_carries, yeti_carry_offers := orm.GetFindAGameType("yeti")
	maze_runs, maze_carries, maze_carry_offers := orm.GetFindAGameType("maze")
	abyssal_runs, abyssal_carries, abyssal_carry_offers := orm.GetFindAGameType("abyssal")
	event_runs, _, _ := orm.GetFindAGameType("event")
	coop_runs, _, _ := orm.GetFindAGameType("coop")
	members_count := orm.GetMembersCount()
	message_embed := []*discordgo.MessageEmbed{
		{
			Title: ":scroll: Requests Board",
			Description: "Current open requests count for dungeon runs and carries:" +
				"\n\n:dragon: **Dragon:** " + fmt.Sprintf("Runs: %d, Carry Requests: %d, Carry Offers: %d", dragon_runs, dragon_carries, dragon_carry_offers) +
				"\n:octopus: **Kraken:** " + fmt.Sprintf("Runs: %d, Carry Requests: %d, Carry Offers: %d", kraken_runs, kraken_carries, kraken_carry_offers) +
				"\n:snowman: **Yeti:** " + fmt.Sprintf("Runs: %d, Carry Requests: %d, Carry Offers: %d", yeti_runs, yeti_carries, yeti_carry_offers) +
				"\n:european_castle: **Maze:** " + fmt.Sprintf("Runs: %d, Carry Requests: %d, Carry Offers: %d", maze_runs, maze_carries, maze_carry_offers) +
				"\n:smiling_imp: **Abyssal:** " + fmt.Sprintf("Runs: %d, Carry Requests: %d, Carry Offers: %d", abyssal_runs, abyssal_carries, abyssal_carry_offers) +
				"\n:speech_balloon: **Event:** " + fmt.Sprintf("Runs: %d", event_runs) +
				"\n:busts_in_silhouette: **Co-Op:** " + fmt.Sprintf("Runs: %d", coop_runs),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  ":video_game: Players registered",
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
			Name:  ":clock1: Last updated:",
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
						Name:  ":clock1: Last updated:",
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

func FindAGameStatsPoster(s *discordgo.Session) {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {

		for _, guild := range s.State.Guilds {
			var board_message *discordgo.Message

			// get board channel id
			guild_config := orm.GetGuildConfig(guild.ID)
			channelID := guild_config.ChannelBoard
			messageID := guild_config.ChannelBoardPost
			if channelID == "" {
				// no channelID configured
				continue
			}

			if messageID == "" {
				embed_message := FindAGameStats_EmbedMessage()
				embed_message[0].Fields[1] = &discordgo.MessageEmbedField{
					Name:  ":clock1: Last updated:",
					Value: string(FindAGameStats_Time()),
				}

				board_message, err := s.ChannelMessageSendEmbeds(
					channelID,
					embed_message,
				)
				if err != nil {
					continue
				}
				// board_md5 = md5.Sum([]byte(FindAGameStats_Message()))
				orm.UpdateGuildConfig(guild.ID, map[string]interface{}{"channel_board_post": board_message.ID})
			} else {
				// check if our message exists
				chan_msg, err := s.ChannelMessage(channelID, messageID)
				if err != nil {
					embed_message := FindAGameStats_EmbedMessage()
					embed_message[0].Fields[1] = &discordgo.MessageEmbedField{
						Name:  ":clock1: Last updated:",
						Value: string(FindAGameStats_Time()),
					}

					board_message, err = s.ChannelMessageSendEmbeds(
						channelID,
						embed_message,
					)
					if err != nil {
						continue
					}
					// board_md5 = md5.Sum([]byte(FindAGameStats_Message()))
					orm.UpdateGuildConfig(guild.ID, map[string]interface{}{"channel_board_post": board_message.ID})
				} else {
					current_msg := []byte(chan_msg.Embeds[0].Description)
					current_msg_fields := []byte(chan_msg.Embeds[0].Fields[0].Value)
					current_msg_combined := append(current_msg, current_msg_fields...)
					new_msg := []byte(FindAGameStats_EmbedMessage()[0].Description)
					new_msg_fields := []byte(FindAGameStats_EmbedMessage()[0].Fields[0].Value)
					new_msg_combined := append(new_msg, new_msg_fields...)

					if md5.Sum(current_msg_combined) != md5.Sum(new_msg_combined) {
						embed_message := FindAGameStats_EmbedMessage()
						embed_message[0].Fields[1] = &discordgo.MessageEmbedField{
							Name:  ":clock1: Last updated:",
							Value: string(FindAGameStats_Time()),
						}
						_, err := s.ChannelMessageEditEmbeds(
							channelID,
							messageID,
							embed_message,
						)
						if err != nil {
							log.Error(err)
							continue
						}
						// orm.UpdateGuildConfig(guild.ID, map[string]interface{}{"channel_board_post": board_message.ID})
					}
				}
			}

		}
	}
}

var game_types = map[string]string{
	"🐉": "Dragon's Dungeon",
	"🐙": "Kraken's Ship",
	"⛄": "Yeti's Tundra",
	"🏰": "Maze",
	"😈": "Abyssal Maze",
	"":  "Chaos Dungeon",
	"💬": "Event",
	"👥": "Co-Op",
}

var game_name = map[string]string{
	"🐉": "dragon",
	"🐙": "kraken",
	"⛄": "yeti",
	"🏰": "maze",
	"😈": "abyssal",
	"":  "chaos",
	"💬": "event",
	"👥": "coop",
}

func deleteMessage(s *discordgo.Session, msg *discordgo.Message) {
	err := s.ChannelMessageDelete(msg.ChannelID, msg.ID)
	if err != nil {
		return
	}
}

func deleteMessageAfterTimeout(s *discordgo.Session, msg *discordgo.Message, request_timeout time.Duration) {
	time.AfterFunc(request_timeout, func() {
		deleteMessage(s, msg)
	})
}

func ReactToFindAGame(s *discordgo.Session, member_id string, member_id_poster string, guild *orm.Guilds, channel_id string, game_type string) {
	guild_config := orm.GetGuildConfig(guild.GuildID)
	player_id, _ := orm.GetMemberWithPlayerID(member_id)
	if player_id.PlayerID == "" {
		if guild_config.ChannelPlayerID == "" {
			noid_msg, err := s.ChannelMessageSend(channel_id, fmt.Sprintf("<@%s> Sorry, but you need to set your Player ID first! However, the server admin has not yet configured the channel!", member_id))
			if err != nil {
				return
			}
			deleteMessageAfterTimeout(s, noid_msg, 30*time.Second)
			return
		}
		noid_msg, err := s.ChannelMessageSend(channel_id, fmt.Sprintf("<@%s> Sorry, but you need to set your Player ID first! <#%s>", member_id, guild_config.ChannelPlayerID))
		if err != nil {
			return
		}
		deleteMessageAfterTimeout(s, noid_msg, 30*time.Second)
		return
	}
	priv_chan, err := s.UserChannelCreate(member_id_poster)
	if err != nil {
		// user probably has dms disabled
		return
	}
	msg, err := s.ChannelMessageSend(priv_chan.ID, fmt.Sprintf("Message from server %s:\n<@%s> wants to play a Hunt Royale **%s** with you! :id: %s", guild.GuildName, member_id, game_type, player_id.PlayerID))
	if err != nil {
		return
	}
	deleteMessageAfterTimeout(s, msg, 10*time.Minute)
}

func FindAGameEmojiResponse(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	guild_info := orm.GetGuild(r.GuildID)
	// this means we are in the find a game channel, and we care what is being reacted

	// skip if the bot is reacting
	if r.UserID == s.State.User.ID {
		return
	}

	// check if we care about this emoji at all
	allowedemojis := []string{"🐉", "🐙", "⛄", "🏰", "😈", "💬", "👥"}
	if !slices.Contains(allowedemojis, r.Emoji.Name) {
		return
	}

	// check if user already responded to the message
	reaction, err := orm.GetFindAGameReaction(r.UserID, r.GuildID, r.MessageID)
	if err != nil {
		log.Error(err)
	}
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
		deleteMessageAfterTimeout(s, react_msg, 30*time.Second)
		return
	}

	// prevent user responding to their own request
	fag, err := orm.GetFindAGameByMsgID(r.MessageID)
	if err != nil {
		log.Error(err)
	}
	if r.UserID == fag.UserID {
		return
	}
	ReactToFindAGame(s, r.UserID, fag.UserID, guild_info, r.ChannelID, game_types[r.Emoji.Name])
	orm.AddFindAGameReaction(r.MessageID, r.GuildID, r.UserID, game_name[r.Emoji.Name])
}
