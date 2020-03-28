package exec

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"time"
)

var ErrExecTimeout = fmt.Errorf("exec timed out")

type Cmd struct {
	*exec.Cmd
}

func Command(name string, args ...string) *Cmd {
	return &Cmd{Cmd: exec.Command(name, args...)}
}

func (c *Cmd) GetOutput() (stdout, stderr []byte, err error) {
	return c.GetOutputWithTimeout(-1)
}

func (c *Cmd) GetOutputWithTimeout(d time.Duration) (stdout, stderr []byte, err error) {
	outPipe, _ := c.StdoutPipe()
	errPipe, _ := c.StderrPipe()
	if err = c.Start(); err != nil {
		return
	}

	killed := false
	if d > 0 {
		killer := time.AfterFunc(d, func() {
			c.Process.Kill()
			killed = true
		})
		defer killer.Stop()
	}

	stdout, _ = ioutil.ReadAll(outPipe)
	stderr, _ = ioutil.ReadAll(errPipe)

	if err = c.Wait(); err != nil {
		if killed {
			err = ErrExecTimeout
		}
	}
	return
}
