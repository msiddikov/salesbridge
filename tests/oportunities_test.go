package tests

import (
	"client-runaway-zenoti/internal/runway"
	"fmt"
	"testing"
	"time"
)

func TestTags(t *testing.T) {
	sales, err := runway.GetStats("Sales",
		time.Now().Add(-30*24*time.Hour),
		time.Now(),
		[]string{"TknDfZc7WZF0cq5YE0du"},
		[]string{"Morpheus | Lead", "membership | lead"},
	)

	if err != nil {
		t.Error(err)
	}
	fmt.Println(sales)
}
