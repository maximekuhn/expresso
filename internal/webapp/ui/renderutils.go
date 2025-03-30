package ui

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func FormatUserNameAndId(name string, id uuid.UUID) string {
	return fmt.Sprintf("%s#%s", name, id.String()[:4])
}

func FormatGroupCreatedAt(d time.Time) string {
	return d.UTC().Format("2 January 2006, 15h04 (UTC)")
}
