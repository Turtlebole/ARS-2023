package main

import (
	"errors"
	"mime"
	"net/http"

	"github.com/gorilla/mux"
)

type configServer struct {
	data      map[string]*Config
	groupData map[string]*Group // izigrava bazu podataka
}

// swagger:route POST /config/ config createConfig
// Add a new config
//
// responses:
//
//	415: ErrorResponse
//	400: ErrorResponse
//	201: ResponseConfig
func (ts *configServer) createConfigHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := createId()
	rt.Id = id
	ts.data[id] = rt
	renderJSON(w, rt)
}

// swagger:route GET /configs/ config getConfigs
// Get configs
//
// responses:
//
//	200: []ResponseConfig
func (ts *configServer) getAllHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := []*Config{}
	for _, v := range ts.data {
		allTasks = append(allTasks, v)
	}

	renderJSON(w, allTasks)
}

// swagger:route GET /config/{id}/ config getConfigId
// Get config Id
//
// responses:
//
//	404: ErrorResponse
//	200: ResponseConfig
func (ts *configServer) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, ok := ts.data[id]
	if !ok {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

// swagger:route DELETE /config/{id}/ config delConfig
// Delete config
//
// responses:
//
//	404: ErrorResponse
//	204: NoContentResponse
func (ts *configServer) delConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if v, ok := ts.data[id]; ok {
		delete(ts.data, id)
		renderJSON(w, v)
	} else {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

// swagger:route POST /group/ group createGroup
// Add a new group
//
// responses:
//
//	415: ErrorResponse
//	400: ErrorResponse
//	201: ResponseGroup
func (ts *configServer) createGroupHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	group, err := decodeGroup(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := createId()
	group.Id = id
	ts.groupData[id] = group
	renderJSON(w, group)
}

// swagger:route PUT /group/{id}/config/{id}/ group addGroupConfig
// Add config to group
//
// responses:
//
//	415: ErrorResponse
//	400: ErrorResponse
//	201: ResponseGroup
func (ts *configServer) addGroupConfig(w http.ResponseWriter, req *http.Request) {
	groupId := mux.Vars(req)["groupId"]
	id := mux.Vars(req)["id"]
	task, ok := ts.data[id]
	group, ook := ts.groupData[groupId]
	if !ok || !ook {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	group.Configs = append(group.Configs, *task)
	ts.groupData[groupId] = group

	return
}

// swagger:route GET /groups/ group getGroups
// Get all groups
//
// responses:
//
//	200: []ResponseGroup
func (ts *configServer) getAllGroupsHandler(w http.ResponseWriter, req *http.Request) {
	allGroups := []*Group{}
	for _, v := range ts.groupData {
		allGroups = append(allGroups, v)
	}

	renderJSON(w, allGroups)
}

// swagger:route GET /group/{id}/ group getGroupId
// Get group Id
//
// responses:
//
//	404: ErrorResponse
//	200: ResponseGroup
func (ts *configServer) getGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, ok := ts.groupData[id]
	if !ok {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

// swagger:route DELETE /group/{id}/ group delGroup
// Delete group
//
// responses:
//
//	404: ErrorResponse
//	204: NoContentResponse
func (ts *configServer) delGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	_, ok := ts.groupData[id]
	if !ok {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	delete(ts.groupData, id)
}

// swagger:route DELETE /group/{id}/config/{id}/ group delGroupConfig
// Delete config from group
//
// responses:
//
//	404: ErrorResponse
//	204: NoContentResponse
func (ts *configServer) delGroupHandlerConfig(w http.ResponseWriter, req *http.Request) {
	groupId := mux.Vars(req)["groupId"]
	id := mux.Vars(req)["id"]
	group, ok := ts.groupData[groupId]
	if !ok {
		err := errors.New("group not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	for i, config := range group.Configs {
		if config.Id == id {
			group.Configs = append(group.Configs[:i], group.Configs[i+1:]...)
			ts.groupData[groupId] = group
			return
		}
	}

	err := errors.New("config not found in group")
	http.Error(w, err.Error(), http.StatusNotFound)
	return
}

func (ts *configServer) swaggerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./swagger.yaml")
}
