package main

import (
	&#34;errors&#34;
	&#34;strings&#34;

	&#34;github.com/emicklei/go-restful&#34;
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
	ID           *uint32 `json:&#34;ID,omitempty&#34; db:&#34;ID&#34;`
	Name         *string `json:&#34;Name,omitempty&#34; db:&#34;Name&#34;`
}

var wo = []WorkspaceObject{
	{&#34;Profiles&#34;, &#34;ProfilesID&#34;, &#34;profiles&#34;},
	{&#34;Dashboards&#34;, &#34;DashboardsID&#34;, &#34;dashboards&#34;},
	{&#34;Filtersets&#34;, &#34;FiltersetsID&#34;, &#34;filtersets&#34;},
	{&#34;KeywordGroups&#34;, &#34;KeywordGroupsID&#34;, &#34;keyword-groups&#34;},
	{&#34;Teams&#34;, &#34;TeamsID&#34;, &#34;teams&#34;},
}

func SetupCrudRoutes() {

	ws := new(restful.WebService)
	ws.
		Path(&#34;/workspaces/{WorkspacesID}/&#34;).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	for _, object := range wo {
		SetupWorkspaceObjectCrudRoutes(ws, &amp;object)
	}

	restful.Add(ws)
}

func SetupWorkspaceObjectCrudRoutes(ws *restful.WebService, wo *WorkspaceObject) {

	ws.Route(ws.GET(&#34;/&#34; &#43; wo.Slug &#43; &#34;/&#34;).To(CrudIndex).
		Doc(&#34;get a profile&#34;).
		Writes(Profile{})) // on the response

	ws.Route(ws.GET(&#34;/&#34; &#43; wo.Slug &#43; &#34;/{&#34; &#43; wo.ID &#43; &#34;}/&#34;).To(CrudView).
		Doc(&#34;get a profile&#34;).
		Param(ws.PathParameter(wo.ID, &#34;identifier of the profile&#34;).DataType(&#34;uint32&#34;)).
		Writes(Profile{})) // on the response

	ws.Route(ws.PATCH(&#34;/&#34; &#43; wo.Slug &#43; &#34;/{&#34; &#43; wo.ID &#43; &#34;}/&#34;).To(CrudUpdate).
		Doc(&#34;update a profile&#34;).
		Reads(Profile{})) // from the request

	ws.Route(ws.POST(&#34;/&#34; &#43; wo.Slug &#43; &#34;/&#34;).To(CrudAdd).
		Doc(&#34;create a profile&#34;).
		Param(ws.PathParameter(wo.ID, &#34;identifier of the profile&#34;).DataType(&#34;uint32&#34;)).
		Reads(Profile{})) // from the request

	ws.Route(ws.DELETE(&#34;/&#34; &#43; wo.Slug &#43; &#34;/{&#34; &#43; wo.ID &#43; &#34;}/&#34;).To(CrudRemove).
		Doc(&#34;delete a profile&#34;).
		Param(ws.PathParameter(wo.ID, &#34;identifier of the profile&#34;).DataType(&#34;uint32&#34;)))
}

func CrudIndex(r *restful.Request, w *restful.Response) {
	log.Println(&#34;I haven&#39;t been written yet!&#34;)
}

func CrudView(r *restful.Request, w *restful.Response) {
	log.Println(&#34;I haven&#39;t been written yet!&#34;)
}

func CrudAdd(r *restful.Request, w *restful.Response) {
	log.Println(&#34;I haven&#39;t been written yet!&#34;)
}

func CrudUpdate(r *restful.Request, w *restful.Response) {
	log.Println(&#34;I haven&#39;t been written yet!&#34;)
}

func CrudRemove(r *restful.Request, w *restful.Response) {
	log.Println(&#34;I haven&#39;t been written yet!&#34;)
}