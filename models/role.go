package models

type Role struct {
	ID          uint         `gorm:"column:role_id;primaryKey" json:"id"`
	RoleName    string       `gorm:"unique;not null" json:"role_name"`
	Description string       `json:"description"`
	Users       []User       `gorm:"many2many:user_roles;" json:"-"`
	Permissions []Permission `gorm:"foreignKey:RoleID" json:"permissions"` // Roller ile izinleri ilişkilendiren alan
}
