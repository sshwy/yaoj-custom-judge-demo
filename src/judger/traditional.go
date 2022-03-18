package judger

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
)

// 目前只考虑 builtin policy
type TradJudger struct {
	Executable string
	InpFile    string
	OutFile    string
	ErrFile    string
	Fileio     bool
	Policy     string
}

var _ Judger = (*TradJudger)(nil)

func (r *TradJudger) Perform() (JudgeRes, error) {
	var res JudgeRes

	fileFlag := "std"
	if r.Fileio {
		fileFlag = "file"
	}

	if r.Policy == "" {
		r.Policy = "builtin:free"
	}

	cmd := exec.Command("yaoj-judger", r.Executable, r.InpFile, r.OutFile, r.ErrFile, fileFlag,
		"-j", "traditional", "--log=.log.local", "-p", r.Policy, "--json")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Printf("%s", cmd.String())

	if err := cmd.Run(); err != nil {
		return res, err
	}
	// log.Printf("[stdout] %s", strings.Trim(stdout.String(), " \n\t\r"))
	// log.Printf("[stderr] %s", strings.Trim(stderr.String(), " \n\t\r"))

	if err := json.Unmarshal(stdout.Bytes(), &res); err != nil {
		return res, err
	}

	return res, nil
}
