package web

import (
	"net/http"
	"url_shortener/internals/handler"
	"url_shortener/internals/models"

	"github.com/gin-gonic/gin"
)

type EntryWebHandler struct {
	entryHandler *handler.EntryHandler
}

func NewEntryWebHandler(entryHandler *handler.EntryHandler) *EntryWebHandler {
	return &EntryWebHandler{entryHandler: entryHandler}
}

func (h *EntryWebHandler) GetEntry(c *gin.Context) {
	short_url := c.Param("short_url")
	fullUrl, err := h.entryHandler.RetrieveUrl(short_url)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"url": fullUrl})

}

func (h *EntryWebHandler) AddUrlEntry(c *gin.Context) {
	var newEntry models.Entry
	err := c.BindJSON(&newEntry)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
	}
	entry, err := h.entryHandler.AddUrlEntry(newEntry.Url)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err})
	}

	c.IndentedJSON(http.StatusOK, entry)
}

func (h *EntryWebHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
