//go:build windows

package game

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gocarina/gocsv"
)

const (
	noTasksMatchingCriteria = "INFO: No tasks are running which match the specified criteria."
	statusNotResponding     = "NOT RESPONDING"
)

var errorWindowTitles = [4]string{
	"BF2 Memory Error",
	"BF2 Error",
	"Microsoft Visual C++ Runtime Library",
}

// "Image Name","PID","Session Name","Session#","Mem Usage","Status","User Name","CPU Time","Window Title"
type taskListProc struct {
	Name        string `csv:"Image Name"`
	PID         int    `csv:"PID"`
	SessionName string `csv:"Session Name"`
	SessionID   int    `csv:"Session#"`
	Memory      string `csv:"Mem Usage"`
	Status      string `csv:"Status"`
	UserName    string `csv:"User Name"`
	CPUTime     string `csv:"CPU Time"`
	WindowTitle string `csv:"Window Title"`
}

func processHealth(pid int) (bool, error) {
	args := []string{
		"/FO", "csv",
		"/FI", fmt.Sprintf("PID eq %d", pid),
	}

	cmd := exec.Command("tasklist", args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return false, err
	}

	output := out.String()

	if strings.Contains(output, noTasksMatchingCriteria) {
		return false, nil
	}

	var procs []taskListProc
	if strings.Contains(output, "PID") {
		err := gocsv.Unmarshal(&out, &procs)
		if err != nil {
			return false, err
		}
	}

	if len(procs) == 0 {
		return false, nil
	}

	proc := procs[0]

	if proc.Status == statusNotResponding {
		return false, nil
	}

	for _, title := range errorWindowTitles {
		if strings.Contains(proc.WindowTitle, title) {
			return false, nil
		}
	}

	return true, nil
}
