package main

import (
	"boilersControlInstantView"
	"boilersControlInstantView/pkg/handler"
	"boilersControlInstantView/pkg/repository"
	"boilersControlInstantView/pkg/service"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initialiing configs:%s", err.Error())
	}
	imageMap, err := loadImageMap()
	redisRepo := repository.NewRedisRepository(
		viper.GetString("redis.host"),
		viper.GetString("redis.port"),
		viper.GetString("redis.password"),
	)
	repos := repository.NewRepository()
	services := service.NewService(repos, redisRepo)
	handlers := handler.NewHandler(services, imageMap)
	srv := new(boilersControlInstantView.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error)
	}
	if err != nil {
		log.Fatalf("Failed to load image map: %v", err)
	}
	fmt.Println("Сервер запущен на http://95.142.45.133:23874/")
	fmt.Println("<Debug> Сервер запущен на http://localhost:23874/boiler-room")
	if err != nil {
		fmt.Println("Ошибка запуска сервера: ", err)
	}
}
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func loadImageMap() (map[int]string, error) {
	fileMap := make(map[int]string)
	files, err := ioutil.ReadDir("static/images")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fileName := file.Name()

		id, err := strconv.Atoi(strings.Split(strings.TrimSuffix(fileName, ".png"), "_")[2])
		if err != nil {
			continue
		}
		fileMap[id] = "/static/images/" + fileName
	}

	return fileMap, nil
}
