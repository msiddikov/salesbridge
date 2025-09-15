package meta

import (
	"context"
	"testing"
)

func TestExchangeToken(t *testing.T) {
	short := "EAAK2EsTlzmEBPU4GqUPSTPTIgMVAHrJqFOWfe3ZCw6KLvhpOrPilfpvZASYWLnrVFH40WS3YZBtlJyZA3qCMLh71Iff3btjBXzActeqbe7T6s27ZCnYsg7SnQZCH7P9pTJyZBXM8QI80y1iMYxzemNL8wGAFIBrjA7XZCeZBHbjOA4E7UyXf3XG9VEpNyZCLPM4AW7H76QdKSE4kxqXhAJP3qf6kPpR6aqtk2yO7A0QVayfXvDZCAZDZD"
	long, err := ExchangeForLongLivedToken(context.Background(), "763141682482785", "4a2a73aa201509af22b12621b3d8741f", short)
	if err != nil {
		t.Fatalf("failed to exchange token: %v", err)
	}
	t.Logf("Long lived token: %v", long)

}
