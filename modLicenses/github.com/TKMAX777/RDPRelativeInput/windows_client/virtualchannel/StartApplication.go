package main

import (
	"bytes"
	"os"
	"os/exec"
	"syscall"

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
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env,
		"SERVER_NAME="+serverName,
	)

	err := cmd.Start()
	if err != nil {
		debug.Debugln("CmdStartError", err)
	}

	var MustReadln = func() string {
		var b = make([]byte, 1)
		var res = make([]byte, 0, 100)

		for {
			n, _ := os.Stderr.Read(b)
			res = append(res, b[:n]...)
			if bytes.Contains(b[:n], []byte{'\n'}) {
				break
			}
		}

		res = bytes.TrimSuffix(res, []byte{'\n'})

		return string(res)
	}

	var isActive = true
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

	for {
		select {
		case command := <-commandChan:
			switch command {
			case "CLOSE":
				isActive = false
				stderr.Write([]byte("done\n"))
				debug.Debugf("Kill Application...")
				cmd.Process.Kill()
				<-doneChan
				debug.Debugln("ok")
				close(commandChan)
				return
			default:
				debug.Debugln("Application: ", command)
			}
		case <-doneChan:
			close(commandChan)
			return
		}
	}
}
