package user

type PasswordAndConfirmationDontMatchError struct{}

func (_ PasswordAndConfirmationDontMatchError) Error() string {
	return "password and confirmation don't match"
}
