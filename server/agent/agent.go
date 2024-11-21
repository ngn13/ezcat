package agent

import (
	"time"

	"github.com/ngn13/ezcat/server/util"
)

const AGENT_SLEEP_MAX = 15

type Agent struct {
	Session uint32 `json:"session"`

	// sysem info
	Username string `json:"username"`
	Hostname string `json:"hostname"`
	PID      int32  `json:"pid"`
	OS       string `json:"os"`
	IP       string `json:"ip"`

	LastCon    time.Time `json:"last_con"` // what was the last connection time
	Conneceted bool      `json:"-"`        // is it currently connected
	ShouldKill bool      `json:"-"`

	Job []Job `json:"-"`
}

func (a *Agent) UpdateConnected() {
	cur := time.Now()
	a.Conneceted = false

	if a.ShouldKill {
		return
	}

	if diff := cur.Sub(a.LastCon); diff.Seconds() < AGENT_SLEEP_MAX {
		a.Conneceted = true
	}
}

func (a *Agent) NextJob() *Job {
	var job *Job = nil

	for i := range a.Job {
		if job = &a.Job[i]; !job.Waiting {
			job = nil
			continue
		}
	}

	return job
}

func (a *Agent) AddJob(cmd byte, arg []byte, arg_size uint8, cb func(*Job)) *Job {
	if nil != arg && arg_size == 0 {
		arg_size = uint8(len(arg))
	}

	if nil == arg {
		arg_size = 0
	}

	job := Job{
		ID:           util.Rand16(),
		Command:      cmd,
		Waiting:      true,
		Success:      false,
		Argument:     arg,
		ArgumentSize: arg_size,
		Callback:     cb,
		Agent:        a,
	}

	a.Job = append(a.Job, job)
	return &a.Job[len(a.Job)-1]
}

func (a *Agent) GetJob(id uint16) *Job {
	for i := range a.Job {
		if a.Job[i].ID == id {
			return &a.Job[i]
		}
	}

	return nil
}

func (a *Agent) DelJob(id uint16) {
	for i := range a.Job {
		if a.Job[i].ID == id {
			a.Job = append(a.Job[:i], a.Job[i+1:]...)
			return
		}
	}
}
