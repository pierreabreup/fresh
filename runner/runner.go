package runner

import (
	"io"
	"os/exec"
)

func run() bool {
	runnerLog("Running...")

  root := root()
  executableName = 'main'
  executableFile := make([]byte, len(root)+len(executableName))
  copyIndex := 0
  copyIndex += copy(executableFile[copyIndex:],root)
  copyIndex += copy(executableFile[copyIndex:],executableName)

	cmd := exec.Command(string(executableFile))

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}

	go io.Copy(appLogWriter{}, stderr)
	go io.Copy(appLogWriter{}, stdout)

	go func() {
		<-stopChannel
		pid := cmd.Process.Pid
		runnerLog("Killing PID %d", pid)
		cmd.Process.Kill()
	}()

	return true
}
