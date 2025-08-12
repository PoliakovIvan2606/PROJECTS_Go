package storage

import (
	"context"
	"internet_shop/internal/domain/product/model"
	"internet_shop/pkg/client/postgresql"
	db "internet_shop/pkg/client/postgresql/model"
	"internet_shop/pkg/logging"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ProductStorage struct {
	queryBuilder sq.StatementBuilderType
	client       *pgxpool.Pool
	logger       *logging.Logger
}

func NewProductStorage(client *pgxpool.Pool, logger *logging.Logger) ProductStorage {
	return ProductStorage{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		client: client,
		logger: logger,
	}
}

const table = "public.product"

func (s *ProductStorage) All(ctx context.Context) ([]model.Product, error) {
	query := s.queryBuilder.Select("id").
		Column("name").
		Column("description").
		Column("image_id").
		Column("price").
		Column("currency_id").
		Column("rating").
		Column("category_id").
		Column("specification").
		Column("created_at").
		Column("updated_at").
		From(table)
		
	sql, args, err := query.ToSql()
	logger := s.logger.WithFields(map[string]interface{}{
		"sql": sql,
		"table": table,
		"args": args,
	})

	if err != nil {
		err := db.ErrCreateQuery(err)
		logger.Error(err)
		return nil, err
	}

	logger.Trace("do query")
	rows, err := s.client.Query(ctx, sql, args...)
	if err != nil {
		err := db.ErrDoQuery(err)
		logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	list := make([]model.Product, 0)
	for rows.Next() {
		p := model.Product{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.ImageID,
			 &p.Price, &p.CurrencyID, &p.Rating, &p.CategoryID, &p.Specification, &p.CreatedAt, &p.UpdatedAt,
			); err != nil {
				err := db.ErrScan(postgresql.ParsePgError(err))
				logger.Error(err)
				return nil, err
		}

		list = append(list, p)
	}

	return list, nil
}