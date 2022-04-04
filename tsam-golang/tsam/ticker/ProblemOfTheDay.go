package ticker

import (
	"time"
)

const hourToTickMidnight int = 00
const minuteToTickMidnight int = 00
const secondToTickMidnight int = 00

// MidnightTask calls the function that has to run at midnight.
func MidnightTask() {
}

// MidnightTicker calculates time to tick at midnight.
func MidnightTicker() {

	// Get current time.
	currentTime := time.Now()

	// Get midnight time (00:00:00).
	midnightTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), hourToTickMidnight, minuteToTickMidnight, secondToTickMidnight, 0, currentTime.Location())
	
	// Get difference between current and midnight time.
	differenceTime := midnightTime.Sub(currentTime)

	// If current time is ahead of midnight time then recalculate difference time.
	if differenceTime < 0 {
		midnightTime = midnightTime.Add(5 * time.Second)
		differenceTime = midnightTime.Sub(currentTime)
	}

	// Calling the midnight funtion and sleeping for 24 hours.
	for {
		time.Sleep(differenceTime)
		differenceTime = 5 * time.Second
		MidnightTask()
	}
}

// const INTERVAL_PERIOD time.Duration = 5 * time.Second

// const HOUR_TO_TICK int = 18
// const MINUTE_TO_TICK int = 02
// const SECOND_TO_TICK int = 02

// type jobTicker struct {
//     t *time.Timer
// }

// func getNextTickDuration() time.Duration {
//     now := time.Now()
//     nextTick := time.Date(now.Year(), now.Month(), now.Day(), HOUR_TO_TICK, MINUTE_TO_TICK, SECOND_TO_TICK, 0, time.Local)
//     if nextTick.Before(now) {
//         nextTick = nextTick.Add(INTERVAL_PERIOD)
//     }
//     return nextTick.Sub(time.Now())
// }

// func NewJobTicker() jobTicker {
//     fmt.Println("new tick here")
//     return jobTicker{time.NewTimer(getNextTickDuration())}
// }

// func (jt jobTicker) updateJobTicker() {
//     fmt.Println("next tick here")
//     jt.t.Reset(getNextTickDuration())
// }

// func NewTempTicker() {
//     jt := NewJobTicker()
//     for {
//         <-jt.t.C
//         fmt.Println(time.Now(), "- just ticked")
//         jt.updateJobTicker()
//     }
// }
