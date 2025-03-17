package common

import "time"

type DatetimeProvider struct{}

func (_ *DatetimeProvider) Provide() time.Time {
	return time.Now().UTC()
}
