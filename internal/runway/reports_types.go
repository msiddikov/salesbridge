package runway

import (
	"client-runaway-zenoti/internal/db/models"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
)

type (
	ReportsRes struct {
		Opportunities []locationsRes
		Count         int
		Sum           float64
	}

	locationsRes struct {
		Location      models.Location
		Opportunities []reportsOpp
		Count         int
		Sum           float64
	}

	reportsOpp struct {
		Tag           string
		Opportunities []runwayv2.Opportunity
		Count         int
		Sum           float64
	}
)

func (r *ReportsRes) Add(a ReportsRes) {
	for _, loca := range a.Opportunities {
		rl := r.getLocationData(loca.Location)

		for _, taga := range loca.Opportunities {
			tagData := rl.getTagData(taga.Tag)
			tagData.Opportunities = append(tagData.Opportunities, taga.Opportunities...)

			tagData.Count += taga.Count
			tagData.Sum += taga.Sum
		}

		rl.Count += loca.Count
		rl.Sum += loca.Sum
	}
}

func (r *ReportsRes) getLocationData(location models.Location) *locationsRes {
	for i, loc := range r.Opportunities {
		if loc.Location.Id == location.Id {
			return &r.Opportunities[i]
		}
	}

	r.Opportunities = append(r.Opportunities, locationsRes{
		Location: location,
	})
	return &r.Opportunities[len(r.Opportunities)-1]
}

func (lr *locationsRes) getTagData(tag string) *reportsOpp {
	for i, tagR := range lr.Opportunities {
		if tagR.Tag == tag {
			return &lr.Opportunities[i]
		}
	}

	lr.Opportunities = append(lr.Opportunities, reportsOpp{
		Tag: tag,
	})
	return &lr.Opportunities[len(lr.Opportunities)-1]
}
