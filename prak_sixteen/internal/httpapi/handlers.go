package httpapi

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/CyberGeo335/prak_sixteen/internal/models"
	"github.com/CyberGeo335/prak_sixteen/internal/repo"
	"github.com/CyberGeo335/prak_sixteen/internal/service"
)

type Router struct{ Svc *service.Service }

func (rt Router) Register(r *gin.Engine) {
	r.POST("/notes", rt.createNote)
	r.GET("/notes/:id", rt.getNote)
}

func (rt Router) createNote(c *gin.Context) {
	var in struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}

	n := models.Note{Title: in.Title, Content: in.Content}
	if err := rt.Svc.Create(c.Request.Context(), &n); err != nil {
		// упрощённо: в реале отличаем validation/other
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, n)
}

func (rt Router) getNote(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}

	n, err := rt.Svc.Get(c.Request.Context(), id)
	if err != nil {
		if err == repo.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, n)
}
