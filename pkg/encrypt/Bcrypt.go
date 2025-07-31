package encrypt

import (
	"fmt"
	"golang.org/x/crypto/bcrypt" // bcrypt算法加密库
	"log"
)

func Bcrypt() {
	// 原始密码
	password := "paipai123456"

	// 生成Bcrypt哈希
	hashedPassword, err := HashPassword(password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Hashed Password: %s\n", hashedPassword)

	// 验证密码
	fmt.Println("Check password 'paipai123456':", CheckPassword(password, hashedPassword))
	fmt.Println("Check password 'wrongPassword':", CheckPassword("wrongPassword", hashedPassword))
}

// hashPassword 使用Bcrypt算法对密码进行哈希处理
func HashPassword(password string) (string, error) {
	// GenerateFromPassword的第二个参数是cost值，范围4-31，值越大计算越慢但更安全
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// checkPassword 检查密码是否匹配Bcrypt哈希
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
