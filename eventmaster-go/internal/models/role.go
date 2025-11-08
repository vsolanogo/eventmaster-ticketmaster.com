package models

type Role struct {
	Base
	Name        string  `gorm:"size:100;not null;unique" validate:"required"`
	Description *string `gorm:"type:text"`
	Users       []*User `gorm:"many2many:user_roles;"`
}

// TableName specifies the table name for the Role model
func (Role) TableName() string {
	return "roles"
}

type RoleResponse struct {
	Role        string  `json:"role"`
	Description *string `json:"description,omitempty"`
}

func (r *Role) ToResponse() RoleResponse {
	return RoleResponse{
		Role:        r.Name,
		Description: r.Description,
	}
}
