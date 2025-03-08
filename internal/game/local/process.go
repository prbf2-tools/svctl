package local

import (
	"context"
	"errors"
	"os"

	"github.com/shirou/gopsutil/v3/process"
)

var (
	ErrProcessAlreadyRunning = errors.New("process already running")
)

var commonProcessArgs = []string{
	"+modPath", "mods/pr",
	"+noStatusMonitor", "1",
	"+multi", "1",
	"+dedicated", "1",
}

type spawnedProcess struct {
	proc *os.Process

	ctx context.Context
}

func watchProcess(proc *os.Process) *spawnedProcess {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		_, _ = proc.Wait()
		cancel()
	}()

	return &spawnedProcess{
		proc: proc,
		ctx:  ctx,
	}
}

func (p *spawnedProcess) PID() int {
	return p.proc.Pid
}

func (p *spawnedProcess) Kill() error {
	return p.proc.Kill()
}

func (p *spawnedProcess) IsRunning() (bool, error) {
	select {
	case <-p.ctx.Done():
		return false, nil
	default:
		return true, nil
	}
}

type recoveredProcess struct {
	processExe string
	pid        int
}

func recoverProcess(pid int, processExe string) *recoveredProcess {
	return &recoveredProcess{
		processExe: processExe,
		pid:        pid,
	}
}

func (p *recoveredProcess) PID() int {
	return p.pid
}

func (p *recoveredProcess) Kill() error {
	proc, err := process.NewProcess(int32(p.pid))
	if err != nil {
		return nil
	}

	return proc.Kill()
}

func (p *recoveredProcess) IsRunning() (bool, error) {
	proc, err := process.NewProcess(int32(p.pid))
	if err != nil {
		return false, nil
	}

	procExe, err := proc.Exe()
	if err != nil {
		return false, nil
	}

	if procExe != p.processExe {
		return false, nil
	}

	return proc.IsRunning()
}

func makeFileExecutable(exePath string) error {
	info, err := os.Stat(exePath)
	if err != nil {
		return err
	}

	if info.Mode().Perm()&0100 == 0 {
		err = os.Chmod(exePath, info.Mode()|0100)
		if err != nil {
			return err
		}
	}

	return nil
}
