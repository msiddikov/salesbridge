package zenotiv1

import (
	"fmt"
	"testing"
)

func TestInvoicesGetDetails(t *testing.T) {
	client := getNaplesClient()

	res, err := client.InvoicesGetDetails("4829909a-b565-4ffc-9aeb-65546c456ab2")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}
