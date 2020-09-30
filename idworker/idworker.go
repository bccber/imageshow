package idworker

/*
	模仿雪花算法，简化的一个id生成器
	64bit
	-----------------------------------------------------
	|保留1位| 大概41位时间戳	| 12位的序列号	|2位的分库基因|
	-----------------------------------------------------
*/
import (
	"sync"
	"time"
)

const (
	// 分库基因最大位数，因为分库个数为4
	// 所以用2bit就够了，可以根据业务适当冗余
	// 但分库个数一定要是2的N次方，否则基因法可能失效
	dnaBits = 2

	// 基因最大值
	dnaMask int64 = -1 ^ (-1 << dnaBits)

	// 序列号12位
	sequenceBits int64 = 12

	// 序列号最大值
	sequenceMask int64 = -1 ^ (-1 << sequenceBits)

	//时间戳需要左移位数
	timestampLeftShift int64 = sequenceBits

	// 初始时间戳,一个神奇的数字
	twepoch int64 = 1136185445000
)

var (
	// 锁
	m sync.Mutex
	// 上次时间戳，初始值为负数
	lastTimestamp int64 = -1
	// 16位的序列号
	sequence int64 = 0
)

// GetIdByDNA 基因法获取唯一ID
func GetIdByDNA(dnaNum int64) int64 {
	if dnaNum < 0 {
		return 0
	}

	dnaNum = dnaNum & dnaMask

	id := GetId()
	// 左移2位
	id = (id << dnaBits) | dnaNum

	return id
}

// GetId 获取唯一ID
func GetId() int64 {
	m.Lock()
	defer m.Unlock()

	timestamp := time.Now().UnixNano() / 1e6
	if lastTimestamp == timestamp {
		// 时间戳不同，序列号累加
		sequence = (sequence + 1) & sequenceMask
		if sequence == 0 {
			// 大于sequenceMask，从重新取时间戳
			for timestamp <= lastTimestamp {
				timestamp = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		// 时间戳不同，序列号从0开始
		sequence = 0
	}

	lastTimestamp = timestamp
	timestamp = timestamp - twepoch

	id := (timestamp << timestampLeftShift) | sequence

	return id
}
