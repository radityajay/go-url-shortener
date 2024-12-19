package url

import (
	"fmt"
	"os"
	"time"

	"math/rand"

	"github.com/radityajay/go-url-shortener/db"
	"github.com/radityajay/go-url-shortener/models"
)

func CreateGenerateURL(longURL string, customAlias *string) (*string, error) {
	length := 6
	shortCode := ""
	if customAlias != nil {
		shortCode = *customAlias
	} else {
		shortCode = generateShortCode(length)
	}

	url := models.URL{
		Model: models.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		LongURL:      longURL,
		ShortURL:     shortCode,
		AccessCount:  0,
		LastAccessed: time.Now(),
	}

	err := url.CreateURL(db.Postgres.DB)
	if err != nil {
		return nil, err
	}

	shortenUrl := fmt.Sprintf("%s/s/%s", os.Getenv("BASE_URL"), url.ShortURL)
	return &shortenUrl, nil
}

func generateShortCode(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)

}

func GetRedirectURL(shortCode string) (string, error) {
	res, err := models.GetURL(db.Postgres.DB, shortCode)
	if err != nil {
		return "", err
	}
	return res.LongURL, nil
}
