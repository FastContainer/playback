package main

import (
	"flag"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/carlescere/scheduler"
)

func main() {
	flag.Parse()
	args := flag.Args()

	var normal Experiment
	var tarpit Experiment

	if len(args) > 0 && args[0] == "dryrun" {
		normal = &MailMock{from: "root@smtp-client", to: "root@smtp-rcpt", host: "containers:58026"}
		tarpit = &MailMock{from: "root@smtp-client", to: "ca@smtp-tarpit", host: "containers:58025"}
	} else {
		normal = &Mail{from: "root@smtp-client", to: "root@smtp-rcpt", host: "containers:58026"}
		tarpit = &Mail{from: "root@smtp-client", to: "ca@smtp-tarpit", host: "containers:58025"}
	}

	job1 := func() {
		b1, _ := normal.sendToNormal(1, 1)
		fmt.Printf("%s\n", b1)
	}
	job2 := func() {
		b2, _ := tarpit.sendToTarpit(1, 1)
		fmt.Printf("%s\n", b2)
	}

	scheduler.Every(3).Seconds().Run(job1)
	scheduler.Every(2).Seconds().Run(job2)

	runtime.Goexit()
}

type Experiment interface {
	sendToNormal(int, int) ([]byte, error)
	sendToTarpit(int, int) ([]byte, error)
}

type Mail struct {
	from string
	to   string
	host string
}

func (m *Mail) sendToNormal(sessionCount int, messageCount int) ([]byte, error) {
	args := []string{
		"-s", strconv.Itoa(sessionCount),
		"-m", strconv.Itoa(messageCount),
		"-c",
		"-S", "test mail by cluster account",
		"-f", m.from,
		"-t", m.to,
		m.host,
	}
	return exec.Command("smtp-source", args...).Output()
}

func (m *Mail) sendToTarpit(sessionCount int, messageCount int) ([]byte, error) {
	args := []string{
		"-s", strconv.Itoa(sessionCount),
		"-m", strconv.Itoa(messageCount),
		"-c",
		"-N",
		"-S", "cluster account",
		"-f", m.from,
		"-t", m.to,
		m.host,
	}
	return exec.Command("smtp-source", args...).Output()
}

type MailMock struct {
	from string
	to   string
	host string
}

func (m *MailMock) sendToNormal(sessionCount int, messageCount int) ([]byte, error) {
	args := []string{
		"-s", strconv.Itoa(sessionCount),
		"-m", strconv.Itoa(messageCount),
		"-c",
		"-N",
		"-S", "cluster account",
		"-f", m.from,
		"-t", m.to,
		m.host,
	}
	return []byte(fmt.Sprintln(args)), nil
}

func (m *MailMock) sendToTarpit(sessionCount int, messageCount int) ([]byte, error) {
	args := []string{
		"-s", strconv.Itoa(sessionCount),
		"-m", strconv.Itoa(messageCount),
		"-c",
		"-N",
		"-S", "cluster account",
		"-f", m.from,
		"-t", m.to,
		m.host,
	}
	return []byte(fmt.Sprintln(args)), nil
}
