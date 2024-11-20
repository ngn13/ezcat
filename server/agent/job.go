package agent

const (
	CMD_RUN  byte = 'R'
	CMD_KILL byte = 'K'
)

type Job struct {
	ID       uint32     `json:"id"`
	Command  byte       `json:"cmd"`
	Waiting  bool       `json:"waiting"`
	Success  bool       `json:"success"`
	Argument string     `json:"argument"`
	Response string     `json:"response"`
	Callback func(*Job) `json:"-"`
	Agent    *Agent     `json:"-"`
}
