package application

import (
  "io"
  "encoding/base64"
  "crypto/aes"
  "crypto/rand"
  "crypto/cipher"
)

func Encrypt(key, text []byte) string {
    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }
    ciphertext := make([]byte, aes.BlockSize+len(text))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        panic(err)
    }
    cfb := cipher.NewCFBEncrypter(block, iv)
    cfb.XORKeyStream(ciphertext[aes.BlockSize:], text)
    return encodeBase64(ciphertext)
}

func encodeBase64(b []byte) string {
    return base64.StdEncoding.EncodeToString(b)
}
