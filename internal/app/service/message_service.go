/**
  @author:kk
  @data:2021/12/15
  @note
**/
package service

import "im_app/internal/app/ws"

// SendMessage 系统单独发送
func SendMessage(code int, fId int, tId int, message string) {
	ws.ImManager.SystemMessageDelivery(int64(fId),
		&ws.Msg{Code: code, FromId: fId, Msg: message,
			ToId: tId, Status: 0, MsgType: 1,
			ChannelType: 3})
}
