package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

func NewPGXPool(ctx context.Context, connString string, logger *zerolog.Logger) (*pgxpool.Pool, error) {

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		logger.Error().Err(err).Msg("Ошибка парсинга строки подключения")
		return nil, err
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	logger.Debug().
		Int("max_conns", int(config.MaxConns)).
		Int("min_conns", int(config.MinConns)).
		Msg("Настройки пула соединений применены")

	if config.ConnConfig != nil {
		logger.Debug().
			Str("host", config.ConnConfig.Host).
			Uint16("port", config.ConnConfig.Port).
			Str("database", config.ConnConfig.Database).
			Str("user", config.ConnConfig.User).
			Msg("Параметры подключения к БД")
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Error().Err(err).Msg("Ошибка создания пула соединений")
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		logger.Error().Err(err).Msg("Проверка подключения не пройдена")
		return nil, err
	}

	logger.Info().
		Str("db_host", config.ConnConfig.Host).
		Str("db_name", config.ConnConfig.Database).
		Msg("Пул соединений к PostgreSQL успешно создан")

	return pool, nil
}
