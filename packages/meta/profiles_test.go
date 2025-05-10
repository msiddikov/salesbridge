package meta

import (
	"fmt"
	"testing"
)

func Test_Myself(t *testing.T) {
	cli := getTestClient()
	res, err := cli.Me()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}

func Test_Adaccounts(t *testing.T) {
	cli := getTestClient()
	res, err := cli.MyAdAccounts()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}
