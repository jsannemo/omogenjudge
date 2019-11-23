// Package users contains utilities for OS-level user and group management.
package users

import (
	"os/user"
	"strconv"

	"github.com/google/logger"
)

// OmogenClientsID returns the group ID of the omogenjudge-clients group.
func OmogenClientsID() int {
	group, err := user.LookupGroup("omogenjudge-clients")
	if err != nil {
		logger.Fatalf("could not look up omogenjudge-clients group: %v", err)
	}
	id, err := strconv.Atoi(group.Gid)
	if err != nil {
		logger.Fatalf("could not convert gid to int: %v", err)
	}
	return id
}
