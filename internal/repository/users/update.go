package users

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/noskov-sergey/auth/internal/model"
	"time"
)

func (r *UserRepository) Update(ctx context.Context, u model.User) error {
	builder := sq.
		Update(tableName).
		PlaceholderFormat(sq.Dollar).
		SetMap(map[string]any{
			nameColumn:      u.Name,
			emailColumn:     u.Email,
			updatedAtColumn: time.Now(),
		}).
		Where("id = ?", u.ID)

	sqlQuery, args, err := builder.ToSql()

	if _, err = r.db.Exec(ctx, sqlQuery, args...); err != nil {
		return fmt.Errorf("exec row: %w", err)
	}

	return nil
}
