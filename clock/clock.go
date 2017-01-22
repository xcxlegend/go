package clock

import (
	"time"
)

type ActionChan int

const (
	ACTIONCHAN_NONE ActionChan = iota
	ACTIONCHAN_START
	ACTIONCHAN_PAUSE
	ACTIONCHAN_CLOSE
)

func NewClock() *clock {
	var c = new(clock)
	c.chans = new(ChanGroup)
	c.chans.Init()
	go c.Run()
	return c
}

type clock struct {
	chans *ChanGroup
	time  int
	isRun bool
}

func (c *clock) Start() {
	c.chans.action <- ACTIONCHAN_START
}

func (c *clock) Pause() {
	c.chans.action <- ACTIONCHAN_PAUSE
}

func (c *clock) Reset(t int) {
	c.chans.reset <- t
}

func (c *clock) Close() {
	c.chans.action <- ACTIONCHAN_CLOSE
}

func (c *clock) ShowTime() int {
	return c.time
}

func (c *clock) WaitAlarm() bool {
	return <-c.chans.alarm
}

type ChanGroup struct {
	action chan ActionChan
	reset  chan int
	time   chan int
	alarm  chan bool
}

func (c *ChanGroup) Init() {
	c.action = make(chan ActionChan, 4)
	c.reset = make(chan int)
	c.time = make(chan int)
	c.alarm = make(chan bool)
}

func (this *clock) Run() {

	var timer = time.NewTicker(time.Second)
	for {
		select {
		case <-timer.C:
			if this.isRun {
				this.time--
				if this.time <= 0 {
					this.time = 0
					this.isRun = false
					this.chans.alarm <- true
					break
				}
			}
		case c := <-this.chans.action:
			if c == ACTIONCHAN_NONE {
				return
			}
			switch c {
			case ACTIONCHAN_START:
				this.isRun = true
				break
			case ACTIONCHAN_PAUSE:
				this.isRun = false
				break
			case ACTIONCHAN_CLOSE:
				return
				break
			}

			break
		case t := <-this.chans.reset:
			if t <= 0 {
				continue
			}
			this.time = t
			this.isRun = true
			break
		}

	}
}
