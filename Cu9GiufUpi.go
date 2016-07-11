package main

import (
	"errors"
	"strings"

	"github.com/emicklei/go-restful"
)

func main() {
	SetupCrudRoutes()
}

type WorkspaceObject struct {
	TableName     string
	ID            string
	Slug          string
}

type CrudObject struct {
	ID           *uint32 `json:"ID,omitempty" db:"ID"`
	Name         *string `json:"Name,omitempty" db:"Name"`
}

var wo = []WorkspaceObject{
	{"Profiles", "ProfilesID", "profiles"},
	{"Dashboards", "DashboardsID", "dashboards"},
	{"Filtersets", "FiltersetsID", "filtersets"},
	{"KeywordGroups", "KeywordGroupsID", "keyword-groups"},
	{"Teams", "TeamsID", "teams"},
}

func SetupCrudRoutes() {

	ws := new(restful.WebService)
	ws.
		Path("/workspaces/{WorkspacesID}/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	for _, object := range wo {
		SetupWorkspaceObjectCrudRoutes(ws, &object)
	}

	restful.Add(ws)
}

func SetupWorkspaceObjectCrudRoutes(ws *restful.WebService, wo *WorkspaceObject) {

	ws.Route(ws.GET("/" + wo.Slug + "/").To(CrudIndex).
		Doc("get a profile").
		Writes(Profile{})) // on the response

	ws.Route(ws.GET("/" + wo.Slug + "/{" + wo.ID + "}/").To(CrudView).
		Doc("get a profile").
		Param(ws.PathParameter(wo.ID, "identifier of the profile").DataType("uint32")).
		Writes(Profile{})) // on the response

	ws.Route(ws.PATCH("/" + wo.Slug + "/{" + wo.ID + "}/").To(CrudUpdate).
		Doc("update a profile").
		Reads(Profile{})) // from the request

	ws.Route(ws.POST("/" + wo.Slug + "/").To(CrudAdd).
		Doc("create a profile").
		Param(ws.PathParameter(wo.ID, "identifier of the profile").DataType("uint32")).
		Reads(Profile{})) // from the request

	ws.Route(ws.DELETE("/" + wo.Slug + "/{" + wo.ID + "}/").To(CrudRemove).
		Doc("delete a profile").
		Param(ws.PathParameter(wo.ID, "identifier of the profile").DataType("uint32")))
}

func CrudIndex(r *restful.Request, w *restful.Response) {
	log.Println("I haven't been written yet!")
}

func CrudView(r *restful.Request, w *restful.Response) {
	log.Println("I haven't been written yet!")
}

func CrudAdd(r *restful.Request, w *restful.Response) {
	log.Println("I haven't been written yet!")
}

func CrudUpdate(r *restful.Request, w *restful.Response) {
	log.Println("I haven't been written yet!")
}

func CrudRemove(r *restful.Request, w *restful.Response) {
	log.Println("I haven't been written yet!")
}