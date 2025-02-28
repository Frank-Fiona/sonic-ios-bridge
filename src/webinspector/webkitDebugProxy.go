/*
 *   sonic-ios-bridge  Connect to your iOS Devices.
 *   Copyright (C) 2022 SonicCloudOrg
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU Affero General Public License as published
 *   by the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU Affero General Public License for more details.
 *
 *   You should have received a copy of the GNU Affero General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package webinspector

import (
	"context"
	"fmt"
	"log"
	"net/http"

	giDevice "github.com/Frank-Fiona/sonic-gidevice"
	"github.com/Frank-Fiona/sonic-ios-bridge/src/util"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var webDebug *WebkitDebugService
var localPort = 9222
var isAdapter = false

func SetIsAdapter(flag bool) {
	isAdapter = flag
}

func InitWebInspectorServer(udid string, port int, isProtocolDebug bool, isDTXDebug bool) context.CancelFunc {
	var err error
	var cannel context.CancelFunc
	if webDebug == nil {
		// optimize the initialization process
		ctx := context.Background()
		device := util.GetDeviceByUdId(udid)
		webDebug = NewWebkitDebugService(&device, ctx, util.GetDeviceVersion(device))
		cannel, err = webDebug.ConnectInspector()
		if err != nil {
			log.Fatal(err)
		}
	}
	localPort = port
	if isProtocolDebug {
		SetProtocolDebug(true)
	}
	if isDTXDebug {
		giDevice.SetDebug(true, true)
	}
	return cannel
}

func PagesHandle(c *gin.Context) {

	pages, err := webDebug.GetOpenPages(localPort)
	if err != nil {
		c.JSONP(http.StatusNotExtended, err)
	}
	c.JSONP(http.StatusOK, pages)
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// solve cross domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func PageDebugHandle(c *gin.Context) {
	id := c.Param("id")

	application, page, err := webDebug.FindPagesByID(id)
	if application == nil || page == nil {
		c.Error(fmt.Errorf(fmt.Sprintf("not find page to id:%s", id)))
		log.Println(fmt.Errorf(fmt.Sprintf("not find page to id:%s", id)))
		return
	}
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	err = webDebug.StartCDP(application.ApplicationID, page.PageID, conn, isAdapter)
	if err != nil {
		log.Fatal(err)
	}

	// make sure initialization is complete
	err = webDebug.receiveWebkitProtocolData()
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		for {
			err = webDebug.receiveWebkitProtocolData()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	for {
		err = webDebug.receiveMessageTool()
		if err != nil {
			log.Panic(err)
		}
		if err == nil || err.Error() == "message is null" {
			continue
		}
	}

}
