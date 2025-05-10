package webServer

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

type (
	blockSlotsMsg struct {
		WorkingHoursFrom string
		WorkingHoursTo   string
		BlockSlots       []BlockSlots
	}

	BlockSlots struct {
		From time.Time
		To   time.Time
	}

	bookMsg struct {
		From           time.Time
		To             time.Time
		Customer_name  string
		Customer_phone string
	}
)

var (
	testLocId        = "VNDSQQmMekLitfnlMBqH"
	testCalendarId   = "N3rTXsSEbXzKfltPChch"
	workingHoursFrom = "09:00:00 -0700"
	workingHoursTo   = "18:00:00 -0700"
)

func setChatlyRoutes(router *gin.Engine) {
	router.GET("/chatly/block-slots/:date", getBlockSlots)
	router.POST("/chatly/book", book)
}

func getBlockSlots(c *gin.Context) {

	fromString := c.Param("date")

	from, err := time.Parse("2006-01-02", fromString)
	lvn.GinErr(c, 400, err, "Invalid date")

	blocks, err := getBlockSlotsForPeriod(from, from.AddDate(0, 0, 1))
	lvn.GinErr(c, 500, err, "Error getting block slots")

	res := blockSlotsMsg{
		WorkingHoursFrom: workingHoursFrom,
		WorkingHoursTo:   workingHoursTo,
		BlockSlots:       blocks,
	}

	c.Data(lvn.Res(200, res, "Success"))
}

func getBlockSlotsForPeriod(from, to time.Time) ([]BlockSlots, error) {
	blocks := []models.BlockSlot{}
	db.DB.Where("calendar_id = ? AND (start_time>? and start_time < ? OR end_time > ? and end_time < ?)", testCalendarId, from, to, from, to).Order("start_time asc").Find(&blocks)

	res := []BlockSlots{}
	for _, b := range blocks {
		res = append(res, BlockSlots{
			From: b.StartTime,
			To:   b.EndTime,
		})
	}
	return res, nil
}

func book(c *gin.Context) {

	body := bookMsg{}
	err := c.BindJSON(&body)
	lvn.GinErr(c, 400, err, "Invalid body")

	bs := models.BlockSlot{
		LocationId: testLocId,
		CalendarId: testCalendarId,
		StartTime:  body.From,
		EndTime:    body.To,
		Title:      body.Customer_name,
		Notes:      body.Customer_phone,
	}

	err = db.DB.Save(&bs).Error
	lvn.GinErr(c, 500, err, "Error saving block slot")

	c.Data(lvn.Res(200, "BlockSlot for user has been created", "OK"))

}
