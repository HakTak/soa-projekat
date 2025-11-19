package models

type UserNoPassDTO struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     Role   `json:"role"`
	Blocked  bool   `json:"blocked"`
}

type UserRegisterDTO struct {
	Username string `json:"username" gorm:"not null;unique"`
	Password string `json:"password"`
	Email    string `json:"email" gorm:"unique"`
	Role     Role   `json:"role"`
}

type LoginRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	Token string        `json:"token"`
	User  UserNoPassDTO `json:"user"`
}
