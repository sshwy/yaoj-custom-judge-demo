// 这个包包含了 yaoj-judger 接口和一个基本的实现
package judger

import (
	"fmt"
)

type JudgeRes struct {
	Result  int `json:"result"`
	Signal  int `json:"signal"`
	Ecode   int `json:"exit_code"`
	Rtime   int `json:"real_time"`
	Ctime   int `json:"cpu_time"`
	Rmemory int `json:"memory"`
}

func (r JudgeRes) String() string {
	return fmt.Sprintf("JudgeRes{{ result: %d, signal: %d, ecode: %d, time(real/cpu): %d/%d, memory: %d }}", r.Result,
		r.Signal, r.Ecode, r.Rtime, r.Ctime, r.Rmemory)
}

// 一个 JudgerTask 代表一个评测任务，其中 Perform 方法用于执行评测，返回相应的评测结果或者错误。
type Judger interface {
	Perform() (JudgeRes, error)
}
