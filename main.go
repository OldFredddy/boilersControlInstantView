package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	imageMap, err := loadImageMap()
	if err != nil {
		log.Fatalf("Failed to load image map: %v", err)
	}
	http.HandleFunc("/boiler-room", func(w http.ResponseWriter, r *http.Request) {
		boilerRoomHandler(w, r, imageMap)
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Сервер запущен на http://95.142.45.133:23874/")
	fmt.Println("<Debug> Сервер запущен на http://localhost:23874/boiler-room")
	err = http.ListenAndServe(":23874", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера: ", err)
	}
}
func boilerRoomHandler(w http.ResponseWriter, r *http.Request, imageMap map[int]string) {
	boilers, err := getBoilersFromAPI("http://95.142.45.133:23873/getparams", imageMap)
	if err != nil {
		log.Printf("Не удалось получить данные: %v", err)
		http.Error(w, "Не удалось получить данные: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/boiler-rooms.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, boilers)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		return
	}
}

type Boiler struct {
	IsOk          int    `json:"isOk"`
	TPod          string `json:"tPod"`
	PPod          string `json:"pPod"`
	TUlica        string `json:"tUlica"`
	TPlan         string `json:"tPlan"`
	TAlarm        string `json:"tAlarm"`
	ImageResId    int    `json:"imageResId"`
	PPodLowFixed  string `json:"pPodLowFixed"`
	PPodHighFixed string `json:"pPodHighFixed"`
	TPodFixed     string `json:"tPodFixed"`
	ID            int    `json:"id"`
	Version       int64  `json:"version"`
	LastUpdated   int64  `json:"lastUpdated"`
	ImageURL      string `json:"-"`
}

func getBoilersFromAPI(apiURL string, imageMap map[int]string) ([]Boiler, error) {
	var boilers []Boiler

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
