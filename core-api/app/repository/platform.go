package repository

import (
	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
)

type PlatformRepo struct {
	db *database.DB
}

func NewPlatformRepo(db *database.DB) PlatformRepository {
	return &PlatformRepo{db}
}

func (r *PlatformRepo) GetAllPlatforms() ([]model.Platform, error) {
	query := `
        SELECT id, name, status, created_at, updated_at 
        FROM platforms 
        ORDER BY name ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var platforms []model.Platform
	for rows.Next() {
		var platform model.Platform
		err := rows.Scan(
			&platform.ID,
			&platform.Name,
			&platform.Status,
			&platform.CreatedAt,
			&platform.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		platforms = append(platforms, platform)
	}

	return platforms, nil
}
