package slashcommands

import (
	"discordbot/internal/orm"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/ken/middlewares/ratelimit"
)

type CommandPlayerIDReverse struct {
	ken.EphemeralCommand
}

var (
	_ ken.SlashCommand         = (*CommandPlayerIDReverse)(nil)
	_ ken.DmCapable            = (*CommandPlayerIDReverse)(nil)
	_ ratelimit.LimitedCommand = (*CommandPlayerIDReverse)(nil)
)

func (c *CommandPlayerIDReverse) Name() string {
	return "player_id_reverse"
}

func (c *CommandPlayerIDReverse) Description() string {
	return "Hunt Royale get Discord member from Player ID"
}

func (c *CommandPlayerIDReverse) Version() string {
	return "1.0.0"
}

func (c *CommandPlayerIDReverse) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *CommandPlayerIDReverse) Options() []*discordgo.ApplicationCommandOption {
	var minLength int = 8
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "playerid",
			Required:    true,
			MinLength:   &minLength,
			MaxLength:   8,
			Description: "Hunt Royale player ID",
		},
	}
}

func (c *CommandPlayerIDReverse) LimiterBurst() int {
	return 4
}

func (c *CommandPlayerIDReverse) LimiterRestoration() time.Duration {
	return 30 * time.Second
}

func (c *CommandPlayerIDReverse) IsLimiterGlobal() bool {
	return false
}

func (c *CommandPlayerIDReverse) IsDmCapable() bool {
	return true
}

func (c *CommandPlayerIDReverse) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	var memberid_message string

	val := ctx.Options().GetByName("playerid").StringValue()

	member_id := orm.GetMemberID(val)
	if member_id != "" {
		memberid_message = fmt.Sprintf("Hunt Royale :id: %s belongs to <@%s>", val, member_id)
	} else {
		memberid_message = fmt.Sprintf("Hunt Royale :id: %s is not registered to anyone", val)
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: memberid_message,
		},
	}

	ctx.Respond(response)
	return
}
