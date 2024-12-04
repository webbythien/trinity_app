package repository

import (
	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
)

type EntityRepo struct {
	db *database.DB
}

func NewEntityRepo(db *database.DB) EntityRepository {
	return &EntityRepo{db}
}

func (r *EntityRepo) GetAllEntityTypes() ([]model.EntityType, error) {
	query := `
        SELECT id, type_name, table_name, status
        FROM entity_types 
        WHERE status = true
        ORDER BY type_name ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []model.EntityType
	for rows.Next() {
		var entity model.EntityType
		err := rows.Scan(
			&entity.ID,
			&entity.TypeName,
			&entity.TableName,
			&entity.Status,
		)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}

	return entities, nil
}
