package routers

import (
	"github.com/emicklei/go-restful"
)

// API Exposes the AccessLevel container for API's
// API: /v1/
func API() *restful.Container {
	container := restful.NewContainer()
	ws := WebService{}
	ws.Register(container, "/v1/AccessLevel")
	return container
}
