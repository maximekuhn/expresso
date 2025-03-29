package common

import "time"

type DatetimeProvider struct{}

// Provide provides current date (now)
func (_ *DatetimeProvider) Provide() time.Time {
	return time.Now().UTC()
}
