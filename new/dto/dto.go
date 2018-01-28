package dto

// Trigger represent a json used to send via
// rabbitmq and start the scan proccess
type Trigger struct {
	Timestamp string `json:"timestamp"`
	Execution string `json:"execution"`
}

// File
type File struct {
	MD5Sum string `json:"md5sum"`
	Name   string `json:"name"`
	Data   string `json:"data"`
}

type Seed struct {
	Execution string `json:"execution"`
	Files     []File `json:"files"`
	Hostname  string `json:"hostname"`
	Resultdir string `json:"resultdir"`
	TimeStamp string `json:"timestamp"`
	Type      string `json:"type"`
}
