package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/carlescere/scheduler"
)

type BulkMail struct {
	number       int
	to           string
	subject      string
	sessionCount int
	messageCount int
	interval     int
}

var cmder Cmder = Cmd{}

const subject = "fast container"
const diminutive = "dimi-%d.test:%d"
const monolithic = "mono-%d.test:%d"

func main() {
	flag.Parse()
	args := flag.Args()
	totalTime := 10 * time.Second

	if len(args) > 0 {
		if len(args) > 1 && args[1] == "dryrun" {
			cmder = MockCmd{Out: ""}
		}

		if i, err := strconv.Atoi(args[0]); err == nil {
			totalTime = time.Duration(i) * time.Minute
		}
	}

	case1 := &BulkMail{number: 1, to: "root@smtp-rcpt", subject: subject, sessionCount: 1, messageCount: 10, interval: 10}
	case2 := &BulkMail{number: 2, to: "root@smtp-tarpit", subject: subject, sessionCount: 1, messageCount: 10, interval: 30}
	//case2 := &BulkMail{number: 2, to: "root@smtp-rcpt", subject: subject, sessionCount: 1, messageCount: 1, interval: 3}
	//case4 := &BulkMail{number: 4, to: "root@smtp-rcpt", subject: subject, sessionCount: 1, messageCount: 100, interval: 5}

	// Playback 1: Containers
	dimi1, _ := scheduler.Every(case1.interval).Seconds().NotImmediately().Run(func() { case1.send(fmt.Sprintf(diminutive, 1, 58025)) })
	dimi2, _ := scheduler.Every(case2.interval).Seconds().NotImmediately().Run(func() { case2.send(fmt.Sprintf(diminutive, 2, 58026)) })
	//dimi3, _ := scheduler.Every(case3.interval).Seconds().Run(func() { case3.send(fmt.Sprintf(diminutive, 3, 58027)) })
	//dimi4, _ := scheduler.Every(case4.interval).Seconds().Run(func() { case4.send(fmt.Sprintf(diminutive, 4, 58028)) })

	// Playback 2: Monolithic
	mono1, _ := scheduler.Every(case1.interval).Seconds().NotImmediately().Run(func() { case1.send(fmt.Sprintf(monolithic, 1, 25)) })
	mono2, _ := scheduler.Every(case2.interval).Seconds().NotImmediately().Run(func() { case2.send(fmt.Sprintf(monolithic, 2, 25)) })
	//mono3, _ := scheduler.Every(case3.interval).Seconds().Run(func() { case3.send(fmt.Sprintf(monolithic, 3, 25)) })
	//mono4, _ := scheduler.Every(case4.interval).Seconds().Run(func() { case4.send(fmt.Sprintf(monolithic, 4, 25)) })

	time.Sleep(totalTime)

	dimi1.Quit <- true
	dimi2.Quit <- true
	//dimi3.Quit <- true
	//dimi4.Quit <- true

	mono1.Quit <- true
	mono2.Quit <- true
	//mono3.Quit <- true
	//mono4.Quit <- true

	fmt.Printf("job finish!\n")
}

func (m *BulkMail) send(by string) ([]byte, error) {
	args := []string{
		"-c",
		"-S", m.subject,
		"-f", by,
		"-t", m.to,
		"-s", strconv.Itoa(m.sessionCount),
		"-m", strconv.Itoa(m.messageCount),
		by,
	}
	fmt.Printf("%s %s\n", time.Now().Format("15:04:05"), by)
	return cmder.Do("smtp-source", args...)
}
