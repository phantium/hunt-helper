package handlers

import (
	"discordbot/internal/orm"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

var game_types = map[string]string{
	"🐉": "Dragon's Dungeon",
	"🐙": "Kraken's Ship",
	"⛄": "Yeti's Tundra",
	"🏰": "Maze",
	"😈": "Abyssal Maze",
	"🎲": "Chaos Dungeon",
	"💬": "Event",
	"👥": "Co-Op",
	"🪵": "Dark Forest",
}

var game_name = map[string]string{
	"🐉": "dragon",
	"🐙": "kraken",
	"⛄": "yeti",
	"🏰": "maze",
	"😈": "abyssal",
	"🎲": "chaos",
	"💬": "event",
	"👥": "coop",
	"🪵": "darkforest",
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
	allowedemojis := []string{"🐉", "🐙", "⛄", "🏰", "😈", "💬", "👥", "🪵", "🎲"}
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
