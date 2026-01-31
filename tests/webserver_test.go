package tests

import (
	"client-runaway-zenoti/internal/mcp"
	webServer "client-runaway-zenoti/internal/webserver"
	"testing"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
)

func TestWebServer(t *testing.T) {
	go webServer.Listen()

	lvn.WaitExitSignal()
}

func TestMcpServer(t *testing.T) {
	server := mcp.NewMCPServer(":8090")

	go server.Run(":8090")

	lvn.WaitExitSignal()
}
