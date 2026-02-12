package service

import (
	"crypto/sha256"
	"ozon_link_shortener/internal/storage"
)

type Service struct {
	Cache storage.Storage
	DB    storage.Storage
	Conf  *Config
}

type Config struct {
	DBAddress string
	Storage   string
}

type Request struct {
	LongUrl string `json:"url"`
}

func Shortener(longUrl string) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	sha := sha256.New()
	shortUrl := []byte{}

	sha.Write([]byte(longUrl))
	hash := sha.Sum(nil)
	hash = hash[:10]

	for i := range hash {
		idx := int(hash[i]) % len(chars)
		shortUrl = append(shortUrl, chars[idx])
	}
	return string(shortUrl)
}

func WriteToStorage(s *Service, shortUrl string, req Request) error {
	err := s.Cache.Set(shortUrl, req.LongUrl)
	if err != nil {
		return err
	}
	if s.Conf.Storage == "postgres" {
		err := s.DB.Set(shortUrl, req.LongUrl)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetFromStorage(s *Service, shortUrl string) (string, error) {
	var err error
	longUrl, err := s.Cache.Get(string(shortUrl))
	if err != nil && s.Conf.Storage == "postgres" {
		longUrl, err = s.DB.Get(shortUrl)
		if err != nil {
			return "", err
		}
	}
	return longUrl, err
}
