package tests

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/runway"
	"client-runaway-zenoti/internal/zenotiLegacy/integrations"
	"testing"
)

func Test_wrongColumns(t *testing.T) {

	l := config.GetLocationById("Solon")

	runway.ForceCheckAll()
	integrations.UpdateStatusesDbForLocation(l)
}
