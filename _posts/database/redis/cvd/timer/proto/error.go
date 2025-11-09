package proto

import (
	"errors"
	"fmt"
)

const (
	ECodeNone        = 0
	ECodeErrorParam  = 1
	ECodeErrorSystem = 2
)

var (
	ErrInvalidTimerID   = errors.New("invalid timer id")
	ErrInvalidTimerType = errors.New("invalid timer type")
	ErrInvalidTopic     = errors.New("invalid topic")
	ErrInvalidDelayTime = errors.New(fmt.Sprintf("invalid delay time, min delay seconds=%d", EMinDelaySeconds))
	ErrSystemFailed     = errors.New("system failed")
)
