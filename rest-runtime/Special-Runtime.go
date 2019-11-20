// Proof of Concepts for the Cloud-Barista Multi-Cloud Project.
//      * Cloud-Barista: https://github.com/cloud-barista

package main

import (
	"fmt"
	// "github.com/cloud-barista/cb-store/config"
	"github.com/cloud-barista/cb-tumblebug/src/common"
	"github.com/cloud-barista/poc-specialized_services/rest-mcisvpn"
	"github.com/cloud-barista/poc-specialized_services/rest-mcislb"
	// "github.com/sirupsen/logrus"

	// REST API (echo)
	"net/http"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	fmt.Println("\n[ Cloud-Barista Specialized Services!! ]")
	fmt.Println("\nInitiating REST API Server ...")

	// Run API Server
	ApiServer()
}

func ApiServer() {

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, This is Cloud-Barista Specialized Service!!")
	})
	e.HideBanner = true

	// Route
	g := e.Group("/ns", common.NsValidation())

	// {{ip}}:{{port}}/ns/{{ns_id}}/mcis-special/mcis/{{mcis_id}}/mcis-vpn
	g.POST("/:nsId/mcis-special/mcis/:mcisId/mcis-vpn", restmcisvpn.RestCreatMcisVpn)
	// g.GET("/:nsId/mcis-special/mcis/:mcisId/mcis-vpn", restmcisvpn.RestGetMcisVpnStatus)
	// g.POST("/:nsId/mcis-special/mcis/:mcisId/mcis-lb", restmcislb.RestCreatMcisLb)
	e.Logger.Fatal(e.Start(":1313"))

}