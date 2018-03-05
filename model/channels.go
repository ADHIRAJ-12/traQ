package model

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/satori/go.uuid"
)

var channelPathMap = &sync.Map{}

// Channel :チャンネルの構造体
type Channel struct {
	ID        string    `xorm:"char(36) pk"`
	Name      string    `xorm:"varchar(20) not null unique(name_parent)"`
	ParentID  string    `xorm:"parent_id char(36) not null unique(name_parent)"`
	Topic     string    `xorm:"text"`
	IsForced  bool      `xorm:"bool not null"`
	IsDeleted bool      `xorm:"bool not null"`
	IsPublic  bool      `xorm:"bool not null"`
	IsVisible bool      `xorm:"bool not null"`
	CreatorID string    `xorm:"char(36) not null"`
	CreatedAt time.Time `xorm:"created not null"`
	UpdaterID string    `xorm:"char(36) not null"`
	UpdatedAt time.Time `xorm:"updated not null"`
}

// TableName テーブル名を指定するメソッド
func (channel *Channel) TableName() string {
	return "channels"
}

// Create チャンネル作成を行うメソッド
func (channel *Channel) Create() error {
	if channel.ID != "" {
		return fmt.Errorf("ID is not empty! You can use Update()")
	}

	if channel.Name == "" {
		return fmt.Errorf("name is empty")
	}

	if channel.CreatorID == "" {
		return fmt.Errorf("creatorID is empty")
	}

	if err := validateChannelName(channel.Name); err != nil {
		return err
	}

	channel.ID = CreateUUID()
	channel.IsVisible = true
	channel.UpdaterID = channel.CreatorID

	// ここまでで入力されない要素は初期値(""や0)で格納される
	if _, err := db.Insert(channel); err != nil {
		return fmt.Errorf("failed to create channel: %v", err)
	}

	//チャンネルパスをキャッシュ
	if path, err := channel.Path(); err == nil {
		channelPathMap.Store(uuid.FromStringOrNil(channel.ID), path)
	}

	return nil
}

// Exists 指定したチャンネルがuserIDのユーザーから見えるチャンネルかどうかを確認する
func (channel *Channel) Exists(userID string) (bool, error) {
	if userID != "" {
		has, err := db.Join("LEFT", "users_private_channels", "users_private_channels.channel_id = channels.id").Where("(is_public = true OR user_id = ?) AND is_deleted = false", userID).Get(channel)
		return has, err
	}
	return db.Get(channel)
}

// Update チャンネルの情報の更新を行う
func (channel *Channel) Update() error {
	_, err := db.ID(channel.ID).UseBool().Update(channel)
	if err != nil {
		return fmt.Errorf("failed to update channel: %v", err)
	}

	//チャンネルパスキャッシュの更新
	updateChannelPathWithDescendants(channel)

	return nil
}

// Parent 親チャンネルを取得する
func (channel *Channel) Parent() (*Channel, error) {
	if len(channel.ParentID) == 0 {
		return nil, nil
	}

	parent := &Channel{}
	has, err := db.Where("id = ?", channel.ParentID).Get(parent)
	if !has {
		return nil, fmt.Errorf("parent channel doesn't exist")
	}
	return parent, err
}

// Children userIDのユーザーから見えるchannelIDの子チャンネル
func (channel *Channel) Children(userID string) ([]string, error) {
	var channelIDList []string
	if channel.ID == "" {
		return nil, fmt.Errorf("channelID is empty")
	}
	err := db.Table("channels").Join("LEFT", "users_private_channels", "users_private_channels.channel_id = channels.id").Where("(is_public = true OR user_id = ?) AND parent_id = ? AND is_deleted = false", userID, channel.ID).Cols("id").Find(&channelIDList)

	if err != nil {
		return nil, fmt.Errorf("failed to find channels: %v", err)
	}
	return channelIDList, nil
}

// Path チャンネルのパス文字列を取得する
func (channel *Channel) Path() (string, error) {
	path := channel.Name
	current := channel

	for {
		parent, err := current.Parent()
		if err != nil {
			return "", err
		}
		if parent == nil {
			break
		}

		if parentPath, ok := GetChannelPath(uuid.FromStringOrNil(parent.ID)); ok {
			return parentPath + "/" + path, nil
		}

		path = parent.Name + "/" + path
		current = parent
	}

	return "#" + path, nil
}

// GetChannelByID チャンネルIDによってチャンネルを取得
func GetChannelByID(userID, channelID string) (*Channel, error) {
	channel := &Channel{}
	channel.ID = channelID
	has, err := db.Join("LEFT", "users_private_channels", "users_private_channels.channel_id = channels.id").Where("(is_public = true OR user_id = ?) AND is_deleted = false", userID).Get(channel)

	if err != nil {
		return nil, fmt.Errorf("failed to get channel: %v", err)
	}

	if !has {
		return nil, fmt.Errorf("this channel is not found or forbidden")
	}

	return channel, nil
}

// GetChannelList userIDのユーザーから見えるチャンネルの一覧を取得する
func GetChannelList(userID string) ([]*Channel, error) {
	// TODO: 隠しチャンネルを表示するかどうかをクライアントと決める
	var channelList []*Channel
	err := db.Join("LEFT", "users_private_channels", "users_private_channels.channel_id = channels.id").Where("(is_public = true OR user_id = ?) AND is_deleted = false", userID).Find(&channelList)

	if err != nil {
		return nil, fmt.Errorf("failed to find channels: %v", err)
	}
	return channelList, nil
}

// GetAllChannels 全てのチャンネルを取得する
func GetAllChannels() (channels []*Channel, err error) {
	err = db.Find(&channels)
	return
}

// GetChannelPath 指定したIDのチャンネルのパス文字列を取得する
func GetChannelPath(id uuid.UUID) (string, bool) {
	v, ok := channelPathMap.Load(id)
	if !ok {
		return "", false
	}

	return v.(string), true
}

func updateChannelPathWithDescendants(channel *Channel) error {
	path, err := channel.Path()
	if err != nil {
		return err
	}

	channelPathMap.Store(uuid.FromStringOrNil(channel.ID), path)

	//子チャンネルも
	var children []*Channel
	if err = db.Where("parent_id = ?", channel.ID).Find(&children); err != nil {
		return err
	}

	for _, v := range children {
		if err := updateChannelPathWithDescendants(v); err != nil {
			return err
		}
	}

	return nil
}

func validateChannelName(name string) error {
	if !regexp.MustCompile(`^[a-zA-Z0-9-_]*$`).Match([]byte(name)) {
		return fmt.Errorf("alphabet, hyphen and underscore are only allowed to use")
	}

	if len(name) > 20 {
		return fmt.Errorf("channel name should be up to 20 characters")
	}
	return nil
}
