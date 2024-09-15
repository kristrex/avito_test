package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingHandler обрабатывает запросы на эндпоинт /api/ping
func PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
