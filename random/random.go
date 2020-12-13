// random.go
// Package random implements multiple functions for generating random numbers and strings.
// Based on schollz's gist:
// https://gist.github.com/schollz/156d608e8ec26816cedaf06f14d7d692

package random

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"io"
	"log"
)

var number = []byte("0123456789")
var alpha = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
var alphanum = append(number, alpha...)
var special []byte

func init() {
	// assert that a cryptographically secure PRNG is available
	buf := make([]byte, 1)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		log.Panic(err)
	}
	// initialize special character slice
	special = make([]byte, 95)
	for i := 0; i < 95; i++ {
		special[i] = byte(i + 32)
	}
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(length int) ([]byte, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil { // note that err == nil only if we read len(b) bytes
		return nil, err
	}
	return b, nil
}

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(domain []byte, length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}
	if domain == nil { // all bytes: encode it
		return base64.URLEncoding.EncodeToString(bytes), nil
	}
	for i, b := range bytes {
		bytes[i] = domain[b%byte(len(domain))]
	}
	return string(bytes), nil
}

// Number returns a securely generated ramdom string containing numerical characters only.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func Number(length int) (string, error) {
	return GenerateRandomString(number, length)
}

// Alpha returns a securely generated ramdom string containing alphabetical characters only.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func Alpha(length int) (string, error) {
	return GenerateRandomString(alpha, length)
}

// AlphaNum returns a securely generated ramdom string containing alphanumerical characters only.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func AlphaNum(length int) (string, error) {
	return GenerateRandomString(alphanum, length)
}

// Special returns a securely generated ramdom string containing any 7-bit printable characters.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func Special(length int) (string, error) {
	return GenerateRandomString(special, length)
}

// Any returns a securely generated ramdom string containing any 8-bit characters URL-safe base64-encoded.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func Any(length int) (string, error) {
	return GenerateRandomString(nil, length)
}

// UInt returns a securely generated random unsigned integer.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func UInt() (uint, error) {
	bytes, err := GenerateRandomBytes(8)
	if err != nil {
		return 0, err
	}
	return uint(binary.BigEndian.Uint64(bytes)), nil
}

// Hex returns a securely generated random hexadecimal integer.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue
func Hex(length int) (string, error) {
	r, err := GenerateRandomBytes(length/2 + 1)
	if err != nil {
		return "", err
	}
	s := hex.EncodeToString(r)
	return s[:length], nil
}
