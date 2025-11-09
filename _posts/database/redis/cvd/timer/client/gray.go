package client

import (
	"context"
	"encoding/json"
	pbTimer "gitlab.xiaoduoai.com/cvd/timer/proto/timer"
	"gitlab.xiaoduoai.com/golib/xd_sdk/gray/open"
	"gitlab.xiaoduoai.com/golib/xd_sdk/gray/open/transfer"
	proto "gitlab.xiaoduoai.com/golib/xd_sdk/gray/proto/pb"
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
	"gitlab.xiaoduoai.com/golib/xd_sdk/metadata"
)
type TransferHandle struct {
	ch chan *proto.TransferMsg
}

func (hdl *TransferHandle) OnRecvMsg(tm *proto.TransferMsg) error {
	hdl.ch <- tm
	return nil
}

func (hdl *TransferHandle) Receive(ctx context.Context) interface{} {
	tm := <-hdl.ch
	msg := &pbTimer.TimerInfo{}
	if err := json.Unmarshal(tm.Data, msg); err != nil {
		logger.Error(ctx, "gray timer unmarshal revc msg failed, err = "+err.Error())
		return nil
	}
	return msg
}

func grayTypeName(topic string) string {
	return "timer-" + topic
}

func transferGray(ctx context.Context, ti *pbTimer.TimerInfo) bool {
	tags := metadata.FromContext(ctx)
	if open.Run(ctx, tags, func(rule *proto.GrayRule) bool {
		transferCli := transfer.GetAPI().TransferClient(ctx, rule)
		if transferCli == nil {
			logger.Error(ctx, "gray timer transfer msg failed, err: transfer client is nil")
			return false
		}

		b, err := json.Marshal(ti)
		if err != nil {
			logger.Error(ctx, "gray timer transfer msg  marshal failed")
			return false
		}

		tm := &proto.TransferMsg{
			Type: grayTypeName(ti.Topic),
			Data: b,
		}

		ack, err := transferCli.Send(ctx, tm)
		if err != nil {
			logger.Error(ctx, "gray timer transfer msg failed, err: "+err.Error())
			return false
		}

		if ack.Code != 0 {
			logger.Error(ctx, "gray timer transfer msg failed, err: "+ack.Info)
			return false
		}
		return true
	}) {
		return true
	}
	return false
}
