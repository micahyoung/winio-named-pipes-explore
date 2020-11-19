package main

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
)

const handleKey = "FOO_HANDLE"

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

func runParent() error {
	fh, err := windows.Open(`CON`, os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "open file handle")
	}
	defer windows.Close(fh)
	val := fmt.Sprintf("%d", fh)
	fmt.Println("VAL", val)

	f := os.NewFile(uintptr(fh), "foo")
	if f == nil {
		return errors.New("open file from handle")
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, "PARENT test\n"); err != nil {
		return err
	}

	cmd := exec.Command(os.Args[0], "child")
	// discard stout/stderr. only messages come from writes to parsed handle
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", handleKey, val))

	if err := cmd.Run(); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(f, "PARENT after\n"); err != nil {
		return err
	}


	return nil
}

func runChild() error {
	val := os.Getenv(handleKey)
	if val == "" {
		return fmt.Errorf("getting handle value for %s", handleKey)
	}

	fhi, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("parsing handle value %d", fhi)
	}

	f := os.NewFile(uintptr(windows.Handle(fhi)), "foo")
	if f == nil {
		return errors.New("open file from handle")
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, "CHILD test\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(os.Stdout, "STDOUT test\n"); err != nil {
		return err
	}

	return nil
}
