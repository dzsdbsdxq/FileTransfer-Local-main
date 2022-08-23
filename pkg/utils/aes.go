package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"sort"
)

const AES_KEY = "sykKwe59_q11peDz"

func PostDataCheckAuth(c *gin.Context, fields []string) bool {
	data := make(map[string]string)
	for _, v := range fields {
		data[v] = com.StrTo(c.PostForm(v)).String()
	}
	data["_t"] = com.StrTo(c.PostForm("_t")).String()
	authString, _ := com.Base64Decode(com.StrTo(c.PostForm("sign")).String())
	signStr := customMd5(aesEncryptCBC(sortToStr(data)))
	return signStr == authString
}

func sortToStr(data map[string]string) string {
	var dataParams string
	//ksort
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	//拼接
	for _, k := range keys {
		dataParams = dataParams + k + "=" + data[k] + "&"
	}
	return dataParams[0 : len(dataParams)-1]
}

func aesEncryptCBC(origDataStr string) string {
	origData := []byte(origDataStr)
	key := []byte(AES_KEY)
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)                // 补全码
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
	encrypted := make([]byte, len(origData))                    // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                  // 加密
	return hex.EncodeToString(encrypted)
}
func aesDecryptCBC(ciphertext string) (string, error) {
	encrypted, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	key := []byte(AES_KEY)
	block, _ := aes.NewCipher(key)                              // 分组秘钥
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
	decrypted := make([]byte, len(encrypted))                   // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)                 // 解密
	decrypted = pkcs5UnPadding(decrypted)                       // 去除补全码
	return string(decrypted), nil
}
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

// Md5 返回一个32位md5加密后的字符串
func customMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
