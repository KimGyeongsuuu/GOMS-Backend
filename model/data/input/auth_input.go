package input

type SignUpInput struct {
	Email    string
	Password string
	Name     string
	Major    string
	Gender   string
}

type SignInInput struct {
	Email    string
	Password string
}

type SendEmaiInput struct {
	Email string
}

type VerifyAuthCodeInput struct {
	Email    string
	AuthCode string
}
