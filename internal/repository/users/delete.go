package users

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	builder := sq.
		Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where("id = ?", id)

	sqlQuery, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("to sql: %v", err)
	}

	if _, err = r.db.Exec(ctx, sqlQuery, args...); err != nil {
		return fmt.Errorf("exec row: %w", err)
	}

	return nil
}
