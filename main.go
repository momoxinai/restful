package AccessLevel

import (
	"fmt"
	"net/http"

	"wwwin-github.cisco.com/DevNet/AccessLevel/log"

	"wwwin-github.cisco.com/DevNet/AccessLevel/config"

	"wwwin-github.cisco.com/DevNet/AccessLevel/routers"
)

func main() {
	log.NewLogger(config.String("appname"))
	hostPort := fmt.Sprintf(":%s", config.String("httpport"))
	log.LogInfof("starting %s (%s) webserver on %s",
		config.String("appname"), config.String("version"), hostPort)
	http.ListenAndServe(hostPort, routers.API())
}
