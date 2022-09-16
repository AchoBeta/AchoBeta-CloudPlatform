package util

import (
	"errors"
	"sync"
	"time"

	"github.com/golang/glog"
)

type Snowflake struct {
	mu        sync.Mutex
	timestamp int64
	workerId  int64
	seq       int64
}

const (
	machineId   int64 = 1 << 1
	workerBits  uint8 = 10                      // 节点数
	seqBits     uint8 = 12                      // 1毫秒内可生成的id序号的二进制位数
	workerMax   int64 = -1 ^ (-1 << workerBits) // 节点ID的最大值，用于防止溢出
	seqMax      int64 = -1 ^ (-1 << seqBits)    // 同上，用来表示生成id序号的最大值
	timeShift   uint8 = workerBits + seqBits    // 时间戳向左的偏移量
	workerShift uint8 = seqBits                 // 节点ID向左的偏移量
	epoch       int64 = 1567906170596           // 开始运行时间
)

var snow *Snowflake = nil
var mutex sync.Mutex

func GetNextSnowflakeID() int64 {
	// 这里用个单例吧
	if snow == nil {
		mutex.Lock()
		defer mutex.Unlock()
		if snow == nil {
			var err error
			snow, err = newSnowflake(machineId)
			if err != nil {
				glog.Errorf("create snowflake error ! msg: ", err.Error())
				return -1
			}
		}
	}
	return snow.Next()
}

// 实例化对象
func newSnowflake(workerId int64) (*Snowflake, error) {
	// 要先检测workerId是否在上面定义的范围内
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("worker ID excess of quantity")
	}
	// 生成一个新节点
	return &Snowflake{
		timestamp: 0,
		workerId:  workerId,
		seq:       0,
	}, nil
}

func (w *Snowflake) Next() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	// 获取生成时的时间戳
	now := time.Now().UnixNano() / 1e6 // 纳秒转毫秒
	if w.timestamp == now {
		w.seq = (w.seq + 1) & seqMax
		// 这里要判断，当前工作节点是否在1毫秒内已经生成seqMax个ID
		if w.seq == 0 {
			// 如果当前工作节点在1毫秒内生成的ID已经超过上限 需要等待1毫秒再继续生成
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		// 如果当前时间与工作节点上一次生成ID的时间不一致 则需要重置工作节点生成ID的序号
		w.seq = 0
	}
	w.timestamp = now // 将机器上一次生成ID的时间更新为当前时间
	// 第一段 now - epoch 为该算法目前已经奔跑了xxx毫秒
	// 如果在程序跑了一段时间修改了epoch这个值 可能会导致生成相同的ID
	ID := int64((now-epoch)<<timeShift | (w.workerId << workerShift) | (w.seq))
	return ID
}
