package migrations

import (
	"fmt"

	"github.com/mmaelzer/opencam/store"
)

func Migrate(dbtype string) error {
	switch dbtype {
	case "sqlite3":
		return migrate(GetSQLiteMigrations())
	}
	return fmt.Errorf("Unsupported database type %s", dbtype)
}

func migrate(migrations []string) error {
	for i := range migrations {
		if _, err := store.RawWrite(migrations[i]); err != nil {
			return err
		}
	}
	return nil
}
