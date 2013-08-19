//
// syslog-shim
//
// Simple wrapper for haproxy that lets redirects log messages from
// syslog to stderr.
//

package main

import (
	"github.com/ziutek/syslog"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

type handler struct {
	// To simplify implementation of our handler we embed helper
	// syslog.BaseHandler struct.
	*syslog.BaseHandler
}

// Simple fiter for named/bind messages which can be used with BaseHandler
func filter(m *syslog.Message) bool {
	return true
}

func newHandler() *handler {
	h := handler{syslog.NewBaseHandler(5, filter, false)}
	go h.mainLoop() // BaseHandler needs some gorutine that reads from its queue
	return &h
}

// mainLoop reads from BaseHandler queue using h.Get and logs messages to stdout
func (h *handler) mainLoop() {
	for {
		m := h.Get()
		if m == nil {
			break
		}
		h := m.Hostname + " "
		if h == " " {
			h = "app: "
		}
		log.Printf("%s: %s%s: %s%s\n", m.Facility, h, m.Severity, m.Tag, m.Content)
	}
	log.Println("exit handler")
	h.End()
}

func main() {
	prog := filepath.Base(os.Args[0])

	log.SetFlags(0)
	log.SetPrefix(prog + ": ")

	if len(os.Args) < 2 {
		log.Println("usage: " + prog + " address child")
		os.Exit(2)
	}

	address := os.Args[1]

	// Create a server with one handler and run one listen gorutine
	s := syslog.NewServer()
	s.AddHandler(newHandler())
	log.Println("listening on " + address)
	s.Listen(address)

	log.Println("running: " + strings.Join(os.Args[2:], " "))
	c := exec.Command(os.Args[2], os.Args[3:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Start(); err != nil {
		log.Println(err)
		log.Println("shutting down")
		s.Shutdown()
		log.Println("exiting 1")
		os.Exit(1)
	}

	dc := make(chan error)
	go func() {
		dc <- c.Wait()
	}()

	// Wait for terminating signal
	sc := make(chan os.Signal, 3)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2)

	exitCode := 0

L:
	for {
		select {
		case err := <-dc:
			if err != nil {
				if msg, ok := err.(*exec.ExitError); ok {
					exitCode = msg.Sys().(syscall.WaitStatus).ExitStatus()
				}
			}
			log.Printf("child exited %d\n", exitCode)
			break L
		case sig := <-sc:
			log.Printf("received signal %s\n", sig)
			log.Printf("sending signal %s to child\n", sig)
			if err := c.Process.Signal(sig); err != nil {
				log.Println(err)
				log.Println("shutting down")
				s.Shutdown()
				log.Println("exiting 1")
				os.Exit(1)
			}
		}
	}

	log.Println("shutting down")
	s.Shutdown()
	if exitCode < 0 {
		exitCode = 1
	}
	log.Printf("exiting %d\n", exitCode)
	os.Exit(exitCode)
}
