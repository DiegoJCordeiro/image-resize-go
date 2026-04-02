package models

type ProcessType string

const (
	RESIZE    ProcessType = "RESIZE"
	REFORMAT  ProcessType = "REFORMAT"
	REQUALITY ProcessType = "REQUALITY"
)
