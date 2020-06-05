package forms

type SignUpForm struct {
	Name      string `json:"name" binding:"required,gt=5,lt=20"`
	Password  string `json:"password" binding:"required,gt=5,lt=20"`
	Password2 string `json:"password2" binding:"required,gt=5,lt=20"`
}

type LoginForm struct {
	Name     string `json:"name" binding:"required,gt=5,lt=20"`
	Password string `json:"password" binding:"required,gt=5,lt=20"`
}

type RefreshTokenForm struct {
	Token string `json:"token" binding:"required"`
}
