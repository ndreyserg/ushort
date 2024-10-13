package storage

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ndreyserg/ushort/internal/app/models"
)

type dbStorage struct {
	db *sql.DB
}

func (s *dbStorage) Get(ctx context.Context, key string) (string, error) {

	row := s.db.QueryRowContext(
		ctx,
		"select original from short_urls where short = $1",
		key,
	)

	if row.Err() != nil {
		return "", row.Err()
	}

	var original string
	err := row.Scan(&original)

	if err != nil {
		return "", err
	}

	return original, nil
}

func (s *dbStorage) Set(ctx context.Context, val string) (string, error) {
	key := getUniqKey()
	_, err := s.db.ExecContext(
		ctx,
		"insert into short_urls (short, original) values ($1, $2)",
		key,
		val,
	)

	if err != nil {
		return "", err
	}

	return key, nil
}

func (s *dbStorage) Check(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *dbStorage) Close() error {
	return s.db.Close()
}

func (s *dbStorage) SetBatch(ctx context.Context, batch models.BatchRequest) (models.BatchResult, error) {
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
			"insert into short_urls (short, original) values ($1, $2)",
			resultItem.Short,
			item.Original,
		)

		if err != nil {
			return nil, err
		}
		result = append(result, resultItem)
	}
	tx.Commit()
	return result, nil
}

func createTable(db *sql.DB) error {
	query := `
CREATE TABLE IF NOT EXISTS short_urls (
    short VARCHAR(20) NOT NULL PRIMARY KEY,
    original VARCHAR(1000) NOT NULL
)`
	_, err := db.ExecContext(context.Background(), query)
	return err
}

func NewDBStorage(dsn string) (Storage, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	if err = createTable(db); err != nil {
		return nil, err
	}

	s := &dbStorage{
		db: db,
	}
	return s, nil
}
