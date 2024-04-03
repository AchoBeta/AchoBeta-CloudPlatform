package log

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

/**
 * @Description: 自定义 trace id hook
 */
type TraceIdHook struct {
	TraceId string
}

func (h *TraceIdHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

const (
	TraceIdKey       = "__trace_id"
	HeaderTraceIdKey = "__trace_id"
)

// Fire implements logrus.Hook.Fire
func (m *TraceIdHook) Fire(entry *logrus.Entry) (err error) {
	var tradeId string
	fmt.Printf("entry.Context: %v\n", entry.Context)
	if entry.Context != nil {
		var ok bool
		tradeId, ok = entry.Context.Value(TraceIdKey).(string) // 判断 entry 中的ctx字段是否存在 traceId 信息
		if !ok {
			return
		}
	}

	if tradeId == "" {
		return
	}

	entry.Data["trace_id"] = tradeId // 把 trace_id 信息加入到携带信息中，这样就可以打印出来了

	return
}

func NewTraceIdHook(traceId string) logrus.Hook {
	hook := TraceIdHook{
		TraceId: traceId,
	}
	return &hook
}
