package models

type Role struct {
	Base
	Name  string  `gorm:"size:100;not null;unique" validate:"required"`
	Users []*User `gorm:"many2many:user_roles;"`
}

// TableName specifies the table name for the Role model
func (Role) TableName() string {
	return "roles"
}
