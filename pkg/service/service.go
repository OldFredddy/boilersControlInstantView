package service

import (
	"boilersControlInstantView"
	"boilersControlInstantView/pkg/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Authorization interface {
}

type Service struct {
	Authorization
	RedisRepo *repository.RedisRepository
}

func NewService(repos *repository.Repository, redisRepo *repository.RedisRepository) *Service {
	return &Service{
		Authorization: repos.Authorization,
		RedisRepo:     redisRepo,
	}
}
func (s *Service) GetBoilersFromAPI(apiURL string, imageMap map[int]string) ([]boilersControlInstantView.Boiler, error) {
	var boilers []boilersControlInstantView.Boiler

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&boilers); err != nil {
		fmt.Println(err)
		return nil, err
	}
	for i := range boilers {
		if imageURL, ok := imageMap[boilers[i].ImageResId]; ok {
			boilers[i].ImageURL = imageURL
		} else {
			boilers[i].ImageURL = "/static/images/default.png"
		}
	}
	return boilers, nil
}
func (s *Service) GetRedisData(key string) ([]string, error) {
	return s.RedisRepo.GetData(key)
}
func (s *Service) GetBoilerData(key string) ([]boilersControlInstantView.Boiler, error) {
	rawData, err := s.RedisRepo.GetData(key)
	if err != nil {
		return nil, err
	}

	var boilerData []boilersControlInstantView.Boiler
	for _, jsonData := range rawData {
		var data boilersControlInstantView.Boiler
		err := json.Unmarshal([]byte(jsonData), &data)
		if err != nil {
			continue
		}
		boilerData = append(boilerData, data)
	}

	return boilerData, nil
}
