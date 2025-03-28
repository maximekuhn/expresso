package group

func ValidateName(name string) error {
	if len(name) == 0 {
		return GroupValidationError{
			Field:  "Name",
			Reason: "empty",
		}
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) == 0 {
		return GroupValidationError{
			Field:  "Password",
			Reason: "empty",
		}
	}
	return nil
}
