package client

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"strings"
	"sync"
	"time"

	"gitlab.xiaoduoai.com/golib/xd_sdk/gray/open/transfer"
	proto "gitlab.xiaoduoai.com/golib/xd_sdk/gray/proto/pb"

	"gitlab.xiaoduoai.com/golib/xd_sdk/gray/open"
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"

	//"github.com/vmihailenco/msgpack"

	pb_timer "gitlab.xiaoduoai.com/cvd/timer/proto/timer"
	"google.golang.org/grpc"
)

type TimerLogger interface {
	Info(ctx context.Context, pairs ...interface{})
	Debug(ctx context.Context, pairs ...interface{})
	Warn(ctx context.Context, pairs ...interface{})
	Error(ctx context.Context, pairs ...interface{})
}

type defaultLogger struct{}

func (l *defaultLogger) Info(ctx context.Context, pairs ...interface{}) {
	logger.Info(ctx, pairs)
}

func (l *defaultLogger) Debug(ctx context.Context, pairs ...interface{}) {
	logger.Debug(ctx, pairs)
}

func (l *defaultLogger) Warn(ctx context.Context, pairs ...interface{}) {
	logger.Warn(ctx, pairs)
}

func (l *defaultLogger) Error(ctx context.Context, pairs ...interface{}) {
	logger.Error(ctx, pairs)
}

type TimerEvent struct {
	SubTopic         string
	ID               string
	Type             string
	ExpiredIn        int64
	NotOverWriteData bool
	Data             string
}

type TimerEventHdl func(context.Context, *TimerEvent)

type TimerClient struct {
	Topic string
	Addr  string
	Conn  *grpc.ClientConn
	Hdl   TimerEventHdl
	Err   error
	L     TimerLogger
	lock  sync.RWMutex

	// for waiting graceful shutdown
	WaitingForGracefulClose func()
}

func NewTimer(ctx context.Context, addr, topic string, hdl TimerEventHdl, logger TimerLogger) *TimerClient {
	if logger == nil {
		logger = &defaultLogger{}
	}

	rand.Seed(time.Now().Unix())

	tc := &TimerClient{Topic: topic, Addr: addr, Hdl: hdl, L: logger}

	// set default
	tc.WaitingForGracefulClose = func() {
		tc.L.Info(ctx, "default WaitingForGracefulClose func")
	}

	tc.Connect(ctx)

	return tc
}

func (t *TimerClient) Connect(ctx context.Context) *TimerClient {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// check scheme, format is scheme://authority/endpoint
	// no scheme, use grpc scheme
	if !strings.Contains(t.Addr, "://") {
		t.Addr = "grpc://auth/" + t.Addr
	}

	var err error
	t.Conn, err = grpc.DialContext(ctx, t.Addr, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	if err != nil {
		t.L.Error(ctx, "grpc.Dial failed", "err", err)
	}

	return t
}

func (t *TimerClient) UpsertTimeContext(ctx context.Context, event *TimerEvent) (string, error) {
	if t.Conn == nil {
		return "", errors.New("timer client is not connected")
	}
	info, err := setTimer(ctx, t.Conn, pb_timer.TimerType_ETimerTypeUpsert, event.ID, event.Type, t.Topic, event.SubTopic, 0, event.ExpiredIn, event.Data, false, event.NotOverWriteData)
	id := ""
	if info != nil {
		id = info.Id
	}
	return id, err
}

func (t *TimerClient) SetTimerContext(ctx context.Context, event *TimerEvent) (error, string) {
	if t.Conn == nil {
		return errors.New("timer client is not connected"), ""
	}

	info, err := setTimer(ctx, t.Conn, pb_timer.TimerType_ETimerTypeNew, event.ID, event.Type, t.Topic, event.SubTopic, 0, event.ExpiredIn, event.Data, false, event.NotOverWriteData)
	id := ""
	if info != nil {
		id = info.Id
	}
	//
	//if err != nil {
	//}
	return err, id
}

func (t *TimerClient) SetTimer(ctx context.Context, event *TimerEvent) (error, string) {
	return t.SetTimerContext(ctx, event)
}

func (t *TimerClient) DeleteTimer(ctx context.Context, subTopic, id string) error {
	return deleteTimer(ctx, t.Conn, t.Topic, subTopic, id, false)
}

func (t *TimerClient) GetTimers(ctx context.Context, params []*pb_timer.GetTimerReq) ([]*pb_timer.TimerInfo, error) {
	for _, v := range params {
		v.Topic = t.Topic
	}
	return getTimers(ctx, t.Conn, params)
}

func (t *TimerClient) genGrayTypeName() string {
	return grayTypeName(t.Topic)
}

func (t *TimerClient) Watch(ctx context.Context) *TimerClient {

	recvFinishCH := make(chan bool, 0)
	t.WaitingForGracefulClose = func() {
		<-recvFinishCH
		t.L.Info(ctx, "timer Watch conn shutdown gracefully !")
	}

	defer func() {
		close(recvFinishCH)
	}()

	open.GetAPI().DefaultInit(ctx)

	if open.GetAPI().IsGray() {
		transfer.GetAPI().RegisterTransferHandle(t.genGrayTypeName(), &TransferHandle{ch: make(chan *proto.TransferMsg, 1)})
	}

	if open.GetAPI().IsGray() {
		return t.watchGray(ctx)
	} else {
		return t.watch(ctx)
	}
}

func (t *TimerClient) watchGray(ctx context.Context) *TimerClient {

	var timerInfo *pb_timer.TimerInfo
	var ok bool

	for {
		select {
		case <-ctx.Done():
			return t
		default:
		}

		grayTimerClient := transfer.GetAPI().GetTransferHandle(t.genGrayTypeName())
		if grayTimerClient == nil {
			time.Sleep(time.Second)
			logger.Error(ctx, "gray timer miss transfer msg handle, typ="+t.genGrayTypeName())
			return nil
		}

		ti := grayTimerClient.Receive(ctx)
		if ti == nil {
			logger.Error(ctx, "gray timer received nil msg, typ="+t.genGrayTypeName())
			return nil
		}

		if timerInfo, ok = ti.(*pb_timer.TimerInfo); !ok {
			logger.Error(ctx, "gray pulsar received msg miss-match type, typ="+t.genGrayTypeName())
			return nil
		}
		event := &TimerEvent{SubTopic: timerInfo.SubTopic, ID: timerInfo.Id, Data: timerInfo.Data, Type: timerInfo.EventType, ExpiredIn: timerInfo.DelaySeconds}
		ctx, end := decodeToCtx(timerInfo.Metadata)
		t.Hdl(ctx, event)
		end()
	}
}

func (t *TimerClient) watch(ctx context.Context) *TimerClient {
	for {
		select {
		case <-ctx.Done():
			return t
		default:
		}
		rpcClient := pb_timer.NewTimerClient(t.Conn)

		stream, err := rpcClient.NotifyB(ctx)
		if err != nil {
			t.Err = err
			t.L.Error(ctx, "timer client create watch failed", "err", err)
			time.Sleep(time.Second)
			continue
		}

		t.L.Info(ctx, "timer client create watch ok")
		req := &pb_timer.WatchReq{Topic: t.Topic}
		if err := stream.Send(req); err != nil {
			t.L.Warn(ctx, "timer client send WatchReq failed, err: ", err)
		}

		stop := make(chan bool, 0)
		go func(s pb_timer.Timer_NotifyBClient) {
			select {
			case <-ctx.Done():
			case <-stop:
				return
			}
			t.L.Info(ctx, "timer client graceful shutdown")
			if err := s.Send(&pb_timer.WatchReq{Stop: true}); err != nil {
				t.L.Warn(ctx, "timer client graceful shutdown, send STOP failed, err: ", err)
			} else {
				t.L.Info(ctx, "timer client graceful shutdown, send STOP OK")
			}
		}(stream)

		for {
			timerInfo, err := stream.Recv()
			if err != nil {
				t.Err = err
				if err != io.EOF {
					t.L.Error(ctx, "timer client watch failed", "err", err, "e", timerInfo)
				}
				break
			}
			if timerInfo.Id == "" {
				continue
			}
			event := &TimerEvent{SubTopic: timerInfo.SubTopic, ID: timerInfo.Id, Data: timerInfo.Data, Type: timerInfo.EventType, ExpiredIn: timerInfo.DelaySeconds}
			ctx, end := decodeToCtx(timerInfo.Metadata)
			if transferGray(ctx, timerInfo) {
				continue
			}
			t.Hdl(ctx, event)
			end()
		}
		close(stop)
		time.Sleep(time.Millisecond)
	}
}
