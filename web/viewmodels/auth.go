package viewmodels

type RegisterForm struct {
	Username string `form:"username" json:"username" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required,email"`
}

type RegisterPageData struct {
	Title       string
	Error       string
	Success     string
	FormValues  RegisterForm
	FieldErrors map[string]string
}

type ValidationResponse struct {
	IsValid bool
	Message string
}

type LoginForm struct {
	Username string `form:"username" json:"username" binding:"required"`
}
