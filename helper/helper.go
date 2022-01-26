package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"freq/config"
	"freq/models"
	"io"
	mRand "math/rand"
	"regexp"
	"strings"
	"sync"
	"time"
	"unsafe"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var priceRegex = regexp.MustCompile("(\\d+\\.\\d{1,2})")

var src = mRand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

type Encryption struct {
	Key []byte
}

func (e *Encryption) Encrypt(text string) (string, error) {
	plaintext := []byte(text)

	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (e *Encryption) Decrypt(cryptoText string) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func RandomString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func IsEmail(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func IsValidPrice(e string) bool {
	if len(e) < 1 {
		return false
	}
	return priceRegex.MatchString(e)
}

func ExtractData(token string) ([]string, error) {
	if token == "" {
		return nil, fmt.Errorf("no token provided")
	}
	xs := strings.Split(token, " ")

	if len(xs) != 2 {
		return nil, fmt.Errorf("invalid token provided")
	}

	tokenValue := strings.Split(xs[1], "|")
	return tokenValue, nil
}

func EncryptPI(purchase *models.Purchase) *models.Purchase {
	key := config.Config("KEY")

	encrypt := Encryption{Key: []byte(key)}

	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		defer wg.Done()
		purchase.Shipped = false
		purchase.Delivered = false
		purchase.TrackingId = ""
		purchase.CreatedAt = time.Now()
		purchase.UpdatedAt = time.Now()
	}()

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(purchase.StreetAddress)

		if err != nil {
			panic(err)
		}

		purchase.StreetAddress = pi
	}()

	go func() {
		defer wg.Done()

		if len(purchase.OptionalAddress) > 0 {
			pi, err := encrypt.Encrypt(purchase.OptionalAddress)

			if err != nil {
				panic(err)
			}

			purchase.OptionalAddress = pi
		}
	}()

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(purchase.City)

		if err != nil {
			panic(err)
		}

		purchase.City = pi
	}()

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(purchase.ZipCode)

		if err != nil {
			panic(err)
		}

		purchase.ZipCode = pi
	}()

	wg.Wait()

	return purchase
}

func DecryptPI(purchase *models.Purchase) *models.Purchase {
	key := config.Config("KEY")

	decrypt := Encryption{Key: []byte(key)}

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		pi, err := decrypt.Decrypt(purchase.StreetAddress)

		if err != nil {
			panic(err)
		}

		purchase.StreetAddress = pi
	}()

	go func() {
		defer wg.Done()

		if len(purchase.OptionalAddress) > 0 {
			pi, err := decrypt.Decrypt(purchase.OptionalAddress)

			if err != nil {
				panic(err)
			}

			purchase.OptionalAddress = pi
		}
	}()

	go func() {
		defer wg.Done()
		pi, err := decrypt.Decrypt(purchase.City)

		if err != nil {
			panic(err)
		}

		purchase.City = pi
	}()

	go func() {
		defer wg.Done()
		pi, err := decrypt.Decrypt(purchase.ZipCode)

		if err != nil {
			panic(err)
		}

		purchase.ZipCode = pi
	}()

	wg.Wait()

	return purchase
}

func CreateEmail(email *models.Email, emailDto *models.EmailDto, emailType models.EmailType) *models.Email {
	emailAdd := config.Config("BUSINESS_EMAIL")

	email.From = emailAdd
	email.CustomerEmail = emailDto.Email
	email.Content = emailDto.Content
	email.Subject = emailDto.Subject
	email.Status = models.Pending
	email.Type = emailType
	email.CreatedAt = time.Now()
	email.UpdatedAt = time.Now()

	return email
}

func ValidateData(data string, dataType string, min int, max int) error {
	if len(data) >= min && len(data) <= max {
		return nil
	}

	return fmt.Errorf("%s must be at least %v and at most %v characters", dataType, min, max)
}
