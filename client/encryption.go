package client

import (
    "fmt"
    "os"
    "crypto/aes"
    "crypto/cipher"
    "io"
    "bytes"
    "strings"
    "crypto/sha256"
)

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

func Encrypt(in io.Reader, key string) io.Reader {

    // Convert io.Reader to string
    buf := new(bytes.Buffer)
    buf.ReadFrom(in)
    s := buf.String()

    // Load the plaintext message you want to encrypt.
    plaintext := []byte(s)

    // Setup a key that will encrypt the other text.
    h := sha256.New()
    io.WriteString(h, key)
    key_text := h.Sum(nil)

    // We chose our cipher type here in this case we are using AES.
    c, err := aes.NewCipher([]byte(key_text))
    if err != nil {
        fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key_text), err)
        os.Exit(-1)
    }

    // We use the CFBEncrypter in order to encrypt
    // the whole stream of plaintext using the
    // cipher setup with c and a iv.
    cfb := cipher.NewCFBEncrypter(c, commonIV)
    ciphertext := make([]byte, len(plaintext))
    cfb.XORKeyStream(ciphertext, plaintext)
    return strings.NewReader(string(ciphertext))
}

func Decrypt(in io.Reader, key string) io.Reader {

    buf := new(bytes.Buffer)
    buf.ReadFrom(in)
    s := buf.String()

    // Load the ciphertext message you want to decrypt
    ciphertext := []byte(s)
 
    // Setup a key that will encrypt the other text.
    h := sha256.New()
    io.WriteString(h, key)
    key_text := h.Sum(nil)
 
    // We chose our cipher type here in this case
    // we are using AES.
    c, err := aes.NewCipher([]byte(key_text));
    if err != nil {
        fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key_text), err)
        os.Exit(-1)
    }
 
    // We use the CFBDecrypter in order to decrypt
    // the whole stream of ciphertext using the
    // cipher setup with c and a iv.
    cfb := cipher.NewCFBDecrypter(c, commonIV)
    plaintext := make([]byte, len(ciphertext))
    cfb.XORKeyStream(plaintext, ciphertext)
    return strings.NewReader(string(plaintext))
}