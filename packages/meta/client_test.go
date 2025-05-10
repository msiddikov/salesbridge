package meta

import (
	"fmt"
	"testing"
)

func TestGetFields(t *testing.T) {
	fields := getFields[[]AdAccount]()
	fmt.Println(fields)
}
