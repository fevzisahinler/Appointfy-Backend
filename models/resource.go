package models

type Resource struct {
	ID           uint         `gorm:"column:resource_id;primaryKey" json:"id"`
	ResourceName string       `gorm:"unique;not null" json:"resource_name"`
	Description  string       `json:"description"`
	Permissions  []Permission `gorm:"foreignKey:ResourceID" json:"-"`
}
