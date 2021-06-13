package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/carlescere/scheduler"
)

type BulkMail struct {
	number       int
	to           string
	sessionCount int
	messageCount int
	interval     int
}

var cmder Cmder = Cmd{}

const container = "container-%d.test:%d"
const monolith = "monolith-%d.test:%d"

func (m *BulkMail) send(by string) ([]byte, error) {
	args := []string{
		"-c",
		"-S", "Hello, Gophers",
		"-f", by,
		"-t", m.to,
		"-s", strconv.Itoa(m.sessionCount),
		"-m", strconv.Itoa(m.messageCount),
		by,
	}
	fmt.Printf("%s %s\n", time.Now().Format("15:04:05"), by)
	return cmder.Do("smtp-source", args...)
}

func Bulk(dryrun bool) {
	delay := 5 * time.Minute
	totalTime := 10 * time.Minute

	if dryrun {
		fmt.Printf("dryrun -- delay: %#vmin, total: %#vmin\n", delay.Minutes(), totalTime.Minutes())
		cmder = MockCmd{Out: ""}
	}

	case1 := &BulkMail{
		number:       1,
		to:           "root@recipient",
		sessionCount: 1,
		messageCount: 10,
		interval:     10,
	}
	case2 := &BulkMail{
		number:       2,
		to:           "root@mxtarpit",
		sessionCount: 1,
		messageCount: 10,
		interval:     30,
	}

	// Playback 1: Container
	dimi1, _ := scheduler.Every(case1.interval).Seconds().NotImmediately().Run(func() {
		case1.send(fmt.Sprintf(container, 1, 58025))
	})
	// Playback 2: Monolith
	mono1, _ := scheduler.Every(case1.interval).Seconds().NotImmediately().Run(func() {
		case1.send(fmt.Sprintf(monolith, 1, 25))
	})

	time.Sleep(delay)

	// Playback 1: Container
	dimi2, _ := scheduler.Every(case2.interval).Seconds().NotImmediately().Run(func() {
		case2.send(fmt.Sprintf(container, 2, 58026))
	})
	// Playback 2: Monolith
	mono2, _ := scheduler.Every(case2.interval).Seconds().NotImmediately().Run(func() {
		case2.send(fmt.Sprintf(monolith, 2, 25))
	})

	time.Sleep(totalTime - delay)

	dimi1.Quit <- true
	dimi2.Quit <- true

	mono1.Quit <- true
	mono2.Quit <- true

	fmt.Printf("job finish!\n")
}
