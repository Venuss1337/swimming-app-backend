package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

type Argon2id struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var Argon *Argon2id

func InitConfig(argon2id *Argon2id) {
	Argon = argon2id
}

func generateSalt(length uint32) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

func decodeHash(encodedHash string) (argon *Argon2id, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, errors.New("invalid hash format")
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("invalid version")
	}

	argon = &Argon2id{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &argon.Memory, &argon.Iterations, &argon.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	argon.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	argon.KeyLength = uint32(len(hash))

	return argon, salt, hash, nil
}

func (argon *Argon2id) Hash(cypher *string, plain []byte) error {
	salt, err := generateSalt(argon.SaltLength)
	if err != nil {
		return err
	}
	hash := argon2.IDKey(plain, salt, argon.Iterations, argon.Memory, argon.Parallelism, argon.KeyLength)

	hash = []byte(base64.RawStdEncoding.EncodeToString(hash))
	salt = []byte(base64.RawStdEncoding.EncodeToString(salt))

	*cypher = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argon.Memory, argon.Iterations, argon.Parallelism, salt, hash)
	return nil
}

func (argon *Argon2id) Verify(plain []byte, encoded string) error {
	decodedArgon, decodedSalt, decodedHash, err := decodeHash(encoded)
	if err != nil {
		return err
	}

	newHash := argon2.IDKey(plain, decodedSalt, decodedArgon.Iterations, decodedArgon.Memory, decodedArgon.Parallelism, decodedArgon.KeyLength)
	if subtle.ConstantTimeCompare(newHash, decodedHash) == 1 {
		return nil
	}
	return errors.New("invalid username or password")
}
