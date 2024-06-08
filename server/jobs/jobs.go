package jobs 

import (
	"github.com/ngn13/ezcat/server/util"
)

type Job struct {
  ID      string `json:"id"`
  Message string `json:"message"`
  Success bool   `json:"success"`
  Active  bool   `json:"active"`
}

var jobs []Job

func Add(msg string) *Job {
  var job Job = Job{
    ID: util.MakeRandom(12),
    Message: msg,
    Success: false,
    Active: true,
  }

  jobs = append(jobs, job)
  return &jobs[len(jobs)-1]
}

func Get(id string) *Job {
  for i := range jobs {
    if jobs[i].ID == id {
      return &jobs[i]
    }
  }
  return nil
}

func Del(id string) {
  for i := range jobs {
    if jobs[i].ID == id {
      jobs = append(jobs[:i], jobs[i+1:]...)
      return
    }
  }
}
