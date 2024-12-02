package repository

import (
	"fmt"

	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
)

type RegionRepo struct {
	db *database.DB
}

func NewRegionRepo(db *database.DB) RegionRepository {
	return &RegionRepo{db}
}

func (r *RegionRepo) GetAllProvince() ([]*model.ProvineUnitName, error) {
	var provinces []*model.ProvineUnitName
	query := `SELECT p.code, p."name" , p.full_name , p.full_name_en ,au.full_name as administrative_unit_name
FROM provinces p
INNER JOIN administrative_units au 
ON p.administrative_unit_id = au.id `

	if err := r.db.Select(&provinces, query); err != nil {
		return nil, fmt.Errorf("failed to get all provinces: %w", err)
	}
	return provinces, nil
}

func (r *RegionRepo) GetAllDistrict(provinceCode string) ([]*model.DistrictProvince, error) {
	var districts []*model.DistrictProvince
	query := `SELECT d.code, d."name" , d.full_name , d.full_name_en ,au.full_name as administrative_unit_name
FROM districts d 
INNER JOIN administrative_units au 
ON d.administrative_unit_id = au.id
WHERE d.province_code = $1
ORDER BY d.code;`

	if err := r.db.Select(&districts, query, provinceCode); err != nil {
		return nil, fmt.Errorf("failed to get all districts for province %s: %w", provinceCode, err)
	}
	return districts, nil
}

func (r *RegionRepo) GetAllWard(districtCode string) ([]*model.WardDistrictProvince, error) {
	var wards []*model.WardDistrictProvince
	query := `SELECT w.code, w."name" , w.full_name , w.full_name_en ,au.full_name as administrative_unit_name
FROM wards w 
INNER JOIN administrative_units au 
ON w.administrative_unit_id = au.id
WHERE w.district_code = $1
ORDER BY w.code;`

	if err := r.db.Select(&wards, query, districtCode); err != nil {
		return nil, fmt.Errorf("failed to get all wards for district %s: %w", districtCode, err)
	}
	return wards, nil
}
