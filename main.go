package main

import (
	"flag"
	"fmt"
	"net/smtp"
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
}

func resource() {
	from := fmt.Sprintf("root@mono-%d.test", 1)
	to := []string{"recipient@example.net"}
	msg := []byte("To: recipient@example.net\r\nSubject: discount Gophers!\r\n\r\nThis is the email body.\r\n")
	err := SendMail("monolith:25", from, to, msg)
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
}

func SendMail(addr string, from string, to []string, msg []byte) error {
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}

	defer c.Close()
	if err = c.hello(); err != nil {
		return err
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func inline() {
	flag.Parse()
	args := flag.Args()
	delay := 5 * time.Minute
	totalTime := 10 * time.Minute

	if len(args) > 0 {
		if args[0] == "dryrun" {
			fmt.Printf("dryrun -- delay: %#vmin, total: %#vmin\n", delay.Minutes(), totalTime.Minutes())
			cmder = MockCmd{Out: ""}
		} else {
			if i, err := strconv.Atoi(args[0]); err == nil {
				totalTime = time.Duration(i) * time.Minute
			}
			if len(args) > 1 {
				if i, err := strconv.Atoi(args[1]); err == nil {
					delay = time.Duration(i) * time.Minute
				}
			}
		}
	}

	case1 := &BulkMail{number: 1, to: "root@smtp-rcpt", subject: subject, sessionCount: 1, messageCount: 10, interval: 10}
	case2 := &BulkMail{number: 2, to: "root@smtp-tarpit", subject: subject, sessionCount: 1, messageCount: 10, interval: 30}

	// Playback 1: Containers
	dimi1, _ := scheduler.Every(case1.interval).Seconds().NotImmediately().Run(func() { case1.send(fmt.Sprintf(diminutive, 1, 58025)) })
	// Playback 2: Monolithic
	mono1, _ := scheduler.Every(case1.interval).Seconds().NotImmediately().Run(func() { case1.send(fmt.Sprintf(monolithic, 1, 25)) })

	time.Sleep(delay)

	// Playback 1: Containers
	dimi2, _ := scheduler.Every(case2.interval).Seconds().NotImmediately().Run(func() { case2.send(fmt.Sprintf(diminutive, 2, 58026)) })
	// Playback 2: Monolithic
	mono2, _ := scheduler.Every(case2.interval).Seconds().NotImmediately().Run(func() { case2.send(fmt.Sprintf(monolithic, 2, 25)) })

	time.Sleep(totalTime - delay)

	dimi1.Quit <- true
	dimi2.Quit <- true

	mono1.Quit <- true
	mono2.Quit <- true

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
