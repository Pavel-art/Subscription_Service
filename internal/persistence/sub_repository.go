package persistence

import (
	"SubscriptionService/internal/core/core_interfaces"
	"SubscriptionService/internal/core/models"
	"SubscriptionService/internal/core/ports/filters"
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubRepository struct {
	db *pgxpool.Pool
}

var _ core_interfaces.ISubRepository = (*SubRepository)(nil)

func NewSubRepository(db *pgxpool.Pool) *SubRepository {
	return &SubRepository{db: db}
}

// Таблица
const tableName = "subscriptions"

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func (s *SubRepository) Create(ctx context.Context, sub *models.Subscription) (*models.Subscription, error) {
	query := psql.Insert(tableName).
		Columns("id", "service_name", "price", "user_id", "start_date", "end_date", "created_at", "updated_at").
		Values(sub.Id, sub.ServiceName, sub.Price, sub.UserId, sub.StartDate, sub.EndDate, sub.CreatedAt, sub.UpdatedAt).
		Suffix("RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build insert query: %w", err)
	}

	row := s.db.QueryRow(ctx, sqlStr, args...)
	var result models.Subscription
	err = row.Scan(&result.Id, &result.ServiceName, &result.Price, &result.UserId,
		&result.StartDate, &result.EndDate, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert subscription: %w", err)
	}

	return &result, nil
}

// Update --- UPDATE ---
func (s *SubRepository) Update(ctx context.Context, sub *models.Subscription) (*models.Subscription, error) {
	sub.UpdatedAt = time.Now()

	query := psql.Update(tableName).
		Set("service_name", sub.ServiceName).
		Set("price", sub.Price).
		Set("start_date", sub.StartDate).
		Set("end_date", sub.EndDate).
		Set("updated_at", sub.UpdatedAt).
		Where(squirrel.Eq{"id": sub.Id}).
		Suffix("RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build update query: %w", err)
	}

	row := s.db.QueryRow(ctx, sqlStr, args...)
	var result models.Subscription
	err = row.Scan(&result.Id, &result.ServiceName, &result.Price, &result.UserId,
		&result.StartDate, &result.EndDate, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("update subscription: %w", err)
	}

	return &result, nil
}

// Delete --- DELETE ---
func (s *SubRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := psql.Delete(tableName).Where(squirrel.Eq{"id": id})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build delete query: %w", err)
	}

	cmd, err := s.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// GetById --- GET BY ID ---
func (s *SubRepository) GetById(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	query := psql.Select("id", "service_name", "price", "user_id", "start_date", "end_date", "created_at", "updated_at").
		From(tableName).
		Where(squirrel.Eq{"id": id})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select by id query: %w", err)
	}

	row := s.db.QueryRow(ctx, sqlStr, args...)
	var sub models.Subscription
	err = row.Scan(&sub.Id, &sub.ServiceName, &sub.Price, &sub.UserId,
		&sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get by id: %w", err)
	}
	return &sub, nil
}

// GetAll --- GET ALL ---
func (s *SubRepository) GetAll(ctx context.Context, page, pageSize int64) ([]*models.Subscription, int64, int64, error) {
	offset := (page - 1) * pageSize

	// Общее количество
	countQuery := psql.Select("COUNT(*)").From(tableName)
	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, 0, fmt.Errorf("build count query: %w", err)
	}

	var totalCount int64
	err = s.db.QueryRow(ctx, countSQL, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("count query: %w", err)
	}

	totalPages := int64(math.Ceil(float64(totalCount) / float64(pageSize)))

	query := psql.Select("id", "service_name", "price", "user_id", "start_date", "end_date", "created_at", "updated_at").
		From(tableName).
		OrderBy("created_at DESC").
		Limit(uint64(pageSize)).
		Offset(uint64(offset))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, 0, 0, fmt.Errorf("build get all query: %w", err)
	}

	rows, err := s.db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("get all query: %w", err)
	}
	defer rows.Close()

	var subs []*models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err = rows.Scan(&sub.Id, &sub.ServiceName, &sub.Price, &sub.UserId,
			&sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("scan subscriptions: %w", err)
		}
		subs = append(subs, &sub)
	}

	return subs, totalCount, totalPages, nil
}

// SumSubscriptionsCost --- SUM (Filter) ---
func (s *SubRepository) SumSubscriptionsCost(ctx context.Context, filter *filters.SubFilter) (int64, error) {
	query := psql.Select("COALESCE(SUM(price), 0)").From(tableName)

	if filter.UserID != nil {
		query = query.Where(squirrel.Eq{"user_id": *filter.UserID})
	}
	if filter.ServiceName != nil && *filter.ServiceName != "" {
		query = query.Where(squirrel.Eq{"service_name": *filter.ServiceName})
	}
	if filter.From != nil && !filter.From.IsZero() {
		query = query.Where(squirrel.GtOrEq{"start_date": *filter.From})
	}
	if filter.To != nil && !filter.To.IsZero() {
		query = query.Where(squirrel.LtOrEq{"end_date": *filter.To})
	}

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build sum query: %w", err)
	}

	var total int64
	err = s.db.QueryRow(ctx, sqlStr, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("execute sum query: %w", err)
	}
	return total, nil
}
