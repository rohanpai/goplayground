package controllers

import (
	"fmt"
	"github.com/revel/revel"
	"my/app/models"
)

// ApiWrapper Controller that wraps all JSON replies and errors.
type ApiWrapper struct {
	*revel.Controller
}

// Reply is boilerplate for all JSON replies
type Reply struct {
	Context string      `json:"context,omitempty"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
}

// Unauthorized generates a 401 with an optional error message.
// If none is supplied a default one is added.
func (c ApiWrapper) Unauthorized(msg string, objs ...interface{}) revel.Result {
	return c.renderErrorString(401, fmt.Sprintf(msg, objs))
}

// Forbidden generates a 403 with an optional error message.
// If none is supplied a default one is added.
func (c ApiWrapper) Forbidden(msg string, objs ...interface{}) revel.Result {
	return c.renderErrorString(403, fmt.Sprintf(msg, objs))
}

// NotFound generates a 404 with an optional error message.
// If none is supplied a default one is added.
func (c ApiWrapper) NotFound(msg string, objs ...interface{}) revel.Result {
	return c.renderErrorString(404, fmt.Sprintf(msg, objs))
}

// BadRequest generates a 400 with an optional error message.
// If none is supplied a default one is added.
// It is recommened to add a meaningful error message
func (c ApiWrapper) BadRequest(msg string, objs ...interface{}) revel.Result {
	return c.renderErrorString(400, fmt.Sprintf(msg, objs))
}

// RenderError generates a 500 with the text of the error.
func (c ApiWrapper) renderError(e error) revel.Result {
	return c.renderErrorString(500, e.Error())
}

// RenderErrorString generates an error message the custom message and the supplied status code.
func (c ApiWrapper) renderErrorString(status int, e string) revel.Result {
	r := c.renderJson(nil, status, e)
	// Check if always 200
	var always200 bool
	c.Params.Bind(&always200, "always200")
	if always200 {
		return r
	}
	c.Response.Status = status
	return r
}

// renderErr generates an error message with the text of the error and the supplied status code.
func (c ApiWrapper) renderErr(status int, e error) revel.Result {
	return c.renderErrorString(status, e.Error())
}

// RenderJson renders the content of the interface, and wraps it is appropriate boilerplate.
// If an xml parameter is added, it will be xml instead.
func (c ApiWrapper) RenderJson(o interface{}) revel.Result {
	return c.renderJson(o, 200, "")
}

func (c ApiWrapper) renderJson(o interface{}, status int, e string) revel.Result {
	var xml bool
	c.Params.Bind(&xml, "xml")

	var j Reply
	c.Params.Bind(&j.Context, "context")
	j.Data = o
	j.Status = status
	j.Error = e

	if xml {
		return c.Controller.RenderXml(j)
	}
	// Check Callback
	var callback string = c.Params.Get("callback")
	if len(callback) > 0 {
		return c.Controller.RenderJsonP(callback, j)
	}
	return c.Controller.RenderJson(j)
}

//[CUT private functions]
