package slashcommands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type BotConfig struct {
	ken.EphemeralCommand
}

var (
	_ ken.SlashCommand = (*BotConfig)(nil)
	_ ken.DmCapable    = (*BotConfig)(nil)
)

func (c *BotConfig) Name() string {
	return "hh_config"
}

func (c *BotConfig) Description() string {
	return "Hunt Helper Config"
}

func (c *BotConfig) Version() string {
	return "1.0.0"
}

func (c *BotConfig) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *BotConfig) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *BotConfig) IsDmCapable() bool {
	return false
}

func (c *BotConfig) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	b := ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "Hunt Helper guild configuration",
	})

	channel_select := ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "Setup your channels",
	})
	channel_select.AddComponents(func(cb *ken.ComponentBuilder) {
		cb.AddActionsRow(func(b ken.ComponentAssembler) {
			b.Add(discordgo.Button{
				Label:    "Use Single LFG Channel",
				CustomID: "config_single_lfg",
				Style:    discordgo.PrimaryButton,
			}, func(ctx ken.ComponentContext) bool {
				return true
			})
			b.Add(discordgo.Button{
				Label:    "Use Multiple LFG Channels",
				CustomID: "config_multiple_lfg",
				Style:    discordgo.DangerButton,
			}, func(ctx ken.ComponentContext) bool {
				return true
			})
		})
	})

	b.AddComponents(func(cb *ken.ComponentBuilder) {
		cb.AddActionsRow(func(b ken.ComponentAssembler) {
			b.Add(discordgo.Button{
				CustomID: "config_channels",
				Label:    "Configure channels",
			}, func(ctx ken.ComponentContext) bool {
				channel_select.Send()
				return true
			})
			b.Add(discordgo.Button{
				CustomID: "requests_board",
				Label:    "Requests board",
			}, func(ctx ken.ComponentContext) bool {
				return true
			})

		})
	})

	fum := b.Send()
	return fum.Error
}
