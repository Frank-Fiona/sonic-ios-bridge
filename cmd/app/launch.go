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
package app

import (
	"github.com/Frank-Fiona/sonic-ios-bridge/src/util"
	"os"

	"github.com/spf13/cobra"
)

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch App",
	Long:  "Launch App",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		_, errLaunch := device.AppLaunch(bundleId)
		if errLaunch != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "launch", errLaunch)
		}
		return nil
	},
}

func initAppLaunch() {
	appRootCMD.AddCommand(launchCmd)
	launchCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	launchCmd.Flags().StringVarP(&bundleId, "bundleId", "b", "", "target bundleId")
	launchCmd.MarkFlagRequired("bundleId")
}
