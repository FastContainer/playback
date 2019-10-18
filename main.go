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
const resourceLimit = 25

func main() {
	resource()
}

func resource() {
	totalTime := 10 * time.Minute
	count := 0
	job, _ := scheduler.Every(15).Seconds().NotImmediately().Run(func() {
		SendMails(count * resourceLimit)
		count++
	})
	time.Sleep(totalTime)
	job.Quit <- true
	fmt.Printf("job finish!\n")
}

func SendMails(startNum int) {
	stopNum := startNum + resourceLimit

	for i := startNum; i < stopNum; i++ {
		host := "monolith:25"
		from := fmt.Sprintf("root@mono-%d.test", i)
		to := "root@recipient"
		msg := []byte(fmt.Sprintf("To: %s\r\nSubject: discount Gophers!\r\n\r\nThis is the email body.\r\n", to))
		go SendMail(host, from, []string{to}, msg)
	}
}

func SendMail(addr string, from string, to []string, msg []byte) error {
	c, err := smtp.Dial(addr)
	if err != nil {
		fmt.Printf("%#v\n", err)
		return err
	}

	defer c.Close()

	if err = c.Hello("localhost"); err != nil {
		fmt.Printf("%#v\n", err)
		return err
	}

	if err = c.Mail(from); err != nil {
		fmt.Printf("%#v\n", err)
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			fmt.Printf("%#v\n", err)
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		fmt.Printf("%#v\n", err)
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		fmt.Printf("%#v\n", err)
		return err
	}

	b := make([]int, 600)
	for i, _ := range b {
		w.Write([]byte(fmt.Sprintf("yo%d\r\n", i)))
		time.Sleep(1 * time.Second)
	}

	err = w.Close()
	if err != nil {
		fmt.Printf("%#v\n", err)
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
