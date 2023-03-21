package slashcommands

// import (
// 	"discordbot/internal/common"
// 	"fmt"
// 	"time"

// 	"github.com/bwmarrin/discordgo"
// 	"github.com/zekrotja/ken"
// )

// type DungeonFinder struct{}

// var (
// 	_ ken.SlashCommand = (*DungeonFinder)(nil)
// )

// func (c *DungeonFinder) Name() string {
// 	return "dungeon_finder"
// }

// func (c *DungeonFinder) Description() string {
// 	return "Hunt Royale Dungeon Finder"
// }

// func (c *DungeonFinder) Version() string {
// 	return "1.0.0"
// }

// func (c *DungeonFinder) Type() discordgo.ApplicationCommandType {
// 	return discordgo.ChatApplicationCommand
// }

// func (c *DungeonFinder) Options() []*discordgo.ApplicationCommandOption {
// 	return []*discordgo.ApplicationCommandOption{}
// 	// options := []*discordgo.ApplicationCommandOption{
// 	// 	{
// 	// 		Type:        discordgo.ApplicationCommandOptionSubCommand,
// 	// 		Name:        "dungeon_finder",
// 	// 		Description: "Create a message with attached role select buttons.",
// 	// 		Options: append([]*discordgo.ApplicationCommandOption{{
// 	// 			Type:        discordgo.ApplicationCommandOptionString,
// 	// 			Name:        "content",
// 	// 			Description: "The content of the message.",
// 	// 			Required:    true,
// 	// 		}}),
// 	// 	},
// 	// }
// 	// return options
// }

// func (c *DungeonFinder) delete_emb(*discordgo.MessageEmbed) {

// }

// func (c *DungeonFinder) Run(ctx ken.Context) (err error) {
// 	if err = ctx.Defer(); err != nil {
// 		return
// 	}

// 	// checking the required perms to run the command!
// 	if !common.MemberHasPermission(ctx.GetSession(), ctx.GetEvent().GuildID, ctx.GetEvent().Member.User.ID, discordgo.PermissionAdministrator) {
// 		respond := &discordgo.InteractionResponse{
// 			Type: discordgo.InteractionResponseChannelMessageWithSource,
// 			Data: &discordgo.InteractionResponseData{
// 				Content: "Sorry, you are not allowed to use this command",
// 				Flags:   discordgo.MessageFlagsEphemeral,
// 			},
// 		}
// 		ctx.Respond(respond)
// 		return
// 	}

// 	// we want at least one selection
// 	var selectMenuMinValues int = 1

// 	b := ctx.FollowUpEmbed(&discordgo.MessageEmbed{
// 		Description: "Use the dungeon finder to find a dungeon.",
// 	})

// 	b.AddComponents(func(cb *ken.ComponentBuilder) {
// 		cb.AddActionsRow(func(b ken.ComponentAssembler) {
// 			b.Add(discordgo.SelectMenu{
// 				CustomID:    "dungeon_finder",
// 				Placeholder: "Choose your dungeons for a run:",
// 				// This is where confusion comes from. If you don't specify these things you will get single item select.
// 				// These fields control the minimum and maximum amount of selected items.
// 				MinValues: &selectMenuMinValues,
// 				MaxValues: 3,
// 				Options: []discordgo.SelectMenuOption{
// 					{
// 						Label: "Dragon",
// 						Value: "Dragon",
// 						// Default works the same for multi-select menus.
// 						Default: false,
// 						Emoji: discordgo.ComponentEmoji{
// 							Name: "dragon",
// 							ID:   "1082313506700935199",
// 						},
// 					},
// 					{
// 						Label: "Kraken",
// 						Value: "Kraken",
// 						Emoji: discordgo.ComponentEmoji{
// 							Name: "kraken",
// 							ID:   "1082313504901578822",
// 						},
// 					},
// 					{
// 						Label: "Yeti",
// 						Value: "Yeti",
// 						Emoji: discordgo.ComponentEmoji{
// 							Name: "yeti",
// 							ID:   "1082333118729556038",
// 						},
// 					},
// 					{
// 						Label: "Maze",
// 						Value: "Maze",
// 						Emoji: discordgo.ComponentEmoji{
// 							Name: "maze",
// 							ID:   "1082313502208827422",
// 						},
// 					},
// 					{
// 						Label: "Abyssal",
// 						Value: "Abyssal",
// 						Emoji: discordgo.ComponentEmoji{
// 							Name: "abyssal",
// 							ID:   "1082313499922944000",
// 						},
// 					},
// 				},
// 			}, func(ctx ken.ComponentContext) bool {
// 				ctx.RespondEmbed(&discordgo.MessageEmbed{
// 					Description: fmt.Sprintf("Resp %s to: %s", ctx.GetData().Values, ctx.GetData().CustomID),
// 				})
// 				return true
// 			})
// 		})
// 		cb.AddActionsRow(func(b ken.ComponentAssembler) {
// 			b.Add(discordgo.SelectMenu{
// 				CustomID:    "dungeon_finder_carry",
// 				Placeholder: "Choose your dungeons for a carry:",
// 				// This is where confusion comes from. If you don't specify these things you will get single item select.
// 				// These fields control the minimum and maximum amount of selected items.
// 				MinValues: &selectMenuMinValues,
// 				MaxValues: 3,
// 				Options: []discordgo.SelectMenuOption{
// 					{
// 						Label: "Dragon",
// 						Value: "Dragon",
// 						// Default works the same for multi-select menus.
// 						Default: false,
// 						Emoji: discordgo.ComponentEmoji{
// 							Name: "dragon",
// 							ID:   "1082313506700935199",
// 						},
// 					},
// 					{
// 						Label: "Kraken",
// 						Value: "Kraken",
// 						Emoji: discordgo.ComponentEmoji{
// 							Name: "kraken",
// 							ID:   "1082313504901578822",
// 						},
// 					},
// 					{
// 						Label: "Yeti",
// 						Value: "Yeti",
// 						Emoji: discordgo.ComponentEmoji{
// 							Name: "yeti",
// 							ID:   "1082333118729556038",
// 						},
// 					},
// 					{
// 						Label: "Maze",
// 						Value: "Maze",
// 						Emoji: discordgo.ComponentEmoji{
// 							Name: "maze",
// 							ID:   "1082313502208827422",
// 						},
// 					},
// 					{
// 						Label: "Abyssal",
// 						Value: "Abyssal",
// 						Emoji: discordgo.ComponentEmoji{
// 							Name: "abyssal",
// 							ID:   "1082313499922944000",
// 						},
// 					},
// 				},
// 			}, func(ctx ken.ComponentContext) bool {
// 				ctx.RespondEmbed(&discordgo.MessageEmbed{
// 					Description: fmt.Sprintf("Resp %s to: %s", ctx.GetData().Values, ctx.GetData().CustomID),
// 				})
// 				return true
// 			})
// 		})
// 	})

// 	fum := b.Send()

// 	ticker := time.NewTicker(30 * time.Second)
// 	quit := make(chan struct{})
// 	go func() {
// 		for {
// 			select {
// 			case <-ticker.C:
// 				fum.AddComponents()
// 			case <-quit:
// 				ticker.Stop()
// 				return
// 			}
// 		}
// 	}()
// 	return fum.Error
// }
