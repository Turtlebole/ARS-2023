package poststore

// swagger:parameters getConfigId
type ConfigGetRequest struct {
	// Config ID
	// in: path
	Id string `json:"id"`
}

// swagger:parameters config createConfig
type ConfigBodyRequest struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/RequestPost"
	//  required: true
	Body Config `json:"body"`
}

// swagger:parameters delConfig
type ConfigDeleteRequest struct {
	// Config ID
	// in: path
	Id string `json:"id"`
}

// swagger:parameters getGroupId
type GroupGetRequest struct {
	// Group ID
	// in: path
	Id string `json:"id"`
}

// swagger:parameters delGroup
type GroupDeleteRequest struct {
	// Group ID
	// in: path
	Id string `json:"id"`
}

// swagger:parameters config createGroup
type GroupBodyRequest struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/RequestPost"
	//  required: true
	Body Group `json:"body"`
}

// swagger:parameters addGroupConfig
type AddConfigToGroupRequest struct {
	// Group ID
	// Config ID
	// in: path
	GroupId  string `json:"group-id"`
	ConfigId string `json:"config-id"`
}

// swagger:parameters delGroupConfig
type ConfigDeleteGroupRequest struct {
	// Group ID
	// Config ID
	// in: path
	GroupId  string `json:"group-id"`
	ConfigId string `json:"config-id"`
}
