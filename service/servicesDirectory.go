package service

type ServicesDirectory struct {
	CurrentVersion float32       `json:"currentVersion"`
	Folders        []interface{} `json:"folders"`
	Services       []Service     `json:"services"`
}

type Service struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
