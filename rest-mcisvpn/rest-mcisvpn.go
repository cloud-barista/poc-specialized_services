// PoC-Specialized_Services based on a MCIS
// Create a MCIS-VPN environment on a MCIS consisting of multiple VMs
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// by B.T. Oh, innodreamer@gmail.com, 2019.11.

package restmcisvpn

import (
	"fmt"
	// "os"
	// "io/ioutil"
	// "strconv"
	"github.com/cloud-barista/poc-specialized_services/mcis-vpn-control"
	"github.com/cloud-barista/cb-store/config"

	"github.com/sirupsen/logrus"

	// REST API (echo)
	"net/http"
	"github.com/labstack/echo"
)

var cblog *logrus.Logger

func init() {
	cblog = config.Cblogger
}

// type vpncreationReq struct {
// 	Name string `json:"MCIS-VPN_Name"`
// }


//================ Credential Handler
func RestCreatMcisVpn(c echo.Context) error {
	// func RestCreatMcisVpn(c echo.Context) (mcisStatusInfo, int, []byte, error) {
	// func RestCreatMcisVpn(c echo.Context) (vpnInfo, int, []byte, error) {
	// func RestCreatMcisVpn(c echo.Context) error {

	nsId := c.Param("nsId")
	mcisId := c.Param("mcisId")	

	cblog.Info("call creatMcisVpn()")

	
	fmt.Println("[ Creat Mcis VPN!! ]")
	// req := &vpncreationReq{}
	// if err := c.Bind(req); err != nil {
	// 	return err
	// }
	// content, responseCode, body, err := mcisvpn.CreatMcisVpn(nsId, mcisId, req)
	content, responseCode, body, err := mcisvpn.CreatMcisVpn(nsId, mcisId)
	if err != nil {
		cblog.Error(err)
		fmt.Println("body: ", string(body))
		fmt.Println("responseCode: ", responseCode)
		// return c.JSONBlob(responseCode, body)
		return c.JSON(http.StatusCreated, content)
	}
	fmt.Println("Content!!:", content)
	fmt.Println("===========================")
	return c.JSON(http.StatusCreated, content)
}