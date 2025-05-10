package zenotiv1

import (
	"fmt"
	"testing"
	"time"
)

func TestReports(t *testing.T) {
	client := getTribecaClient()

	report, err := client.ReportsCollections(time.Now().AddDate(0, 0, -5), time.Now())

	if err != nil {
		t.Error(err)
	}

	fmt.Println(len(report))
}
