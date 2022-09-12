package utils

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

/*
* 1                                               42           52             64
* +-----------------------------------------------+------------+---------------+
* | timestamp(ms)                                 | workerid   | sequence      |
* +-----------------------------------------------+------------+---------------+
* | 0000000000 0000000000 0000000000 0000000000 0 | 0000000000 | 0000000000 00 |
* +-----------------------------------------------+------------+---------------+
*
* 1. 41位时间截(毫秒级)，注意这是时间截的差值（当前时间截 - 开始时间截)。可以使用约70年: (1L << 41) / (1000L * 60 * 60 * 24 * 365) = 69
* 2. 10位数据机器位，可以部署在1024个节点
* 3. 12位序列，毫秒内的计数，同一机器，同一时间截并发4096个序号
 */

const (
	twepoch        = int64(1483228800000)             //开始时间截 (2017-01-01)
	workeridBits   = uint(10)                         //机器id所占的位数
	sequenceBits   = uint(12)                         //序列所占的位数
	datacenterId   = int64(2)                         //数据中心
	workeridMax    = int64(-1 ^ (-1 << workeridBits)) //支持的最大机器id数量
	sequenceMask   = int64(-1 ^ (-1 << sequenceBits)) //
	workeridShift  = sequenceBits                     //机器id左移位数
	timestampShift = sequenceBits + workeridBits      //时间戳左移位数
)

// A Snowflake struct holds the basic information needed for a snowflake generator worker
type Snowflake struct {
	sync.Mutex
	lastTimestamp int64
	workerid      int64
	sequence      int64
}

// NewSnowflake NewNode returns a new snowflake worker that can be used to generate snowflake IDs
func NewSnowflake(workerid int64) (*Snowflake, error) {

	if workerid < 0 || workerid > workeridMax {
		return nil, errors.New("workerid must be between 0 and 1023")
	}

	return &Snowflake{
		lastTimestamp: 0,
		workerid:      workerid,
		sequence:      0,
	}, nil
}

// GenerateSnowId Generate creates and returns a unique snowflake ID
func (s *Snowflake) GenerateSnowId() string {
	s.Lock()
	//获取当前毫秒数
	now := time.Now().UnixNano() / 1000000
	//如果上次生成时间和当前时间相同,在同一毫秒内
	if s.lastTimestamp == now {
		//sequence自增，因为sequence只有12bit，所以和sequenceMask相与一下，去掉高位
		s.sequence = (s.sequence + 1) & sequenceMask
		//判断是否溢出,也就是每毫秒内超过4095，当为4096时，与sequenceMask相与，sequence就等于0
		if s.sequence == 0 {
			//自旋等待到下一毫秒
			for now <= s.lastTimestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		//如果和上次生成时间不同,重置sequence，就是下一毫秒开始，sequence计数重新从0开始累加
		s.sequence = 0
	}
	s.lastTimestamp = now
	r := (now-twepoch)<<timestampShift | (s.workerid << workeridShift) | (s.sequence)
	orderId := strconv.FormatInt(r, 10)
	s.Unlock()
	return orderId
}

/*
	Generate creates and returns a custom snowflake ID
	年月日时分秒毫秒+7位随机数
	return 2105271946072968397239
	年月日时分秒毫秒+4位随机数
	return 2105272022427398193
*/

func (s *Snowflake) GenerateTimestampId() string {
	s.Lock()
	now := time.Now().UnixNano() / 1000000
	if s.lastTimestamp == now {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			for now <= s.lastTimestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		s.sequence = 0
	}
	s.lastTimestamp = now
	//计算毫秒数
	diff := now - ((now / 1e3) * 1e3)
	//毫秒时间戳
	timestamp := time.Unix(now/1e3, 0).Format("060102150405") + strconv.FormatInt(diff, 10)
	//生成7位随机数
	r := datacenterId<<timestampShift | (s.workerid << workeridShift) | (s.sequence)
	//生成4位随机数
	//r := (s.workerid << workeridShift) | (s.sequence)
	orderId := timestamp + strconv.FormatInt(r, 10)
	s.Unlock()
	return orderId
}
