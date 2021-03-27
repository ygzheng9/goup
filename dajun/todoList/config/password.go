package config

import "golang.org/x/crypto/bcrypt"

// HashPassword bcrypt 对密码进行加密
func HashPassword(password string) (string, error) {
	// cost 会极大影响速度， 默认是 10，如果增大到 17，就需要 5~10s
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	return string(bytes), err
}

// CheckPasswordHash bcrypt 比较输入密码是否匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// EncryptDecrypt runs a XOR encryption on the input string, encrypting it if it hasn't already been,
// and decrypting it if it has, using the key provided.
func EncryptDecrypt(input, key string) (output string) {
	for i := 0; i < len(input); i++ {
		output += string(input[i] ^ key[i%len(key)])
	}

	return output
}
