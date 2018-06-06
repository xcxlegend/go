package antnet

import (
	"time"
)

func Sleep(ms int) {
	time.Sleep(time.Millisecond * time.Duration(ms))
}

func SetTimeout(inteval int, fn func(...interface{}) int, args ...interface{}) {
	if inteval < 0 {
		LogError("new timerout inteval:%v", inteval)
		return
	}
	LogInfo("new timerout inteval:%v", inteval)
	Go2(func(cstop chan struct{}) {
		for inteval > 0 {
			tick := time.NewTicker(time.Millisecond * time.Duration(inteval))
			select {
			case <-cstop:
				inteval = 0
			case <-tick.C:
				inteval = fn(args...)
			}
		}
	})
}

func timerTick() {
	TimeNanoStamp = time.Now().UnixNano()
	StartTick = time.Now().UnixNano() / 1000000
	NowTick = StartTick
	Timestamp = NowTick / 1000
	Go(func() {
		for IsRuning() {
			Sleep(1)
			NowTick = time.Now().UnixNano() / 1000000
			Timestamp = NowTick / 1000
			TimeNanoStamp = time.Now().UnixNano()
		}
	})
}

/**
* @brief 获得timestamp距离下个小时的时间，单位s
*
* @return uint32_t 距离下个小时的时间，单位s
 */
func GetNextHourIntervalS() int {
	return int(3600 - (Timestamp % 3600))
}

/**
 * @brief 获得timestamp距离下个小时的时间，单位ms
 *
 * @return uint32_t 距离下个小时的时间，单位ms
 */
func GetNextHourIntervalMS() int {
	return GetNextHourIntervalS() * 1000
}

/**
* @brief 时间戳转换为小时，24小时制，0点用24表示
*
* @param timestamp 时间戳
* @param timezone  时区
* @return uint32_t 小时 范围 1-24
 */
func GetHour24(timestamp int64, timezone int) int {
	hour := (int((timestamp%86400)/3600) + timezone)
	if hour > 24 {
		return hour - 24
	}
	return hour
}

/**
 * @brief 时间戳转换为小时，24小时制，0点用0表示
 *
 * @param timestamp 时间戳
 * @param timezone  时区
 * @return uint32_t 小时 范围 0-23
 */
func GetHour23(timestamp int64, timezone int) int {
	hour := GetHour24(timestamp, timezone)
	if hour == 24 {
		return 0 //24点就是0点
	}
	return hour
}

func GetHour(timestamp int64, timezone int) int {
	return GetHour23(timestamp, timezone)
}

/**
* @brief 判断两个时间戳是否是同一天
*
* @param now 需要比较的时间戳
* @param old 需要比较的时间戳
* @param timezone 时区
* @return uint32_t 返回不同的天数
 */
func IsDiffDay(now, old int64, timezone int) int {
	now += int64(timezone * 3600)
	old += int64(timezone * 3600)
	return int((now / 86400) - (old / 86400))
}

/**
* @brief 判断时间戳是否处于一个小时的两边，即一个时间错大于当前的hour，一个小于
*
* @param now 需要比较的时间戳
* @param old 需要比较的时间戳
* @param hour 小时，0-23
* @param timezone 时区
* @return bool true表示时间戳是否处于一个小时的两边
 */
func IsDiffHour(now, old int64, hour, timezone int) bool {
	diff := IsDiffDay(now, old, timezone)
	if diff == 1 {
		if GetHour23(old, timezone) > hour {
			return GetHour23(now, timezone) >= hour
		}
	} else if diff >= 2 {
		return true
	}

	return (GetHour23(now, timezone) >= hour) && (GetHour23(old, timezone) < hour)
}
