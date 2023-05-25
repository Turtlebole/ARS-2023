package poststore

type Config struct {
	Id      string            `json:"id"`
	Entries map[string]string `json:"entries"`
	Labels  string            `json:"labels"`
	Version string            `json:"version"`
}

type Group struct {
	Id      string   `json:"id"`
	Configs []Config `json:"configs"`
	Labels  string   `json:"labels"`
	Version string   `json:"version"`
}
type RequestPost struct {
	Id      string   `json:"id"`
	Title   string   `json:"title"`
	Text    string   `json:"text"`
	Tags    []string `json:"tags"`
	Labels  string   `json:"labels"`
	Version string   `json:"version"`
}
