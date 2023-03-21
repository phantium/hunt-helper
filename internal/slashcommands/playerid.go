package slashcommands

import (
	"discordbot/internal/orm"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/ken/middlewares/ratelimit"
)

type CommandPlayerID struct {
	ken.EphemeralCommand
}

var (
	_ ken.SlashCommand         = (*CommandPlayerID)(nil)
	_ ken.DmCapable            = (*CommandPlayerID)(nil)
	_ ratelimit.LimitedCommand = (*CommandPlayerID)(nil)
)

func (c *CommandPlayerID) Name() string {
	return "player_id"
}

func (c *CommandPlayerID) Description() string {
	return "Hunt Royale Fetch Player ID"
}

func (c *CommandPlayerID) Version() string {
	return "1.0.0"
}

func (c *CommandPlayerID) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *CommandPlayerID) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Required:    true,
			Description: "@ a discord user",
		},
	}
}

func (c *CommandPlayerID) LimiterBurst() int {
	return 4
}

func (c *CommandPlayerID) LimiterRestoration() time.Duration {
	return 30 * time.Second
}

func (c *CommandPlayerID) IsLimiterGlobal() bool {
	return false
}

func (c *CommandPlayerID) IsDmCapable() bool {
	return true
}

func (c *CommandPlayerID) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	var playerid_message string

	val := ctx.Options().GetByName("user").UserValue(ctx)

	player_id := orm.GetPlayerID(val.ID)
	if player_id != "" {
		playerid_message = fmt.Sprintf("<@%s> Hunt Royale :id: %s", val.ID, player_id)
	} else {
		playerid_message = fmt.Sprintf("<@%s> has no Hunt Royale :id: registered", val.ID)
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: playerid_message,
		},
	}

	ctx.Respond(response)
	return
}
