package runaway

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/types"
	"errors"
	"time"
)

func GetLeads(from, to time.Time, l []string) (int, error) {
	res := 0
	for _, locId := range l {
		loc := config.GetLocationById(locId)
		ppln, err := GetPipeline(loc)
		if err != nil {
			return 0, err
		}
		for _, stage := range ppln.Stages {
			count, err := getStageCount(loc, stage.Id, from, to)
			if err != nil {
				return 0, err
			}
			res += count

		}
	}
	return res, nil
}

func GetBookings(from, to time.Time, l []string) (int, error) {
	res := 0
	for _, locId := range l {
		loc := config.GetLocationById(locId)
		count, err := getStageCount(loc, loc.BookId, from, to)
		if err != nil {
			return 0, err
		}
		res += count
	}
	return res, nil
}

func GetNoShows(from, to time.Time, l []string) (int, error) {
	res := 0
	for _, locId := range l {
		loc := config.GetLocationById(locId)
		count, err := getStageCount(loc, loc.NoShowsId, from, to)
		if err != nil {
			return 0, err
		}
		res += count
	}
	return res, nil
}

func GetSales(from, to time.Time, l []string) (float64, error) {
	total := float64(0)
	meta := types.Meta{}
	fetchedAll := false
	for _, locId := range l {
		loc := config.GetLocationById(locId)
		for !fetchedAll {
			opps, meta1, err := getOpportunitiesByMetaPeriod(from, to, loc, meta, "open", loc.SalesId)
			if err != nil {
				return 0, err
			}
			meta = meta1
			if meta.StartAfter == 0 {
				fetchedAll = true
			}

			for _, o := range opps {
				total += o.MonetaryValue
			}
		}

	}
	return total, nil
}

func getStageCount(l config.Location, stageId string, from, to time.Time) (int, error) {
	if stageId == "" {
		return 0, errors.New("stageId is not filled")
	}
	count := 0
	meta := types.Meta{}
	fetchedAll := false

	for !fetchedAll {
		opps, meta1, err := getOpportunitiesByMetaPeriod(from, to, l, meta, "open", stageId)
		if err != nil {
			return 0, err
		}
		meta = meta1
		if meta.StartAfter == 0 {
			fetchedAll = true
		}
		count += len(opps)
	}
	return count, nil
}
