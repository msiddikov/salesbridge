package automator

import (
	"fmt"
	"time"
)

func getTotalPagesToProcess(totalItems, itemsPerPage, pageFrom, pageTo int) int {
	if itemsPerPage <= 0 {
		return 0
	}

	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage

	if pageFrom < 1 {
		pageFrom = 1
	}
	if pageTo > totalPages || pageTo <= 0 {
		pageTo = totalPages
	}

	if pageFrom > pageTo {
		return 0
	}

	return pageTo - pageFrom + 1
}

func getTotalItemsToProcess(totalItems, itemsPerPage, pageFrom, pageTo int) int {
	if itemsPerPage <= 0 {
		return 0
	}

	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage

	if pageFrom < 1 {
		pageFrom = 1
	}
	if pageTo > totalPages || pageTo <= 0 {
		pageTo = totalPages
	}

	if pageFrom > pageTo {
		return 0
	}

	pagesToProcess := pageTo - pageFrom + 1
	itemsToProcess := 0

	if pageTo == totalPages {
		// Adjust for the last page which might not be full
		lastPageItems := totalItems - (totalPages-1)*itemsPerPage
		pagesToProcess -= 1
		itemsToProcess = pagesToProcess*itemsPerPage + lastPageItems
	} else {
		itemsToProcess = pagesToProcess * itemsPerPage
	}

	if itemsToProcess > totalItems {
		return totalItems
	}

	return itemsToProcess
}

func parseTime(timeStr string) (time.Time, error) {
	layouts := []string{
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"01/02/2006 15:04:05",
		"01/02/2006 15:04",
		"01/02/2006",
		"01/02/06",
		"1/2/2006 3:04:05 PM",
		"1/2/2006 3:04 PM",
		"1/2/2006",
		"01/02/2006 03:04:05 PM",
		"01/02/2006 03:04 PM",
		"01/02/2006 03:04PM",
		"Monday, January 2, 2006 3:04 PM",
		"January 2, 2006 3:04 PM",
		"January 2, 2006",
		time.RFC3339,
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, timeStr)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", timeStr)
}
