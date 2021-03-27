package config

import (
	"encoding/base64"
	"fmt"
	"log"
	"testing"
)

func Test_EncryptDecrypt(t *testing.T) {
	data := "mysql"

	encrypted := EncryptDecrypt(data, passwordKey)
	fmt.Printf("Encrypted: %s\n", encrypted)

	decrypted := EncryptDecrypt(encrypted, passwordKey)
	fmt.Printf("Decrypted: %s\n", decrypted)

	// fmt.Printf("%+v\n", appConfig)
}

func Test_DecodeString(t *testing.T) {
	input := []byte("Z&9I^0x)9D*6")

	// 演示base64编码
	encodeString := base64.StdEncoding.EncodeToString(input)
	fmt.Println(encodeString)

	// 对上面的编码结果进行base64解码
	decodeBytes, err := base64.StdEncoding.DecodeString(encodeString)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(decodeBytes))
}
