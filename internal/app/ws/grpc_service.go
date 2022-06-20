// Package ws
/**
  @author:kk
  @data:2021/11/9
  @note
**/
package ws

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"im_app/api/protobuf-spec/im"
	conf "im_app/config"
	"im_app/pkg/zaplog"
	"log"
	"net"
	"net/http"
	"strconv"
)

var RpcServer = grpc.NewServer()

type ImRpcServerHandler interface {
	StartRpc()
}

type ImRpcServer struct {
}

// StartRpc 启动rpc服务
func StartRpc() {

	im.RegisterImRpcServiceServer(RpcServer, new(ImRpcServer))

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(conf.Conf.GrpcPort))
	if err != nil {
		log.Fatal("服务监听端口失败", err)
	}
	_ = RpcServer.Serve(listener)
}

// SendMessage rpc消息投递
func (ps *ImRpcServer) SendMessage(ctx context.Context, request *im.MessageRequest) (*im.MessageResponse, error) {
	jsonMessageFrom, _ := json.Marshal(&RpcMsg{Code: int(request.Code), Msg: request.Msg,
		FromId: int(request.FromId),
		ToId:   int(request.ToId), Status: 1, MsgType: int(request.MsgType), ChannelType: int(request.ChannelType)})

	zaplog.Info(jsonMessageFrom)
	if data, ok := ImManager.ImClientMap[int64(request.ToId)]; ok {
		data.Send <- jsonMessageFrom
	} else {
		MqPersonalPublish(jsonMessageFrom, int(request.ToId))
	}
	return &im.MessageResponse{Code: http.StatusOK}, nil
}
