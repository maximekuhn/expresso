package auth

import "golang.org/x/crypto/bcrypt"

func hashPassword(p string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
}

// return true if the password is correct
func checkPassword(clear string, hash []byte) bool {
	return bcrypt.CompareHashAndPassword(hash, []byte(clear)) == nil
}
