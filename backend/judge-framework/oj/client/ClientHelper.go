package client

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/skmonir/mango-gui/backend/judge-framework/utils"
	"io"
)

func genFtaa() string {
	return utils.RandString(18)
}

func genBfaa() string {
	return "f1b3f18c715565b589b7823cda7448ce"
}

func createHash(key string) []byte {
	hashFunc := md5.New()
	hashFunc.Write([]byte(key))
	return hashFunc.Sum(nil)
}

func encrypt(handle, password string) (ret string, err error) {
	block, err := aes.NewCipher(createHash("tiutiu" + handle + "477"))
	if err != nil {
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}
	text := gcm.Seal(nonce, nonce, []byte(password), nil)
	ret = hex.EncodeToString(text)
	return
}

func decrypt(handle, password string) (ret string, err error) {
	data, err := hex.DecodeString(password)
	if err != nil {
		err = errors.New("Cannot decode the password")
		return
	}
	block, err := aes.NewCipher(createHash("tiutiu" + handle + "477"))
	if err != nil {
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}
	nonceSize := gcm.NonceSize()
	nonce, text := data[:nonceSize], data[nonceSize:]
	plain, err := gcm.Open(nil, nonce, text, nil)
	if err != nil {
		return
	}
	ret = string(plain)
	return
}

func getSubmitUrl(host, url string) string {
	oj, cid, _, ctype := utils.ExtractInfoFromUrl(url)
	if oj == "codeforces" {
		return fmt.Sprintf(host+"/%v/%v/submit", ctype, cid)
	}
	return ""
}
