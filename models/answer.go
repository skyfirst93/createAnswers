package models

type Answer struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Event struct {
	Type   string `json:"event"`
	Answer Answer `json:"data"`
}
