package agent

import (
	"math/rand"
	"time"

	"github.com/ngn13/ezcat/server/log"
)

type List struct {
	Agents []Agent
}

func (l *List) add(agent Agent) {
	l.Agents = append(l.Agents, agent)
}

func (l *List) len() int {
	return len(l.Agents)
}

func (l *List) get(i int) *Agent {
	return &l.Agents[i]
}

func (l *List) del(i int) {
	l.Agents = append(l.Agents[:i], l.Agents[i+1:]...)
}

func (l *List) New() *Agent {
	agent := Agent{
		Session:    rand.Uint32(),
		LastCon:    time.Now(),
		Conneceted: true,
	}

	log.Debg("registered a new agent (session: %v)", agent.Session)

	l.add(agent)
	return l.get(l.len() - 1)
}

func (l *List) Find(s uint32) *Agent {
	for i := 0; i < l.len(); i++ {
		if cur := l.get(i); cur.Session == s {
			return cur
		}
	}

	return nil
}

func (l *List) Remove(s uint32) {
	for i := 0; i < l.len(); i++ {
		if cur := l.get(i); cur.Session == s {
			l.del(i)
			return
		}
	}
}

func (l *List) Update() {
	var cur *Agent

	for i := 0; i < l.len(); i++ {
		cur = l.get(i)
		cur.UpdateConnected()
	}
}

func (l *List) GetJob(id uint16) (job *Job) {
	var cur *Agent

	for i := 0; i < l.len(); i++ {
		cur = l.get(i)

		if job = cur.GetJob(id); job != nil {
			return job
		}
	}

	return nil
}

func (l *List) DelJob(id uint16) {
	for i := 0; i < l.len(); i++ {
		l.get(i).DelJob(id)
	}
}

func (l *List) Ready() []*Agent {
	var res []*Agent

	for i := 0; i < l.len(); i++ {
		if cur := l.get(i); cur.Conneceted && cur.Hostname != "" && cur.Username != "" {
			res = append(res, cur)
		}
	}

	return res
}
