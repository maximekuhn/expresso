package sqlite

import (
	"database/sql"
	"fmt"
)

func checkRowsAffected(res sql.Result, expected int) error {
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c != int64(expected) {
		return fmt.Errorf("expected to affect %d row(s), affected %d", expected, c)
	}
	return nil
}
