package handlers

import (
	"discordbot/internal/orm"

	"github.com/bwmarrin/discordgo"

	log "github.com/sirupsen/logrus"
)

// executed when added into a new server (guild)

func GuildCreate(s *discordgo.Session, gc *discordgo.GuildCreate) {
	guild := orm.GetGuild(gc.Guild.ID)
	if guild.GuildID == "" {
		orm.CreateGuild(&orm.Guilds{
			GuildID:   gc.Guild.ID,
			GuildName: gc.Guild.Name,
			OwnerID:   gc.Guild.OwnerID,
		})
		guildconfig := orm.GetGuildConfig(gc.Guild.ID)
		if guildconfig.GuildID == "" {
			orm.CreateGuildConfig(gc.Guild.ID)
		}
	}
	log.Infof("Joined server with guild ID: %s (Name: %s) (Owner: %s)!", gc.Guild.ID, gc.Guild.Name, gc.Guild.OwnerID)

	// dungeon_finder command
	_, err := s.ApplicationCommandCreate(s.State.Application.ID, gc.Guild.ID, &discordgo.ApplicationCommand{
		Name:        "dungeon_finder",
		Description: "Hunt Royale Dungeon Finder",
	})
	if err != nil {
		log.Fatalf("Cannot create slash command /dungeon_finder: %v", err)
	}

	// defer s.ApplicationCommandDelete(s.State.Application.ID, gc.Guild.ID, df.ID)

	log.Infof("Created command /dungeon_finder on server with guild ID: %s (Name: %s) (Owner: %s)!", gc.Guild.ID, gc.Guild.Name, gc.Guild.OwnerID)

	// ef := []*discordgo.MessageEmbedField{
	// 	{
	// 		Name:   "Greetings!",
	// 		Value:  "I am " + s.State.Application.Name + " and I am here to help with Hunt Royale!",
	// 		Inline: true,
	// 	},
	// }
	// e := &discordgo.MessageEmbed{
	// 	Type:   discordgo.EmbedTypeRich,
	// 	Title:  s.State.Application.Name + " has arrived!",
	// 	Fields: ef,
	// }
	// if _, err := s.ChannelMessageSendEmbed(gc.Guild.SystemChannelID, e); err != nil {
	// 	log.Errorf("Failed to send introduction message for %s (GuildID: %s)", gc.Guild.Name, gc.Guild.ID)
	// }
}
