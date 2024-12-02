package model

type ProvineUnitName struct {
	Code       string `db:"code" json:"code"`
	Name       string `db:"name" json:"name"`
	FullName   string `db:"full_name" json:"full_name"`
	FullNameEn string `db:"full_name_en" json:"full_name_en"`
	UnitName   string `db:"administrative_unit_name" json:"unit_name"`
}

type DistrictProvince struct {
	Code       string `db:"code" json:"code"`
	Name       string `db:"name" json:"name"`
	FullName   string `db:"full_name" json:"full_name"`
	FullNameEn string `db:"full_name_en" json:"full_name_en"`
	UnitName   string `db:"administrative_unit_name" json:"unit_name"`
}

type WardDistrictProvince struct {
	Code       string `db:"code" json:"code"`
	Name       string `db:"name" json:"name"`
	FullName   string `db:"full_name" json:"full_name"`
	FullNameEn string `db:"full_name_en" json:"full_name_en"`
	UnitName   string `db:"administrative_unit_name" json:"unit_name"`
}
