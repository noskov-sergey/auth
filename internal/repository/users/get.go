package users

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/noskov-sergey/auth/internal/model"
)

func (r *UserRepository) Get(ctx context.Context, filter model.UserFilter) (*model.User, error) {
	builder := sq.Select(
		idColumn,
		nameColumn,
		emailColumn,
		roleColumn,
		createdAtColumn,
		updatedAtColumn,
	).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: filter.UserID})

	sqlQuery, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("to sql: %w", err)
	}

	var user model.User

	err = r.db.QueryRow(ctx, sqlQuery, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.NotFoundErr
	} else if err != nil {
		return nil, fmt.Errorf("get query: %w", err)
	}

	return &user, nil
}
