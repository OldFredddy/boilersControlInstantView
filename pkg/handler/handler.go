package handler

import (
	"boilersControlInstantView/pkg/service"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"reflect"
)

type Handler struct {
	services *service.Service
	imageMap map[int]string
}

func NewHandler(services *service.Service, imageMap map[int]string) *Handler {
	return &Handler{services: services, imageMap: imageMap}
}
func (h *Handler) boilerRoomHandler(c *gin.Context) {
	boilers, err := h.services.GetBoilersFromAPI("http://95.142.45.133:23873/getparams", h.imageMap)
	if err != nil {
		log.Printf("Не удалось получить данные: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	tmpl, err := template.ParseFiles("templates/boiler-rooms.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = tmpl.Execute(c.Writer, boilers)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
func (h *Handler) getBoilerDataByField(c *gin.Context, field string) {
	boilerID := c.Param("id")
	key := "history:boiler:" + boilerID
	data, err := h.services.GetBoilerData(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var values []string
	for _, entry := range data {
		v := reflect.ValueOf(entry)
		valueField := v.FieldByName(field)
		if valueField.IsValid() {
			values = append(values, valueField.String())
		}
	}
	c.JSON(http.StatusOK, values)
}

func (h *Handler) getTPod(c *gin.Context) {
	h.getBoilerDataByField(c, "TPod")
}

func (h *Handler) getPPod(c *gin.Context) {
	h.getBoilerDataByField(c, "PPod")
}
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Static("/static", "./static")
	router.Use(CORSMiddleware())
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}
	mobileView := router.Group("/mobile")
	{
		mobileView.GET("/", h.updateItem)
	}
	router.GET("/boiler-room", h.boilerRoomHandler)
	router.GET("/boilers/gettpod/:id", h.getTPod)
	router.GET("/boilers/getppod/:id", h.getPPod)
	return router
}
