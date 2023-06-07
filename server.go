package main

import (
	"errors"
	"mime"
	"net/http"

	ps "github.com/Turtlebole/ARS-2023/poststore"

	"github.com/gorilla/mux"
)

type postServer struct {
	store *ps.PostStore
	//data      map[string]*ps.Config
	//groupData map[string]*ps.Group
}

// swagger:route POST /config/ config createConfig
// Add a new config
//
// responses:
//
//	415: ErrorResponse
//	400: ErrorResponse
//	201: ResponseConfig
func (ts *postServer) createConfigHandler(w http.ResponseWriter, req *http.Request) {
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

	task, err := ts.store.Post(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	renderJSON(w, task)
}

// swagger:route GET /configs/ config getConfigs
// Get configs
//
// responses:
//
//	200: []ResponseConfig
func (ts *postServer) getAllHandler(w http.ResponseWriter, req *http.Request) {
	allTasks, err := ts.store.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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
func (ts *postServer) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	task, err := ts.store.Get(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
func (ts *postServer) delConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]

	msg, err := ts.store.Delete(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, msg)
}

// swagger:route POST /group/ group createGroup
// Add a new group
//
// responses:
//
//	415: ErrorResponse
//	400: ErrorResponse
//	201: ResponseGroup
func (ts *postServer) createGroupHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	group, err := decodeGroup(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := ts.store.PostGroup(group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	renderJSON(w, task)
}

// swagger:route PUT /group/{id}/config/{id}/ group addGroupConfig
// Add config to group
//
// responses:
//
//	415: ErrorResponse
//	400: ErrorResponse
//	201: ResponseGroup
func (ts *postServer) addGroupConfigHandler(w http.ResponseWriter, req *http.Request) {
	groupId := mux.Vars(req)["Id"]
	configId := mux.Vars(req)["Id"]
	err := ts.store.AddConfigToGroup(groupId, configId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

//	 swagger:route GET /groups/ group getGroups
//	 Get all groups
//
//	 responses:
//
//		200: []ResponseGroup
func (ts *postServer) getAllGroupsHandler(w http.ResponseWriter, req *http.Request) {
	allGroups, err := ts.store.GetAllGroups()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderJSON(w, allGroups)
}

//	 swagger:route GET /group/{id}/ group getGroupId
//	 Get group Id
//
//	 responses:
//
//		404: ErrorResponse
//		200: ResponseGroup
func (ts *postServer) getGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]

	group, err := ts.store.GetGroupById(id, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, group)
}

// // swagger:route DELETE /group/{id}/ group delGroup
// // Delete group
// //
// // responses:
// //
// //	404: ErrorResponse
// //	204: NoContentResponse
func (ts *postServer) delGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["Id"]
	version := mux.Vars(req)["version"]

	msg, err := ts.store.DeleteGroup(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, msg)
}

// 	delete(ts.groupData, id)
// }

// // swagger:route DELETE /group/{id}/config/{id}/ group delGroupConfig
// // Delete config from group
// //
// // responses:
// //
// //	404: ErrorResponse
// //	204: NoContentResponse
// func (ts *postServer) delGroupHandlerConfig(w http.ResponseWriter, req *http.Request) {
// 	groupId := mux.Vars(req)["groupId"]
// 	id := mux.Vars(req)["id"]
// 	group, ok := ts.groupData[groupId]
// 	if !ok {
// 		err := errors.New("group not found")
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	for i, config := range group.Configs {
// 		if config.Id == id {
// 			group.Configs = append(group.Configs[:i], group.Configs[i+1:]...)
// 			ts.groupData[groupId] = group
// 			return
// 		}
// 	}

// 	err := errors.New("config not found in group")
// 	http.Error(w, err.Error(), http.StatusNotFound)
// 	return
// }

func (ts *postServer) swaggerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./swagger.yaml")
}
