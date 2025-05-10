package zenotiv1

import (
	"fmt"
	"testing"
	"time"
)

func TestBlockOuts(t *testing.T) {
	client := getClient()

	start, _ := time.Parse(time.RFC3339, "2023-06-14T00:00:00Z")
	end, _ := time.Parse(time.RFC3339, "2023-06-15T00:00:00Z")

	res, err := client.EmployeesListAllBlockOutTimes(start, end)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}
