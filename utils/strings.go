package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"time"
)

// GenerateID generates a string to be used as the resource ID
func GenerateID() string {
	idLength := 15
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	id := make([]rune, idLength)
	for i := range idLength {
		randomInt, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		id[i] = chars[int(randomInt.Int64())]
	}
	return string(id)
}

// GenerateEtag generates a random etag which is generated on every write action
func GenerateEtag() (string, error) {
	// Get a random number and the time in seconds and milliseconds
	timeNow := time.Now()
	randomInt, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		return "", fmt.Errorf("error generating random int: [%w]", err)
	}
	randomNumber := fmt.Sprintf("%07d", randomInt)
	timeInMillis := strconv.FormatInt(timeNow.UnixNano()/int64(time.Millisecond), 10)
	timeInSeconds := strconv.FormatInt(timeNow.UnixNano()/int64(time.Second), 10)
	// Calculate a SHA-512 truncated digest
	shaDigest := sha512.New512_224()
	_, err = shaDigest.Write([]byte(randomNumber + timeInMillis + timeInSeconds))
	if err != nil {
		return "", fmt.Errorf("error writing sha digest: [%w]", err)
	}
	sha1Hash := hex.EncodeToString(shaDigest.Sum(nil))
	return sha1Hash, nil
}
