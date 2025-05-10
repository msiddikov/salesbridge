package tgbot

import "testing"

func TestPolling(t *testing.T) {
	svc, err := NewTestService()
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	svc.StartPolling()
}

func TestSendNotify(t *testing.T) {
	svc, err := NewTestService()
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	svc.Notify("test", "Hello", true)
}
