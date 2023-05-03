package main

type Config struct {
	Id      string            `json:"id"`
	Entries map[string]string `json:"entries"`
}

type Group struct {
	Id      string   `json:"id"`
	Configs []Config `json:"configs"`
}
