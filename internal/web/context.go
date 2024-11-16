package web

import (
	"encoding/json"
	"net/http"

	"github.com/lorenzogood/blog/internal/templates"
)

type Ctx struct {
	r          *http.Request
	w          http.ResponseWriter
	statusCode HttpStatusCode
}

func (c *Ctx) Request() *http.Request {
	return c.r
}

func (c *Ctx) Response() http.ResponseWriter {
	return c.w
}

func (c *Ctx) SetStatus(s HttpStatusCode) {
	c.statusCode = s

	c.Response().WriteHeader(int(s))
}

func (c *Ctx) Header() http.Header {
	return c.w.Header()
}

func (c *Ctx) StatusCode() HttpStatusCode {
	return c.statusCode
}

func (c *Ctx) Respond(status HttpStatusCode, b []byte) error {
	c.SetStatus(status)
	_, err := c.Response().Write(b)
	return err
}

func (c *Ctx) RespondString(status HttpStatusCode, b string) error {
	c.Header().Set("Content-Type", "text/plain")
	return c.Respond(status, []byte(b))
}

func (c *Ctx) RespondTemplate(t *templates.Renderer, status HttpStatusCode, name string, data any) error {
	c.SetStatus(status)
	c.Header().Set("Content-Type", "text/html")

	return t.Render(c.w, name, data)
}

func (c *Ctx) RespondJson(status HttpStatusCode, data any) error {
	c.SetStatus(status)
	c.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(c.w).Encode(data); err != nil {
		return err
	}

	return nil
}

func (c *Ctx) Method() RouteMethod {
	return RouteMethod(c.r.Method)
}

func (c *Ctx) PathValue(p string) string {
	return c.r.PathValue(p)
}
