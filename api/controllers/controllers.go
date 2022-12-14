package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/Patrick564/url-shortener-backend/internal/models"
	"github.com/Patrick564/url-shortener-backend/utils"
	"github.com/gin-gonic/gin"
)

type Controllers interface {
	All() ([]models.Url, error)
	Add(url string) (models.Url, error)
	GoTo(id string) (string, error)
}

type Env struct {
	Urls Controllers
}

func (e *Env) UrlsIndex(ctx *gin.Context) {
	u, err := e.Urls.All()
	if err != nil {
		log.Println(err)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"urls": u})
}

func (e *Env) UrlsAdd(ctx *gin.Context) {
	var body struct {
		Url string `json:"url"`
	}

	err := ctx.BindJSON(&body)
	if err != nil {
		if err.Error() == "EOF" {
			log.Println(err)
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": utils.ErrEmptyBody.Error()},
			)
			return
		}

		log.Println(err)
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	u, err := e.Urls.Add(body.Url)
	if err != nil {
		if errors.Is(err, utils.ErrEmptyBody) {
			log.Println(err)
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
			return
		}

		log.Println(err)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{"url": u},
	)
}

func (e *Env) UrlsGoToID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status": http.StatusNotFound,
				"error":  utils.ErrEmptyID.Error(),
			},
		)
		return
	}

	url, err := e.Urls.GoTo(id)
	if err != nil {
		log.Println(err)
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"status": http.StatusBadRequest,
				"error":  err.Error(),
			},
		)
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, url)
}
