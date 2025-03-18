package sqlite

import "fmt"

type DataCorruptedError struct {
	Type     string
	Original error
}

func (e DataCorruptedError) Error() string {
	return fmt.Sprintf(
		"DataCorruptedError[type=%s, original=%s]",
		e.Type, e.Original.Error(),
	)
}
