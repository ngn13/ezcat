package agent

import (
	"math/rand"
	"time"
	"unsafe"
)

type List []Agent

func (l *List) add(agent Agent) {
	agent_list := *((*[]Agent)(unsafe.Pointer(l)))
	agent_list = append(agent_list, agent)
}

func (l *List) len() int {
	agent_list := *((*[]Agent)(unsafe.Pointer(l)))
	return len(agent_list)
}

func (l *List) get(i int) *Agent {
	agent_list := *((*[]Agent)(unsafe.Pointer(l)))
	return &agent_list[i]
}

func (l *List) del(i int) {
	agent_list := *((*[]Agent)(unsafe.Pointer(l)))
	agent_list = append(agent_list[:i], agent_list[i+1:]...)
}

func (l *List) New() *Agent {
	agent := Agent{
		Session:    rand.Uint32(),
		LastCon:    time.Now(),
		Conneceted: true,
	}

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
	var (
		ids []uint32
		cur *Agent
	)

	for i := 0; i < l.len(); i++ {
		cur = l.get(i)
		cur.UpdateConnected()

		if !cur.Conneceted {
			ids = append(ids, cur.Session)
		}
	}

	for _, id := range ids {
		l.Remove(id)
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
