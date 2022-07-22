package converts

import (
	"github.com/google/uuid"
)

func ConvertUUIDToString(from uuid.UUID) string {
	return from.String()
}

func ConvertStringToUUID(from string) (uuid.UUID, error) {
	return uuid.Parse(from)
}
