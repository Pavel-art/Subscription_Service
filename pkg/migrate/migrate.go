package migrate

import (
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
)

// Up выполняет миграции вверх (создаёт/обновляет таблицы).
// dbURL — строка подключения postgres (из DB_URL).
func Up(dbURL string, logger *zerolog.Logger) error {
	if dbURL == "" {
		return fmt.Errorf("DB_URL is empty")
	}

	absPath, err := filepath.Abs("migrations")
	if err != nil {
		return fmt.Errorf("resolve migrations path: %w", err)
	}

	sourceURL := "file://" + absPath

	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}

	if err == migrate.ErrNoChange {
		logger.Debug().Msg("migrations: no new migrations to apply")
	} else {
		logger.Info().Msg("migrations: applied successfully")
	}

	return nil
}
