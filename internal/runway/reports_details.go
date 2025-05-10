package runway

import (
	"client-runaway-zenoti/internal/db/models"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
	"golang.org/x/exp/slices"
)

type (
	exportData struct {
		Location    models.Location
		Opportunity runwayv2.Opportunity
		StageName   string
	}

	sheetsSet struct {
		Sheet   string
		Display string
		Stages  []string
	}
)

func GetDetails(input ReportsRes) ([]byte, error) {
	////////////////////////////////////////////////////////////
	//
	//					Setup
	//
	///////////////////////////////////////////////////////////

	sheets := []sheetsSet{
		{
			Sheet:   "all",
			Display: "ALL LEADS",
			Stages: []string{
				"New Leads / Other",
				"Services Sold",
				"Booking Confirmed",
				"No Showed",
				"Showed But Didn't Purchase",
			},
		},
		{
			Sheet:   "sales",
			Display: "SERVICES SOLD",
			Stages: []string{
				"Services Sold",
			},
		},
		{
			Sheet:   "bookings",
			Display: "BOOKINGS",
			Stages: []string{
				"Services Sold",
				"Booking Confirmed",
				"No Showed",
				"Showed But Didn't Purchase",
			},
		},
		{
			Sheet:   "noshows",
			Display: "NO SHOWED",
			Stages: []string{
				"No Showed",
			},
		},
		{
			Sheet:   "showNoSales",
			Display: "SHOWED BUT DIDN'T PURCHASE",
			Stages: []string{
				"Showed But Didn't Purchase",
			},
		},
	}

	f := excelize.NewFile()
	for _, s := range sheets {
		f.NewSheet(s.Sheet)
	}
	f.DeleteSheet("Sheet1")

	///////////////////////////////////////////////////////////
	//
	//						Format data
	//
	//////////////////////////////////////////////////////////

	dataArray := []exportData{}

	for _, loc := range input.Opportunities {
		for _, tag := range loc.Opportunities {
			for _, opp := range tag.Opportunities {
				dataArray = append(dataArray, exportData{
					Location:    loc.Location,
					Opportunity: opp,
					StageName:   getStageName(loc.Location, opp),
				})
			}
		}
	}

	for _, sheet := range sheets {
		fillOpportunities(dataArray, f, sheet)
	}

	buffer, _ := f.WriteToBuffer()
	return buffer.Bytes(), nil
}

func fillOpportunities(dataArray []exportData, file *excelize.File, sheet sheetsSet) error {
	data := [][]interface{}{}
	data = append(data, []interface{}{
		sheet.Display,
	})
	data = append(data, []interface{}{
		"Location",
		"Name",
		"Monetary value",
		"Stage",
		"Tags",
		"Created at",
		"link",
	})

	for _, d := range dataArray {
		if !slices.Contains(sheet.Stages, d.StageName) {
			continue
		}

		data = append(data, []interface{}{
			d.Location.Name,
			d.Opportunity.Name,
			d.Opportunity.MonetaryValue,
			d.StageName,
			strings.Join(d.Opportunity.Contact.Tags, ", "),
			d.Opportunity.CreatedAt,
			fmt.Sprintf("https://app.clientrunway.com/v2/location/%s/contacts/detail/%s", d.Location.Id, d.Opportunity.Contact.Id),
		})
	}

	return appendToXls(file, sheet.Sheet, "A", 1, data)
}

func appendToXls(file *excelize.File, sheet, col string, row int, data [][]interface{}) error {
	var addr string
	var err error

	for r, dataRow := range data {
		if addr, err = excelize.JoinCellName(col, r+row); err != nil {
			fmt.Println(err)
			return err
		}
		if err = file.SetSheetRow(sheet, addr, &dataRow); err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func getStageName(loc models.Location, opp runwayv2.Opportunity) string {
	switch opp.PipelineStageId {
	case loc.SalesId:
		return "Services Sold"
	case loc.BookId:
		return "Booking Confirmed"
	case loc.NoShowsId:
		return "No Showed"
	case loc.ShowNoSaleId:
		return "Showed But Didn't Purchase"
	default:
		return "New Leads / Other"
	}
}
