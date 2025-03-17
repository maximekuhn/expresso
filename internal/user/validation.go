package user

func ValidateName(name string) error {
	if len(name) == 0 {
		return UserValidationError{
			Field:  "Name",
			Reason: "empty",
		}
	}
	return nil
}
