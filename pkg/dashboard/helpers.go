package dashboard

import (
	uuid "github.com/satori/go.uuid"
)

func getUUID() string {
	return uuid.NewV4().String()
}
