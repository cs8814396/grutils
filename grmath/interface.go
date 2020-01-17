package grmath

type Cryptor interface {
	Encrypt(data []byte) (string, error)

	Decrypt(data string) ([]byte, error)
}
