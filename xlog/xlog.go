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
var User = &channel{channel: "user"}
var Prop = &channel{channel: "prop"}
var UserLevel = &channel{channel: "user_level"}
var UserProp = &channel{channel: "user_prop"}
var UserEarning = &channel{channel: "user_earning"}
var Team = &channel{channel: "team"}
var Game = &channel{channel: "game"}
var Mine = &channel{channel: "mine"}
var Task = &channel{channel: "task"}
var Wallet = &channel{channel: "wallet"}
var Family = &channel{channel: "family"}
var Bot = &channel{channel: "bot"}
var Ton = &channel{channel: "ton"}
var SystemConfig = &channel{channel: "system_config"}
var Rank = &channel{channel: "rank"}
var RedisToSql = &channel{channel: "redisToSql"}
var ReportData = &channel{channel: "report_data"}
var Ad = &channel{channel: "ad"}
var Heist = &channel{channel: "heist"}

type channel struct {
	channel string
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
