package main

import (
	"github.com/Microsoft/go-winio"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

func main() {
	var err error
	if len(os.Args) > 1 && os.Args[1] == "child" {
		log.SetPrefix("child: ")
		err = runChild()
	} else {
		log.SetPrefix("parent: ")
		err = runParent()
	}
	if err != nil {
		log.Fatal(err)
	}
}

var namedPipe = `\\.\pipe\cnb_exec_d`

func runParent() error {
	log.Printf("opening pipe for writing %s\n", namedPipe)

	l, err := winio.ListenPipe(namedPipe, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to create stdin pipe %s", namedPipe)
	}
	defer func(l net.Listener) {
		log.Printf("deferred closing of pipe %s\n", namedPipe)
		if err != nil {
			l.Close()
		}
	}(l)

	go func() {
		log.Printf("listening on pipe %s\n", namedPipe)
		c, err := l.Accept()
		if err != nil {
			log.Printf("failed to accept stdin connection on %s\n", namedPipe)
			return
		}

		body, err := ioutil.ReadAll(c)
		if err != nil {
			log.Printf("failed to read %s\n", namedPipe)
			return
		}

		// print out child's pipe message to parent's stdout
		log.Printf("pipe content %s\n", string(body))

		c.Close()
		l.Close()
	}()

	log.Printf("running child subprocess %s\n", namedPipe)
	cmd := exec.Command(os.Args[0], "child")
	// attach stdout/err so child log debug messages show up
	// NOTE: actual pipe messages won't be logged as "child"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)

	return nil
}

func runChild() error {
	log.Printf("dialing pipe %s\n", namedPipe)
	d, err := winio.DialPipe(namedPipe, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to open pipe for writing %s", namedPipe)
	}

	log.Printf("writing message to pipe %s\n", namedPipe)
	if _, err = d.Write([]byte("Hello World")); err != nil {
		return errors.Wrapf(err, "failed to write to pipe %s", namedPipe)
	}

	if err = d.Close(); err != nil {
		return errors.Wrapf(err, "failed to close dialer for pipe %s", namedPipe)
	}

	return nil
}
