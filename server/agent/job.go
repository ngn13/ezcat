package agent

type Job struct {
	ID           uint16
	Command      byte
	Waiting      bool
	Success      bool
	Argument     []byte
	ArgumentSize uint8
	Response     []byte
	ResponseSize uint8
	Callback     func(*Job)
	Agent        *Agent
}
