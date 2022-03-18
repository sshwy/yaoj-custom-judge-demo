package judger

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
)

type GeneJudger struct {
	Executable string
	Policy     string
}

var _ Judger = (*GeneJudger)(nil)

func (r *GeneJudger) Perform() (JudgeRes, error) {
	var res JudgeRes

	if r.Policy == "" {
		r.Policy = "builtin:free"
	}

	cmd := exec.Command("yaoj-judger", r.Executable,
		"-j", "general", "--log=.log.local", "-p", r.Policy, "--json")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Printf("%s", cmd.String())

	if err := cmd.Run(); err != nil {
		return res, err
	}

	if err := json.Unmarshal(stdout.Bytes(), &res); err != nil {
		return res, err
	}

	return res, nil
}
