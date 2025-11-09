package client

import (
	"context"
	"errors"

	"gitlab.xiaoduoai.com/cvd/timer/proto"
	pb_timer "gitlab.xiaoduoai.com/cvd/timer/proto/timer"
	"google.golang.org/grpc"
)

/*
cc: grpc client
setType: 类型：new / upsert
id：选填，若不填由服务器生成 【************重要：若填了已有的id会被修改************】
topic：时钟到期回调通知的topic 【************重要：请使用topic生成函数获取************】
expireSeconds：超时时间，不得小于 proto.ErrInvalidDelayTime
msgNum：消息编号，用来配置cvd这边的nsq消息路由，不需要回调路由可填0
data：自定义数据
notify: 更新时是否发送通知 （true：发送，false：不发送）
*/
func setTimer(ctx context.Context, cc *grpc.ClientConn, setType pb_timer.TimerType, id, eventType, topic, subTopic string, msgNum, expireSeconds int64, data string, notify, notOverWriteData bool) (*pb_timer.TimerInfo, error) {
	switch setType {
	case pb_timer.TimerType_ETimerTypeNew, pb_timer.TimerType_ETimerTypeUpsert:
	default:
		return nil, proto.ErrInvalidTimerType
	}
	if len(topic) <= 0 {
		return nil, proto.ErrInvalidTopic
	}
	if expireSeconds < proto.EMinDelaySeconds {
		return nil, proto.ErrInvalidDelayTime
	}
	mdstr := encodeCtx(ctx)

	t := pb_timer.TimerInfo{
		Id:             proto.GetTimerID(id),
		Type:           setType,
		DelaySeconds:   expireSeconds,
		Topic:          topic,
		SubTopic:       subTopic,
		TopicMsg:       int32(msgNum),
		Data:           data,
		NotNotify:      !notify,
		NotReplaceData: notOverWriteData,
		EventType:      eventType,
		Metadata:       mdstr,
	}
	req := &pb_timer.TimerReq{TimerInfo: &t}
	c := pb_timer.NewTimerClient(cc)
	ack, err := c.SetTimer(ctx, req)
	if err != nil {
		return nil, err
	}
	if ack.ErrInfo.ErrCode != proto.ECodeNone {
		return nil, errors.New(ack.ErrInfo.ErrMsg)
	}
	return ack.TimerInfo, nil
}

func getTimers(ctx context.Context, cc *grpc.ClientConn, params []*pb_timer.GetTimerReq) ([]*pb_timer.TimerInfo, error) {
	req := &pb_timer.TimerListReq{
		TimerList: params,
	}
	c := pb_timer.NewTimerClient(cc)
	ack, err := c.GetTimers(ctx, req)
	if err != nil {
		return nil, err
	}
	if ack.ErrInfo.ErrCode != proto.ECodeNone {
		return nil, errors.New(ack.ErrInfo.ErrMsg)
	}
	return ack.TimerList, nil
}

func deleteTimer(ctx context.Context, cc *grpc.ClientConn, topic, subTopic, id string, notify bool) error {
	t := pb_timer.TimerInfo{
		Id:        id,
		Topic:     topic,
		SubTopic:  subTopic,
		Type:      pb_timer.TimerType_ETimerEndTypeDelete,
		NotNotify: !notify,
	}
	req := &pb_timer.TimerReq{TimerInfo: &t}
	c := pb_timer.NewTimerClient(cc)
	ack, err := c.DeleteTimer(ctx, req)
	if err != nil {
		return err
	}
	if ack.ErrInfo.ErrCode != proto.ECodeNone {
		return errors.New(ack.ErrInfo.ErrMsg)
	}
	return nil
}
