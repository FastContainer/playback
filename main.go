package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/carlescere/scheduler"
)

type BulkMail struct {
	from         string
	to           string
	subject      string
	sessionCount int
	messageCount int
	interval     int
}

var cmder Cmder = Cmd{}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) > 0 && args[0] == "dryrun" {
		cmder = MockCmd{Out: ""}
	}

	totalTime := 10 * time.Second
	case1 := &BulkMail{from: "root@example-a.local", to: "root@smtp-rcpt", subject: "aaa", sessionCount: 10, messageCount: 100, interval: 2}
	case2 := &BulkMail{from: "root@example-b.local", to: "root@smtp-rcpt", subject: "bbb", sessionCount: 1, messageCount: 1, interval: 3}
	case3 := &BulkMail{from: "root@example-c.local", to: "ca@smtp-tarpit", subject: "ccc", sessionCount: 1, messageCount: 1, interval: 5}
	case4 := &BulkMail{from: "root@example-d.local", to: "root@smtp-rcpt", subject: "ddd", sessionCount: 100, messageCount: 1000, interval: 5}

	// Playback 1: Containers
	job1, _ := scheduler.Every(case1.interval).Seconds().Run(func() { case1.send("containers:58025") })
	job2, _ := scheduler.Every(case2.interval).Seconds().Run(func() { case2.send("containers:58026") })
	job3, _ := scheduler.Every(case3.interval).Seconds().Run(func() { case3.send("containers:58027") })
	job4, _ := scheduler.Every(case4.interval).Seconds().Run(func() { case4.send("containers:58028") })

	// Playback 2: Monolithic
	monolithic := "monolithic:58025"
	job5, _ := scheduler.Every(case1.interval).Seconds().Run(func() { case1.send(monolithic) })
	job6, _ := scheduler.Every(case2.interval).Seconds().Run(func() { case2.send(monolithic) })
	job7, _ := scheduler.Every(case3.interval).Seconds().Run(func() { case3.send(monolithic) })
	job8, _ := scheduler.Every(case4.interval).Seconds().Run(func() { case4.send(monolithic) })

	time.Sleep(totalTime)
	job1.Quit <- true
	job2.Quit <- true
	job3.Quit <- true
	job4.Quit <- true
	job5.Quit <- true
	job6.Quit <- true
	job7.Quit <- true
	job8.Quit <- true
	fmt.Printf("job finish!\n")
}

func (m *BulkMail) send(by string) ([]byte, error) {
	args := []string{
		"-c",
		"-S", m.subject,
		"-f", m.from,
		"-t", m.to,
		"-s", strconv.Itoa(m.sessionCount),
		"-m", strconv.Itoa(m.messageCount),
		by,
	}
	return cmder.Do("smtp-source", args...)
}
