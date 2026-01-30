package svc_cerbo

import (
	"client-runaway-zenoti/internal/db/models"
	"strconv"
)

func ListUsers(location models.Location) (map[string]uint, error) {
	cli, err := clientForLocation(location)
	if err != nil {
		return nil, err
	}

	users, err := cli.GetUsers()
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]uint)
	for _, user := range users {
		name := ""
		if user.DisplayNameForMessages != nil {
			name = *user.DisplayNameForMessages
		} else {
			name = user.FirstName + " " + user.LastName
		}
		id, err := strconv.Atoi(user.Id)
		if err != nil {
			return nil, err
		}
		userMap[name] = uint(id)
	}

	return userMap, nil
}
