package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ndreyserg/ushort/internal/app/models"
)

type dbStorage struct {
	db *sql.DB
}

// Get получение оригинала ссылки по ее ID.
func (s *dbStorage) Get(ctx context.Context, key string) (string, error) {

	row := s.db.QueryRowContext(
		ctx,
		"select original, is_deleted from short_urls where short = $1",
		key,
	)

	if row.Err() != nil {
		return "", row.Err()
	}

	var original string
	var isDeleted bool
	err := row.Scan(&original, &isDeleted)

	if err != nil {
		return "", err
	}

	if isDeleted {
		return "", ErrIsGone
	}

	return original, nil
}

// Set - сохранение ссылки в хранилище и получение ee ID.
func (s *dbStorage) Set(ctx context.Context, val string, userID string) (string, error) {
	short := getUniqKey()

	row := s.db.QueryRowContext(
		ctx,
		`insert into short_urls (short, original, user_id) values ($1, $2, $3)
		on conflict (original) do UPDATE SET original = EXCLUDED.original returning short`,
		short,
		val,
		userID,
	)

	if row.Err() != nil {
		return "", row.Err()
	}

	var savedShort string

	err := row.Scan(&savedShort)

	if err != nil {
		return "", err
	}

	if savedShort != short {
		return savedShort, ErrConflict
	}

	return savedShort, nil
}

// Check проверка доступности БД.
func (s *dbStorage) Check(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// Close закрытие хранилища.
func (s *dbStorage) Close() error {
	return s.db.Close()
}

// SetBatch сохранение массива ссылок хранилище и получение их ID.
func (s *dbStorage) SetBatch(ctx context.Context, batch models.BatchRequest, userID string) (models.BatchResult, error) {
	result := make(models.BatchResult, 0, len(batch))
	tx, err := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err != nil {
		return nil, err
	}

	for _, item := range batch {
		resultItem := models.BatchResultItem{
			ID:    item.ID,
			Short: getUniqKey(),
		}
		_, err := tx.ExecContext(
			ctx,
			"insert into short_urls (short, original, user_id) values ($1, $2, $3)",
			resultItem.Short,
			item.Original,
			userID,
		)

		if err != nil {
			return nil, err
		}
		result = append(result, resultItem)
	}
	tx.Commit()
	return result, nil
}

// GetUserUrls получения сохраненных ссылок по пользователю.
func (s *dbStorage) GetUserUrls(ctx context.Context, userID string) ([]StorageItem, error) {

	res := []StorageItem{}
	query := "select short, original from short_urls where user_id = $1 and is_deleted = false"

	rows, err := s.db.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := StorageItem{}
		errN := rows.Scan(&item.Short, &item.Original)
		if errN != nil {
			return nil, errN
		}
		res = append(res, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteUserData удаление ссылок по пользователя.
func (s *dbStorage) DeleteUserData(ctx context.Context, ids []string, userID string) error {

	if len(ids) == 0 {
		return nil
	}

	query := "UPDATE short_urls SET is_deleted = true where user_id = $1 and short in ("
	ph := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids)+1)
	args = append(args, userID)
	for i := 0; i < len(ids); i++ {
		ph = append(ph, fmt.Sprintf("$%d", i+2))
		args = append(args, ids[i])
	}
	query = query + strings.Join(ph, ",") + ");"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func runMigrations(db *sql.DB) error {
	query := `
create table if not exists short_urls (
    short VARCHAR(20) NOT NULL PRIMARY KEY,
    original VARCHAR(1000) NOT NULL,
	user_id VARCHAR(100) NOT NULL,
	is_deleted boolean default false
)`
	_, err := db.ExecContext(context.Background(), query)

	if err != nil {
		return err
	}

	query = "create unique index if not exists short_urls_original_uniq on short_urls (original)"

	_, err = db.ExecContext(context.Background(), query)

	return err
}

// NewDBStorage создание БД хранилища
func NewDBStorage(dsn string) (Storage, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	if err = runMigrations(db); err != nil {
		return nil, err
	}

	s := &dbStorage{
		db: db,
	}
	return s, nil
}
