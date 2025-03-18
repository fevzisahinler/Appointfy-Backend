package models

type Permission struct {
	ID         uint     `gorm:"column:id;primaryKey" json:"id"`
	RoleID     uint     `gorm:"not null" json:"role_id"`
	ResourceID uint     `gorm:"not null" json:"resource_id"`
	CanView    bool     `gorm:"default:false" json:"can_view"`
	CanEdit    bool     `gorm:"default:false" json:"can_edit"`
	CanDelete  bool     `gorm:"default:false" json:"can_delete"`
	CanCreate  bool     `gorm:"default:false" json:"can_create"`
	Role       Role     `gorm:"foreignKey:RoleID" json:"-"`
	Resource   Resource `gorm:"foreignKey:ResourceID" json:"-"`
}
