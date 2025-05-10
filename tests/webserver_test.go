package tests

import (
	webServer "client-runaway-zenoti/internal/webserver"
	"testing"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
)

func TestWebServer(t *testing.T) {
	go webServer.Listen()

	lvn.WaitExitSignal()
}
