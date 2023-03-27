package main

import (
	"bytes"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/TKMAX777/RDPRelativeInput/debug"
)
import "C"

func StartApplication(rw *VirtualChannelReadWriteCloser, serverName string) {
	var stderr = new(bytes.Buffer)
	var cmd = exec.Command(os.Getenv("ProgramW6432") + `\RDPRelativeInput\RelativeInputClient.exe`)

	cmd.Stdin = rw
	cmd.Stdout = rw
	cmd.Stderr = stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	cmd.Env = append(
		os.Environ(),
		"SERVER_NAME="+serverName,
	)

	err := cmd.Start()
	if err != nil {
		debug.Debugln("CmdStartError", err)
	}

	var isActive = true

	var MustReadln = func() string {
		var b = make([]byte, 1)
		var res = make([]byte, 0, 100)

		for isActive {
			n, _ := os.Stderr.Read(b)
			if n == 0 {
				time.Sleep(200 * time.Millisecond)
			}
			res = append(res, b[:n]...)
			if bytes.Contains(b[:n], []byte{'\n'}) {
				break
			}
		}
		if !isActive {
			return ""
		}

		res = bytes.TrimSuffix(res, []byte{'\n'})
		return string(res)
	}

	var commandChan = make(chan string)
	go func() {
		for isActive {
			commandChan <- MustReadln()
		}
	}()

	var doneChan = make(chan bool)
	go func() {
		cmd.Wait()
		doneChan <- true
		close(doneChan)
	}()

	for isActive {
		select {
		case command := <-commandChan:
			switch command {
			case "CLOSE":
				isActive = false

				// stop read stderr routine
				debug.Debugf("Stop stderr routine...")
				stderr.Write([]byte("done\n"))
				<-commandChan
				close(commandChan)
				debug.Debugln("ok")

				debug.Debugf("Kill Application...")
				cmd.Process.Kill()
				<-doneChan
				debug.Debugln("ok")

				return
			default:
				debug.Debugln("Application: ", command)
			}
		case <-doneChan:
			isActive = false

			// stop read stderr routine
			debug.Debugf("Stop stderr routine...")
			stderr.Write([]byte("done\n"))
			<-commandChan
			close(commandChan)
			debug.Debugln("ok")

			return
		}
	}
}
