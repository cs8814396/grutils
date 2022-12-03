package grmath

import (
	"crypto/rc4"
	"encoding/base64"
	"encoding/hex"
)

type Aes256Rc4Base64Cryptor struct {
	AesKey []byte
	AesIv  []byte
	Rc4Key []byte
}

func (this *Aes256Rc4Base64Cryptor) Encrypt(data []byte) (encryptData string, err error) {
	return Aes256Rc4Base64Encrypt(data, this.AesKey, this.AesIv, this.Rc4Key)
}

func (this *Aes256Rc4Base64Cryptor) Decrypt(data string) (decryptData []byte, err error) {
	return Aes256Rc4Base64Decrypt(data, this.AesKey, this.AesIv, this.Rc4Key)
}

func Aes256Rc4Base64Encrypt(data []byte, key []byte, iv []byte, rc4Key []byte) (encryptData string, err error) {
	aesData, err := AesCBCZeroPaddingEncrypt(data, key, iv)
	if err != nil {
		return
	}

	rc4obj, err := rc4.NewCipher(rc4Key)
	if err != nil {
		return
	}

	rc4Result := make([]byte, len(aesData))

	rc4obj.XORKeyStream(rc4Result, aesData)

	encryptData = base64.StdEncoding.EncodeToString(rc4Result)

	return

}
func Aes256Rc4Base64Decrypt(data string, key []byte, iv []byte, rc4Key []byte) (decryptData []byte, err error) {
	decryptBase64, err := base64.StdEncoding.DecodeString(data)
	if err != nil {

		return
	}

	decryptRc4 := make([]byte, len(decryptBase64))

	rc4obj, err := rc4.NewCipher(rc4Key)
	if err != nil {

		return
	}

	rc4obj.XORKeyStream(decryptRc4, decryptBase64)

	decryptData, err = AesCBCZeroPaddingDecrypt(decryptRc4, key, iv)

	if err != nil {

		return
	}
	return

}

/*
import hashlib


class RC4:
    def __init__(self, public_key):
        self.public_key = public_key
        self.public_key = self.public_key.encode("utf-8")
        self.public_key = hashlib.md5(self.public_key).hexdigest()

        print(self.public_key)

    def encode(self, text):
        return self.__docrypt(text)

    def decode(self, text):
        return self.__docrypt(text)

    def __docrypt(self, text):

        result = ''
        box = list(range(256))
        randkey = []
        key_lenth = len(self.public_key)

        for i in range(255):
            randkey.append(ord(self.public_key[i % key_lenth]))

        for i in range(255):
            j = 0
            j = (j + box[i] + randkey[i]) % 256
            tmp = box[i]
            box[i] = box[j]
            box[j] = tmp
        for i in range(len(text)):
            a = j = 0
            a = (a + 1) % 256
            j = (j + box[a]) % 256
            tmp = box[a]
            box[a] = box[j]
            box[j] = tmp
            result += chr(ord(text[i]) ^ (box[(box[a] + box[j]) % 256]))
        return result

class AesRC4Base64(object):
    """
     aes+rc4+base64
    """

    def __init__(self, key='key for aes333', rc4_key='key for rc4'):
        self.key = key

        # if we can't get special key for rc4, then we use key
        self.rc4 = RC4(rc4_key if rc4_key else self.key)

        # we use key utf8 format as the key for aes
        self.aes_total_key = hashlib.md5(self.key.encode("utf-8")).hexdigest()
        self.aes_iv = self.aes_total_key[16:32].encode('utf-8')
        self.aes_key = self.aes_total_key[:16].encode('utf-8')

        self.length = 16

        self.aes = Aes(self.aes_key, self.aes_iv)

    def encode(self, text):

        aes_text = b2a_hex(self.aes.encode(text))
        ciphertext = self.rc4.encode(aes_text.decode('utf8')).encode('utf8')
        return base64.b64encode(ciphertext).decode('utf8')

    def decode(self, text):
        text = base64.b64decode(text.encode('utf8')).decode('utf8')

        text = self.rc4.decode(text)

        plain_text = self.aes.decode(a2b_hex(text))
        return plain_text

*/
type Aes256Rc5Base64Cryptor struct {
	AesKey []byte
	AesIv  []byte
	Rc4Key []byte
}

func (this *Aes256Rc5Base64Cryptor) Encrypt(data []byte) (encryptData string, err error) {
	return Aes256Rc5Base64Encrypt(data, this.AesKey, this.AesIv, this.Rc4Key)
}

func (this *Aes256Rc5Base64Cryptor) Decrypt(data string) (decryptData []byte, err error) {
	return Aes256Rc5Base64Decrypt(data, this.AesKey, this.AesIv, this.Rc4Key)
}

func Aes256Rc5Base64Encrypt(data []byte, key []byte, iv []byte, rc4Key []byte) (encryptData string, err error) {
	tmpBytes, err := AesCBCZeroPaddingEncrypt(data, key, iv)
	if err != nil {
		return
	}
	aesDataEncodeToString := hex.EncodeToString(tmpBytes)
	tmpBytes = []byte(aesDataEncodeToString)

	/*rc4obj, err := rc4.NewCipher(rc4Key)
	if err != nil {
		return
	}

	rc4Result := make([]byte, len(tmpBytes))

	rc4obj.XORKeyStream(rc4Result, tmpBytes)*/

	tmpBytes = Rc4Decode(tmpBytes, rc4Key)

	tmpBytes = EncodeUnicodeToUtf8(tmpBytes)

	encryptData = base64.StdEncoding.EncodeToString(tmpBytes)

	return

}
func Aes256Rc5Base64Decrypt(data string, key []byte, iv []byte, rc4Key []byte) (decryptData []byte, err error) {
	tmpBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {

		return
	}

	tmpBytes = DecodeUtf8ToUnicode(tmpBytes)

	/*

		decryptRc4 := make([]byte, len(tmpBytes))

		rc4obj, err := rc4.NewCipher(rc4Key)
		if err != nil {

			return
		}

		rc4obj.XORKeyStream(decryptRc4, tmpBytes)*/

	tmpBytes = Rc4Decode(tmpBytes, rc4Key)

	tmpBytes, err = hex.DecodeString(string(tmpBytes))
	if err != nil {
		return
	}

	decryptData, err = AesCBCZeroPaddingDecrypt(tmpBytes, key, iv)

	if err != nil {

		return
	}
	return

}
