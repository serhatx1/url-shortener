package service

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"url-shortener/cache"
)

type URLService struct {
	Repo URLRepository
}

type URLRepository interface {
	Save(originalURL string, shortURL string) error
	IfLongUrlExist(originalURL string) (string, bool)
	FindOriginal(shortURL string) (string, bool)
}

func NewURLService(repo URLRepository) *URLService {
	return &URLService{Repo: repo}
}

func (s *URLService) ShortenURL(originalURL *string) (string, error) {

	if *originalURL == "" {
		return "", errors.New("original URL cannot be empty")
	}
	*originalURL = strings.Replace(*originalURL, "www.", "", 1)
	s.ensureHTTPS(originalURL)

	existUrl, exist := s.Repo.IfLongUrlExist(*originalURL)
	if exist {
		return existUrl, errors.New("the url already exist")
	}

	shortURL := s.Encode(originalURL)
	s.Repo.Save(*originalURL, shortURL)
	return shortURL, nil

}
func (s *URLService) Redirect(shortURL *string) (string, error) {

	if *shortURL == "" {
		return "", errors.New("URL cannot be empty")
	}
	parsedURL, err := url.Parse(*shortURL)
	if err != nil {
		return "", errors.New("caused in parse")
	}
	log.Print(parsedURL)
	pathArr := make([]string, 0)
	trimmedPath := strings.TrimSuffix(parsedURL.Path, "/")
	pathSegments := strings.Split(trimmedPath, "/")
	joinedPath := strings.Join(pathSegments, "/")
	log.Print("joined", joinedPath)
	var originalURL string
	for len(pathSegments) > 0 {
		var url string
		var exist bool
		url, exist = cache.Get(joinedPath)
		if !exist {
			url, exist = s.Repo.FindOriginal(joinedPath)
		}
		if !exist {
			pathSegments = strings.Split(joinedPath, "/")
			pathArr = append(pathArr, pathSegments[len(pathSegments)-1])
			pathSegments = pathSegments[:len(pathSegments)-1]
			joinedPath = strings.Join(pathSegments, "/")
		}
		if exist {
			originalURL = url
			break
		}

	}
	if originalURL == "" {
		return "", errors.New("url could not find")
	}
	originalURL = originalURL + "/" + strings.Trim(strings.Join(reverse(pathArr), "/"), "/")
	cache.Add(*shortURL, originalURL)
	return originalURL, nil
}
func (s *URLService) GetOriginalURL(shortURL string) (string, error) {

	return "", nil
}
func (h *URLService) ensureHTTPS(url *string) string {
	if !strings.HasPrefix(*url, "http://") && !strings.HasPrefix(*url, "https://") {
		*url = "https://" + *url
		return *url
	}
	return *url
}

func (s *URLService) Encode(originalURL *string) string {
	hash := sha1.Sum([]byte(*originalURL))
	encoded := base64.RawURLEncoding.EncodeToString(hash[:])
	shortURL := fmt.Sprintf("%s/%s", "http://localhost:8080", encoded[:8])
	return shortURL
}

func (s *URLService) Decode(shortURL string) string {
	return ""
}
func reverse(slice []string) []string {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}
