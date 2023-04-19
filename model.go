package main

type Config struct {
	Id      string            `json:"id"`
	Title   string            `json:"title"`
	Text    string            `json:"text"`
	Tags    []string          `json:"tags"`
	Entries map[string]string `json:"entries"`
}

type Service struct {
	Data map[string]*Config `json:"data"`
}
