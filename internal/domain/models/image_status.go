package models

type Status string

const (
	InProgress Status = "In Progress"
	Success           = "Success"
	Warning           = "Warning"
	Failed            = "Failed"
)
