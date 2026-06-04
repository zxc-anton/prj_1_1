package types

type RequestUser struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=35"`
}

type LoginUser struct {
	Password string `json:"password" validate:"required,min=3,max=35"`
	Email    string `json:"email" validate:"required,email"`
}
