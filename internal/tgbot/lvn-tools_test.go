package tgbot

import (
	"fmt"
	"testing"
)

func TestSave(t *testing.T) {
	err := saveConfig("12345", "Hello3")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}

func TestGet(t *testing.T) {
	cfg, err := getConfig("12345", "Hello")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	fmt.Println(cfg)
}
