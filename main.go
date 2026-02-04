package main

import (
	webServer "client-runaway-zenoti/internal/webserver"
	"client-runaway-zenoti/packages/grafana"

	integrations_zenoti "client-runaway-zenoti/internal/integrations/zenoti"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
)

func main() {
	grafana.Init()

	go webServer.Listen()

	//runaway.UpdateLocations()

	//go integrationsOld.StartScheduledJobs()
	go integrations_zenoti.StartScheduledJobs()

	lvn.WaitExitSignal()
}
