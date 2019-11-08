package routers

import (
	"net/http"

	"wwwin-github.cisco.com/DevNet/restful/config"

	"github.com/emicklei/go-restful"
)

type WebService struct {
}

// Register the API
// prefix: /v1/querys
func (w *WebService) Register(container *restful.Container, prefix string) {
	ws := &restful.WebService{}
	ws.
		Path(prefix).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		ApiVersion(config.String("version"))

	ws.Route(ws.GET("/").To(w.get))
	container.Add(ws)
}

// list all indexes
func (w *WebService) get(req *restful.Request, res *restful.Response) {
	res.WriteHeaderAndEntity(http.StatusOK, req.Request.Header)
}
