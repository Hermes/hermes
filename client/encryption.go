package main

import (
    "fmt"
    "os"
    "crypto/aes"
    "crypto/cipher"
    "io"
)

func Encrypt(in io.Reader, key []byte) io.Reader {

    // Load the plaintext message you want to encrypt.
    plaintext := []byte("hello, world")
    if len(os.Args) > 1 {
        plaintext = []byte(os.Args[1])
    }

    // Setup a key that will encrypt the other text.
    key_text := "32o4908go293hohg98fh40gh"
    if len(os.Args) > 2 {
        key_text = os.Args[2]
    }

    // We chose our cipher type here in this case
    // we are using AES.
    c, err := aes.NewCipher([]byte(key_text))
    if err != nil {
        fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key_text), err)
        os.Exit(-1)
    }

    // We use the CFBEncrypter in order to encrypt
    // the whole stream of plaintext using the
    // cipher setup with c and a iv.
    cfb := cipher.NewCFBEncrypter(c, key)
    ciphertext := make([]byte, len(plaintext))
    cfb.XORKeyStream(ciphertext, plaintext)
    fmt.Printf("%s=>%x\n", plaintext, ciphertext)


    var r io.Reader
    return r
}

func Decrypt(in io.Reader, key []byte) io.Reader {

    // Load the plaintext message you want to encrypt.
    plaintext := []byte("hello, world")
    if len(os.Args) > 1 {
    plaintext = []byte(os.Args[1])
    }

    // Setup a key that will encrypt the other text.
    key_text := "32o4908go293hohg98fh40gh"
    if len(os.Args) > 2 {
    key_text = os.Args[2]
    }

    // We chose our cipher type here in this case
    // we are using AES.
    c, err := aes.NewCipher([]byte(key_text))
    if err != nil {
    fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key_text), err)
    os.Exit(-1)
    }

    // We decrypt it here just for the purpose of
    // showing the fact that it is decryptable.
    cfbdec := cipher.NewCFBDecrypter(c, commonIV)
    plaintextCopy := make([]byte, len(plaintext))
    cfbdec.XORKeyStream(plaintextCopy, ciphertext)
    fmt.Printf("%x=>%s\n",ciphertext, plaintextCopy)

    var r io.Reader
    return r
}

func main() {
    var a io.Reader
    var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
    Encrypt(a, commonIV)
    Decrypt(a, commonIV)
}