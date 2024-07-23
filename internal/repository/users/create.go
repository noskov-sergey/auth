package users

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/noskov-sergey/auth/internal/model"
)

func (r *UserRepository) Create(ctx context.Context, u model.User) (int, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, roleColumn, passwordColumn, confirmPasswordColumn).
		Values(u.Name, u.Email, u.Role, u.Password, u.ConfirmPassword).
		Suffix("RETURNING id")

	sqlQuery, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("to sql: %w", err)
	}

	var insertedID int

	if err = r.db.QueryRow(ctx, sqlQuery, args...).Scan(&insertedID); err != nil {
		return 0, fmt.Errorf("query row: %w", err)
	}

	return insertedID, nil
}
