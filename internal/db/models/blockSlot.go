package models

import (
	runwayv2 "client-runaway-zenoti/packages/runwayV2"

	"gorm.io/gorm"
)

func (b *BlockSlot) BeforeCreate(tx *gorm.DB) (err error) {

	client, err := Svc.NewClientFromId(b.LocationId)
	if err != nil {
		return err
	}

	slot, err := client.CalendarCreateBlockSlot(runwayv2.BlockSlot{
		LocationId:    b.LocationId,
		CalendarId:    b.CalendarId,
		StartTime:     b.StartTime,
		EndTime:       b.EndTime,
		Title:         b.Title,
		CalendarNotes: b.Notes,
	})

	b.Id = slot.Id
	return err
}

func (b *BlockSlot) BeforeUpdate(tx *gorm.DB) (err error) {

	client, err := Svc.NewClientFromId(b.LocationId)
	if err != nil {
		return err
	}

	_, err = client.CalendarEditBlockSlot(b.Id, runwayv2.BlockSlot{
		LocationId:    b.LocationId,
		CalendarId:    b.CalendarId,
		StartTime:     b.StartTime,
		EndTime:       b.EndTime,
		Title:         b.Title,
		CalendarNotes: b.Notes,
	})

	return err
}

func (b *BlockSlot) BeforeDelete(tx *gorm.DB) (err error) {

	client, err := Svc.NewClientFromId(b.LocationId)
	if err != nil {
		return err
	}

	err = client.CalendarDeleteEvent(b.Id)

	return err
}

func BlockSlotHasDuplicates(blocks []BlockSlot) (hasDuplicates bool, res []BlockSlot) {

	for _, bs := range blocks {
		for _, bs2 := range blocks {
			if bs.ZenotiId == bs2.ZenotiId && bs.Id != bs2.Id {
				hasDuplicates = true
				res = append(res, bs)
			}
		}
	}
	return
}
