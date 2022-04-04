package community

import (
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// Channel Contain ID, ChannelName
// @Tags community-forum
type Channel struct {
	general.TenantBase
	// Logo string `json:"logo" gorm:"type:varchar(100)"`
	ChannelName string `json:"channelName" example:"Golang Developer" gorm:"type:varchar(100)"`
}

// Validate Check Channel is Valid or Not
func (channel *Channel) Validate() error {
	if util.IsEmpty(channel.ChannelName) {
		return errors.NewValidationError("Channel Name Must be specified")
	}
	return nil
}

// TableName sets table name as community_channels.
func (*Channel) TableName() string {
	return "community_channels"
}
