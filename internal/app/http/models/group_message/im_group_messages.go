/**
  @author:kk
  @data:2021/12/8
  @note
**/
package group_message

import (
	"github.com/golang-module/carbon"
	"gorm.io/gorm"
	userModel "im_app/internal/app/http/models/user"
)

type ImGroupMessages struct {
	ID        int64           `json:"id"`
	MsgType   int             `json:"msg_type"`
	Msg       string          `json:"msg"`
	GroupId   int64           `json:"group_id"`
	FromId    int64           `json:"user_id"`
	Status    int             `json:"status" gorm:"-" valid:"status"` //忽略一下该字段写入
	CreatedAt string          `json:"created_at"`
	Users     userModel.Users `json:"users,omitempty" gorm:"foreignKey:FromId;references:ID"`
}

func (ImGroupMessages) TableName() string {
	return "im_group_messages"
}

func (a *ImGroupMessages) AfterFind(tx *gorm.DB) (err error) {
	if a.CreatedAt != "" {
		a.CreatedAt = carbon.Parse(a.CreatedAt).SetLocale("zh-CN").DiffForHumans()
	}
	a.Status = 0
	return
}
