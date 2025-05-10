package runway

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/reports"
	"client-runaway-zenoti/internal/zenotiLegacy/zenoti"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	"fmt"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
)

func GetStats(stat string, from, to time.Time, locationIds, tags []string) (interface{}, error) {

	switch stat {
	case "Expenses":
		return reports.GetExpenses(from, to, locationIds)
	case "Sales":
		sales, err := GetOpportunitiesForLocations(from, to, locationIds, tags, "SalesId")
		return sales.Sum, err
	case "SalesNo":
		sales, err := GetOpportunitiesForLocations(from, to, locationIds, tags, "SalesId")
		return sales.Count, err
	case "ROI":
		sales, err := GetStats("Sales", from, to, locationIds, tags)
		if err != nil {
			return 0, err
		}
		expenses, err := GetStats("Expenses", from, to, locationIds, tags)
		return sales.(float64) - expenses.(float64), err
	case "NewLeads":
		leads, err := GetOpportunitiesForLocations(from, to, locationIds, tags, "")
		return leads.Count, err
	case "Bookings":
		bookings, err := GetOpportunitiesForLocations(from, to, locationIds, tags, "BookId")
		if err != nil {
			return 0, err
		}
		sales, err := GetStats("SalesNo", from, to, locationIds, tags)
		if err != nil {
			return 0, err
		}
		noShows, err := GetStats("NoShows", from, to, locationIds, tags)
		if err != nil {
			return 0, err
		}
		showNoSales, err := GetStats("ShowNoSale", from, to, locationIds, tags)
		return bookings.Count + sales.(int) + noShows.(int) + showNoSales.(int), err
	case "NoShows":
		data, err := GetOpportunitiesForLocations(from, to, locationIds, tags, "NoShowsId")
		return data.Count, err
	case "ShowNoSale":
		data, err := GetOpportunitiesForLocations(from, to, locationIds, tags, "ShowNoSaleId")
		return data.Count, err

	case "LeadsConv":
		leads, err := GetStats("NewLeads", from, to, locationIds, tags)
		if err != nil {
			return 0, err
		}

		bookings, err := GetStats("Bookings", from, to, locationIds, tags)

		if leads.(int) == 0 {
			return 0, err
		} else {
			return 100 * bookings.(int) / leads.(int), err
		}

	case "BookingsConv":
		sales, err := GetStats("SalesNo", from, to, locationIds, tags)
		if err != nil {
			return 0, err
		}

		bookings, err := GetStats("Bookings", from, to, locationIds, tags)

		if bookings.(int) == 0 {
			return 0, err
		} else {
			return 100 * sales.(int) / bookings.(int), err
		}
	case "ShowRate":
		noShows, err := GetStats("NoShows", from, to, locationIds, tags)
		if err != nil {
			return 0, err
		}

		bookings, err := GetStats("Bookings", from, to, locationIds, tags)

		if bookings.(int) == 0 {
			return 0, err
		} else {
			return 100 * (bookings.(int) - noShows.(int)) / bookings.(int), err
		}
	case "Rank":
		return GetRanking(from, to, locationIds)
	case "ZenotiMembersNo":
		return GetZentotiActiveMembers(locationIds)
	case "MembershipConv":
		members, err := GetOpportunitiesForLocations(from, to, locationIds, tags, "MemberId")
		if err != nil {
			return 0, err
		}

		bookings, err := GetStats("Bookings", from, to, locationIds, tags)

		if bookings.(int) == 0 {
			return 0, err
		} else {
			return 100 * members.Count / bookings.(int), err
		}
	default:
		return nil, fmt.Errorf("unknown stat (not implemented)")
	}
}

func GetLeads(from, to time.Time, locationIds, tags []string) (ReportsRes, error) {
	return GetOpportunitiesForLocations(from, to, locationIds, tags, "")
}

func SetRanksBySales() {
	fmt.Println("Starting ranking locations")
	type item struct {
		location models.Location
		sale     float64
	}

	to := time.Now()
	from := to.Add(-30 * 24 * time.Hour)

	locs := []models.Location{}
	sales := []item{}

	db.DB.Find(&locs)

	for _, l := range locs {
		sale, err := GetStats("Sales", from, to, []string{l.Id}, []string{})
		if err != nil {
			fmt.Println(err)
		}
		db.DB.Where("id", l.Id).First(&l)
		sales = append(sales, item{
			location: l,
			sale:     sale.(float64),
		})
	}

	for k, l := range sales {
		l.location.Rank = int64(k) + 1
		db.DB.Save(&l.location)
	}
	fmt.Println("Finished ranking locations")
}

func GetRanking(from, to time.Time, locationIds []string) (int64, error) {
	if len(locationIds) == 0 {
		return 0, nil
	}
	loc := models.Location{}

	loc.Get(locationIds[0])
	return loc.Rank, nil
}

func GetTags() []string {
	return []string{
		"Body Contouring | Lead",
		"Facial | Lead",
		"Botox | Lead",
		"Filler | Lead",
		"Laser | Lead",
		"Hair Restoration | Lead",
		"Semaglutide | Lead",
		"Membership | Lead",
		"Morpheus | Lead",
		"IV Therapy | Lead",
		"Weight loss | Lead",
		"Microneedling | Lead",
		"Graveyard",
		"Website | Lead",
		"Email | Lead",
		"Email List ",
		"Walk in | Lead",
	}
}

func getOpportunitiesWithTags(from, to time.Time, tags []string, locId, stageIdField string) (locationsRes, error) {
	res := locationsRes{}
	loc := models.Location{}
	loc.Get(locId)
	client, err := svc.NewClientFromId(locId)
	if err != nil {
		return res, err
	}

	res.Location = loc

	filter := runwayv2.OpportunitiesFilter{
		StartDate:  from,
		EndDate:    to,
		PipelineId: loc.PipelineId,
	}

	if stageIdField != "" {
		stage := lvn.GetValue[string](loc, stageIdField)
		if stage == "" {
			return res, nil
		}
		filter.StageId = lvn.GetValue[string](loc, stageIdField)
	}

	if len(tags) == 0 {
		tags = append(tags, "")
	}

	for _, tag := range tags {
		tagOpps := reportsOpp{}
		filter.Query = tag
		opps, err := client.OpportunitiesGetAll(filter)

		if err != nil {
			return res, err
		}

		tagOpps.Opportunities = opps
		for _, o := range opps {
			tagOpps.Sum += o.MonetaryValue
		}
		tagOpps.Count = len(opps)
		res.Opportunities = append(res.Opportunities, tagOpps)
		res.Count += tagOpps.Count
		res.Sum += tagOpps.Sum
	}

	return res, nil
}

func GetOpportunitiesForLocations(from, to time.Time, locationIds, tags []string, stageIdField string) (ReportsRes, error) {
	res := ReportsRes{}

	for _, locId := range locationIds {
		loc, err := getOpportunitiesWithTags(from, to, tags, locId, stageIdField)
		if err != nil {
			return res, err
		}
		res.Count += loc.Count
		res.Sum += loc.Sum
		res.Opportunities = append(res.Opportunities, loc)
	}

	return res, nil
}

func GetZentotiActiveMembers(locationIds []string) (int, error) {
	res := 0
	for _, locId := range locationIds {
		loc := config.GetLocationByRWID(locId)
		guests, err := zenoti.GetMembersNo(loc)
		if err != nil {
			return 0, err
		}
		res += guests
	}
	return res, nil
}
