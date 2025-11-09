package proto

import (
	"encoding/hex"
	"strings"

	"gitlab.xiaoduoai.com/cvd/common/id"
)

const (
	EMinDelaySeconds = 2
)

func GetTopic(filter ...string) string {
	const gap = "-"
	if len(filter) <= 0 {
		filter = []string{"NoneTopic"}
	}
	return "service_timer_topic" + gap + strings.Join(filter, gap)
}

func GetTimerID(selfID string) string {
	if selfID == "" {
		selfID = hex.EncodeToString(id.ObjectID())
	}
	return selfID
}
