package xlog

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"sync"
)

const (
	LogLevelInfo  = "info"
	LogLevelError = "error"
)

var Default = &channel{channel: "default"}

type channel struct {
	channel string
}

func NewChannel(name string) *channel {
	return &channel{channel: name}
}

var logMaps = new(sync.Map)

func (receiver *channel) Info(msg string, data ...interface{}) {
	Logger(LogLevelInfo, receiver.channel).Info(context.Background(), msg, data)
}

func (receiver *channel) Error(msg string, data ...interface{}) {
	Logger(LogLevelError, receiver.channel).Error(context.Background(), msg, data)
}

func Logger(level string, channel ...string) *glog.Logger {
	var _channel = "default"
	if len(channel) > 0 {
		_channel = channel[0]
	}

	key := fmt.Sprintf("%s:%s", level, _channel)
	if _val, ok := logMaps.Load(key); ok {
		return _val.(*glog.Logger)
	}

	_log := glog.New()
	err := _log.SetConfigWithMap(g.Map{
		"path":   fmt.Sprintf("./runtime/logs/%s/%s", _channel, level),
		"level":  level,
		"stdout": false,
	})

	if err != nil {
		fmt.Println("log init error", level, channel)
		//return nil
	}
	logMaps.Store(key, _log)
	return _log
}
