package zenotiv1

import (
	"fmt"
	"testing"
)

type ()

func TestCentersList(t *testing.T) {
	client := getClient()

	res, err := CentersListAll(client.cfg.apiKey)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}

func TestListServices(t *testing.T) {
	client := getTrainingClient()
	services, err := client.CenterServicesGetAll(CenterServicesFilter{})
	if err != nil {
		t.Error(err)
	}
	if len(services) == 0 {
		t.Error("No services found")
	}
	fmt.Println(services)
}

func TestListRooms(t *testing.T) {
	client := getFairlawnClient()
	rooms, err := client.CenterRoomsGet()
	if err != nil {
		t.Error(err)
	}
	if len(rooms) == 0 {
		t.Error("No rooms found")
	}
}
