package status

type Status string

const (
	InProgress Status = "In Progress"
	Success    Status = "Success"
	Warning    Status = "Warning"
	Failed     Status = "Failed"
	Deleted    Status = "Deleted"
)
