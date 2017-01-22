# Clock

countdown clock

## Installation

Install Clock using the "go get" command:

```
go get github.com/xcxlegend/go
```

## Example

```
var c = clock.NewClock()
c.Pause()
<-time.After(time.Second)
c.Reset(10)
<-time.After(1 * time.Second)
c.Start()
<-time.After(time.Second)
c.Reset(5)

fmt.Println("time: ", c.ShowTime())
c.WaitAlarm()  // <- chan bool

```