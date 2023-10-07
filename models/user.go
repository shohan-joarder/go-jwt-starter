package models

type UserValidations struct {
	FirstName string `json:"first_name" form:"first_name" validate:"required,min=5"`
	LastName  string `json:"last_name" form:"last_name" validate:"required,min=5"`
	Email     string `json:"email" form:"email" validate:"required,email"`
	Phone     string `json:"phone" form:"phone" validate:"required,min=10,max=12"`
	Password  string `json:"password" form:"password" validate:"required,min=8,max=20"`
	// DOB       *time.Time `json:"dob" form:"date" validate:"omitempty"`
}
type UserRegistrations struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	// DOB       *time.Time `json:"dob" form:"date" validate:"omitempty"`
}

type Login struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8,max=20"`
}