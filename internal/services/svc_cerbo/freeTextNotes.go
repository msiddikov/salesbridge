package svc_cerbo

import (
	"client-runaway-zenoti/internal/db/models"
	"fmt"
	"strings"
)

func UpsertSectionIntoFreeTextNote(location models.Location, patientID string, noteTypeID uint, sectionHeader, sectionContent string) error {

	cli, err := clientForLocation(location)
	if err != nil {
		return err
	}

	existingContent, err := cli.FtnGetFreeTextNote(patientID, noteTypeID)
	if err != nil {
		return err
	}

	newSection := fmt.Sprintf("---<b>%s</b>---:<br>%s<br>---END---<br>", sectionHeader, sectionContent)

	// Check if section already exists
	startMarker := fmt.Sprintf("---<b>%s</b>---:<br>", sectionHeader)
	start := strings.Index(existingContent, startMarker)
	if start != -1 {
		sectionStart := start
		if start >= 4 && existingContent[start-4:start] == "<br>" {
			sectionStart = start - 4
		}
		searchFrom := start + len(startMarker)
		if endRel := strings.Index(existingContent[searchFrom:], "---END---<br>"); endRel != -1 {
			end := searchFrom + endRel + len("---END---<br>")
			replacement := newSection
			if sectionStart == 0 {
				replacement = strings.TrimPrefix(newSection, "<br>")
			}
			updatedContent := existingContent[:sectionStart] + replacement + existingContent[end:]
			return cli.FtnUpdateFreeTextNote(patientID, noteTypeID, updatedContent)
		}
	}

	updatedContent := newSection + existingContent

	return cli.FtnUpdateFreeTextNote(patientID, noteTypeID, updatedContent)
}

func ListFreeTextNoteTypes(location models.Location) (map[string]uint, error) {
	cli, err := clientForLocation(location)
	if err != nil {
		return nil, err
	}

	noteTypes, err := cli.FtnGetAllAvailableTypes()
	if err != nil {
		return nil, err
	}

	noteTypeMap := make(map[string]uint)
	for _, nt := range noteTypes {
		noteTypeMap[nt.Name] = nt.Id
	}

	return noteTypeMap, nil
}
