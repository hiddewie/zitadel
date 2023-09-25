package setup

import (
	"context"
	"embed"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	//go:embed 14/cockroach/*.sql
	//go:embed 14/postgres/*.sql
	changeEvents embed.FS
)

type ChangeEvents struct {
	dbClient *database.DB
}

func (mig *ChangeEvents) Execute(ctx context.Context) error {
	migrations, err := changeEvents.ReadDir("14/" + mig.dbClient.Type())
	if err != nil {
		return err
	}
	for _, migration := range migrations {
		stmt, err := readStmt(changeEvents, "14", mig.dbClient.Type(), migration.Name())
		if err != nil {
			return err
		}

		logging.WithFields("migration", mig.String(), "file", migration.Name()).Debug("execute statement")

		_, err = mig.dbClient.ExecContext(ctx, stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mig *ChangeEvents) String() string {
	return "14_events_push"
}