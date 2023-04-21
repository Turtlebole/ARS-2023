package main

type Config struct {
	Id      string            `json:"id"`
	Entries map[string]string `json:"entries"`
}

type Service struct {
	Data map[string]*Config
}
