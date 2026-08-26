package dto

type Event struct {
	SuiteName string `json:"suiteName"`
	Path      string `json:"path"`
	Event     string `json:"event"`
	Contents  []byte `json:"contents"`
	FileType  string `json:"fileType"`
	RenamedTo string `json:"renamedTo"`
	Ip        string `json:"ip"`
}

type Suite struct {
	Name  string  `json:"name"`
	Files []Event `json:"files"`
}

type NewSuiteReq struct {
	Suites []Suite `json:"suites"`
}

type SuiteReq struct {
	Suites []Suite `json:"suites"`
}

type EventsReq struct {
	Events []Event `json:"events"`
}
