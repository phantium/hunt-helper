package orm

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/exp/slices"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"discordbot/internal/configuration"
)

type PlayerIDs struct {
	gorm.Model
	MemberID string `gorm:"unique"`
	PlayerID string
}

// type DungeonFinders struct {
// 	gorm.Model
// 	MemberID string `gorm:"unique"`
// 	Dungeons string
// }

type BoardPost struct {
	gorm.Model
	MessageID string `gorm:"unique"`
	ChannelID string
	GuildID   string
}

type GuildConfig struct {
	gorm.Model
	GuildID               string `gorm:"unique"`
	ChannelFindAGame      string
	ChannelMultiple       bool `gorm:"default:true"`
	ChannelPlayerID       string
	ChannelBoard          string
	ChannelBoardPost      string
	FAGRequestTime        int `gorm:"default:60"`
	FAGRequestTimeout     int `gorm:"default:60"`
	FAGDungeonSelectLimit int `gorm:"default:3"`
	RoleDragon            string
	RoleKraken            string
	RoleYeti              string
	RoleMaze              string
	RoleAbyssal           string
	RoleCoop              string
	RoleEvent             string
	ChannelDragon         string
	ChannelKraken         string
	ChannelYeti           string
	ChannelMaze           string
	ChannelAbyssal        string
	ChannelCoop           string
	ChannelEvent          string
}

type Guilds struct {
	gorm.Model
	GuildID         string `gorm:"unique"`
	GuildName       string
	OwnerID         string
	SystemChannelID string
}

type Dungeons []string

type FindAGameMessage struct {
	gorm.Model
	MessageID string
	ChannelID string
	GuildID   string
	UserID    string
	RunType   string
	Dungeons  Dungeons `gorm:"serializer:json"`
}

type FindAGameReaction struct {
	gorm.Model
	MessageID string
	GuildID   string
	UserID    string
	Dungeon   string
}

var _db *gorm.DB

var cfg configuration.DiscordConfig

const discord_config string = "discord.yml"

func init() {
	// as init is standalone, it cannot rely on being given a cfg interface
	configuration.ReadConfig(&cfg, discord_config)

	username := cfg.Database.Username
	password := cfg.Database.Password
	host := cfg.Database.Host
	port := cfg.Database.Port
	dbname := cfg.Database.Database
	timeout := cfg.Database.Timeout
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, dbname, timeout)

	var err error

	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Fatal("failed to connect to database, error: ", err)
	}

	// migrate our model if necessary
	_db.AutoMigrate(&PlayerIDs{})
	_db.AutoMigrate(&FindAGameReaction{})
	_db.AutoMigrate(&GuildConfig{})
	_db.AutoMigrate(&Guilds{})
	_db.AutoMigrate(&FindAGameMessage{})
	_db.AutoMigrate(&BoardPost{})

	sqlDB, _ := _db.DB()
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

// member management

func DelMembersExistingPlayerID(member_id string) {
	_db.Unscoped().Where("member_id = ?", member_id).Delete(&PlayerIDs{})
}

func AddMemberWithPlayerID(member_id string, player_id string) {
	_db.Create(&PlayerIDs{
		MemberID: member_id,
		PlayerID: player_id,
	})
}

func GetMemberWithPlayerID(member_id string) (*PlayerIDs, error) {
	var playerids *PlayerIDs
	if err := _db.Where("member_id = ?", member_id).Find(&playerids).Error; err != nil {
		return nil, err
	}
	return playerids, nil
}

func GetMemberIDWithPlayerID(player_id string) (*PlayerIDs, error) {
	var playerids *PlayerIDs
	if err := _db.Where("player_id = ?", player_id).Find(&playerids).Error; err != nil {
		return nil, err
	}
	return playerids, nil
}

func GetMembersCount() string {
	var playerids []*PlayerIDs
	result := _db.Find(&playerids)
	return fmt.Sprint(result.RowsAffected)
}

// guild registration

func GetGuild(guild_id string) *Guilds {
	var guild *Guilds
	if err := _db.Where("guild_id = ?", guild_id).Find(&guild).Error; err != nil {
		return nil
	}
	return guild
}

func GetGuilds() *[]Guilds {
	var guild *[]Guilds
	if err := _db.Find(&guild).Error; err != nil {
		return nil
	}
	return guild
}

func CreateGuild(g *Guilds) {
	_db.Create(&Guilds{
		GuildID:         g.GuildID,
		GuildName:       g.GuildName,
		OwnerID:         g.OwnerID,
		SystemChannelID: g.SystemChannelID,
	})
}

func DeleteGuild(guild_id string) {
	_db.Unscoped().Where("guild_id = ?", guild_id).Delete(&Guilds{})
}

// guild management

func GetGuildConfig(guild_id string) *GuildConfig {
	var guildconfig *GuildConfig
	if err := _db.Where("guild_id = ?", guild_id).Find(&guildconfig).Error; err != nil {
		return nil
	}
	return guildconfig
}

func CreateGuildConfig(guild_id string) {
	_db.Create(&GuildConfig{
		GuildID: guild_id,
	})
}

func UpdateGuildConfig(guild_id string, config map[string]interface{}) {
	if config != nil {
		_db.Model(&GuildConfig{}).Where("guild_id = ?", guild_id).Updates(config)
	}
}

// board management

func CreateBoardPost(bp *BoardPost) {
	_db.Create(&BoardPost{
		MessageID: bp.MessageID,
		ChannelID: bp.ChannelID,
		GuildID:   bp.GuildID,
	})
}

func GetBoardPost(bp *BoardPost) *BoardPost {
	var boardpost *BoardPost
	if err := _db.Where("message_id = ?", bp.MessageID).Where("guild_id = ?", bp.GuildID).Find(&boardpost).Error; err != nil {
		log.Println(err)
		return boardpost
	}
	return boardpost
}

// find a game reaction

func DeleteFindAGameReactions() {
	_db.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&FindAGameReaction{})
}

func AddFindAGameReaction(message_id string, guild_id string, user_id string, dungeon string) {
	_db.Create(&FindAGameReaction{
		MessageID: message_id,
		GuildID:   guild_id,
		UserID:    user_id,
		Dungeon:   dungeon,
	})
}

func GetFindAGameReaction(user_id string, guild_id string, message_id string) (*FindAGameReaction, error) {
	var reaction *FindAGameReaction
	if err := _db.Where("user_id = ?", user_id).Where("guild_id = ?", guild_id).Where("message_id = ?", message_id).Last(&reaction).Error; err != nil {
		// if err := _db.Where("user_id = ?", user_id).Last(&reaction).Error; err != nil {
		return reaction, err
	}
	return reaction, nil
}

// find a game

func AddFindAGame(message_id string, channel_id string, guild_id string, user_id string, dungeons []string, run_type string) {
	_db.Create(&FindAGameMessage{
		MessageID: message_id,
		ChannelID: channel_id,
		GuildID:   guild_id,
		UserID:    user_id,
		Dungeons:  dungeons,
		RunType:   run_type,
	})
}

func DeleteFindAGame(user_id string, guild_id string) {
	_db.Unscoped().Where("user_id = ?", user_id).Where("guild_id = ?", guild_id).Delete(&FindAGameMessage{})
}

func DeleteFindAGameByMessageID(message_id string) {
	_db.Unscoped().Where("message_id = ?", message_id).Delete(&FindAGameMessage{})
}

func GetFindAGame(user_id string) *FindAGameMessage {
	var findagame *FindAGameMessage
	if err := _db.Where("user_id = ?", user_id).Last(&findagame).Error; err != nil {
		log.Println(err)
		return findagame
	}
	return findagame
}

func GetFindAGameByMsgID(message_id string) (*FindAGameMessage, error) {
	var findagame *FindAGameMessage
	// if err := _db.Where("message_id = ?", message_id).Where("guild_id = ?", guild_id).Last(&findagame).Error; err != nil {
	if err := _db.Where("message_id = ?", message_id).Last(&findagame).Error; err != nil {
		return findagame, err
	}
	return findagame, nil
}

func GetFindAGameByUserID(user_id string) (*FindAGameMessage, error) {
	var findagame *FindAGameMessage
	// if err := _db.Where("message_id = ?", message_id).Where("guild_id = ?", guild_id).Last(&findagame).Error; err != nil {
	if err := _db.Where("user_id = ?", user_id).Last(&findagame).Error; err != nil {
		return findagame, err
	}
	return findagame, nil
}

func GetFindAGameType(dungeon string) (int, int, int) {
	// returns counter_run, counter_carry
	var findagame []*FindAGameMessage
	var user_ids []string
	var counter_run int = 0
	var counter_carry int = 0
	var counter_carry_offers int = 0
	if err := _db.Find(&findagame).Error; err != nil {
		return 0, 0, 0
	}
	for _, r := range findagame {
		if slices.Contains(r.Dungeons, dungeon) && !slices.Contains(user_ids, r.UserID) {
			user_ids = append(user_ids, r.UserID)
			if r.RunType == "run" {
				counter_run += 1
			}
			if r.RunType == "carry" {
				counter_carry += 1
			}
			if r.RunType == "carry_offer" {
				counter_carry_offers += 1
			}
		}
	}
	return counter_run, counter_carry, counter_carry_offers
}

// common functions
func GetPlayerID(member_id string) string {
	player_id, _ := GetMemberWithPlayerID(member_id)
	return player_id.PlayerID
}

func GetMemberID(player_id string) string {
	member_id, _ := GetMemberIDWithPlayerID(player_id)
	return member_id.MemberID
}
