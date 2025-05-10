package zenotiv1

import (
	"fmt"
	"testing"
)

func TestInvoicesGetDetails(t *testing.T) {
	client := getTribecaClient()

	res, err := client.InvoicesGetDetails("d663c848-2317-4ac9-bee4-5d6746d54506")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}
