package main

// swagger:response ResponseConfig
type ResponseConfig struct {
	// Id of the config
	// in: body
	Id string `json:"id"`
	// Map of config entries
	// in: body
	Entries map[string]string `json:"entries"`
}

// swagger:response ResponseGroup
type ResponseGroup struct {
	// Id of the group
	// in: string
	Id string `json:"id"`

	// List of group configs
	// in: []Config
	Configs []Config `json:"configs"`
}

// swagger:response ErrorResponse
type ErrorResponse struct {
	// Error status code
	// in: int64
	Status int64 `json:"status"`
	// Message of the error
	// in: string
	Message string `json:"message"`
}

// swagger:response NoContentResponse
type NoContentResponse struct{}
