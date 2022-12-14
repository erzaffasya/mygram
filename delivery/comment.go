package delivery

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/erzaffasya/mygram/middlewares"
	"github.com/erzaffasya/mygram/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type commentRoutes struct {
	cuc models.CommentUsecase
	puc models.PhotoUsecase
}

func NewCommentRoute(handlers *gin.Engine, cuc models.CommentUsecase, puc models.PhotoUsecase) {
	route := &commentRoutes{cuc, puc}

	handler := handlers.Group("/comments")
	{
		handler.Use(middlewares.Authentication())
		handler.GET("/", route.Fetch)
		handler.POST("/", route.Store)
		handler.PUT("/:id", middlewares.CommentEditDelAuthorization(route.cuc), route.Update)
		handler.DELETE("/:id", middlewares.CommentEditDelAuthorization(route.cuc), route.Delete)
	}
}

// ? perlu melewati proses autentikasi dan autorisasi terlebih dahulu
// ? hanya get comment sendiri?
// Fetch godoc
// @Summary      Fetch comments
// @Description  get comments
// @Tags         comments
// @Accept       json
// @Produce      json
// @Success      200	{object}	[]models.Comment
// @Failure      400	{object}	ErrorResponse
// @Failure      401	{object}	ErrorResponse
// @Security     Bearer
// @Router       /comments      [get]
func (route *commentRoutes) Fetch(c *gin.Context) {
	var (
		comments []models.Comment
		err      error
	)

	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	err = route.cuc.Fetch(c.Request.Context(), &comments, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// Store godoc
// @Summary      Create an comment
// @Description  create and store an comment
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        json  body  models.Comment true  "Comment"
// @Success      201  {object}  models.Comment
// @Failure      400	{object}	ErrorResponse
// @Failure      401	{object}	ErrorResponse
// @Security     Bearer
// @Router       /comments      [post]
func (route *commentRoutes) Store(c *gin.Context) {
	var (
		comment models.Comment
		photo   models.Photo
		err     error
	)

	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	err = c.ShouldBindJSON(&comment)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	photoID := uint(comment.PhotoID)

	err = route.puc.GetByID(c.Request.Context(), &photo, photoID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": fmt.Sprintf("Photo with id %d doesn't exist", photoID),
		})
		return
	}

	comment.UserID = userID

	err = route.cuc.Store(c.Request.Context(), &comment)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         comment.ID,
		"user_id":    comment.UserID,
		"photo_id":   comment.PhotoID,
		"message":    comment.Message,
		"created_at": comment.CreatedAt,
	})
}

// Update godoc
// @Summary      Update an comment
// @Description  update an comment by ID
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Comment ID"
// @Param        json  body  models.Comment true  "Comment"
// @Success      200  {object}  models.Photo
// @Failure      400  {object}	ErrorResponse
// @Failure      401  {object}	ErrorResponse
// @Failure      404  {object}	ErrorResponse
// @Security     Bearer
// @Router       /comments/{id} [put]
func (route *commentRoutes) Update(c *gin.Context) {
	var (
		comment models.Comment
		photo   models.Photo
		err     error
	)

	commentIDStr := c.Param("id")
	commentIDInt, err := strconv.Atoi(commentIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Failed to cast comment id to int",
		})
		return
	}

	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	err = c.ShouldBindJSON(&comment)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	commentID := uint(commentIDInt)

	updatedComment := models.Comment{
		UserID:  userID,
		Message: comment.Message,
	}

	photo, err = route.cuc.Update(c.Request.Context(), updatedComment, commentID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         photo.ID,
		"user_id":    photo.UserID,
		"title":      photo.Title,
		"photo_url":  photo.PhotoUrl,
		"caption":    photo.Caption,
		"updated_at": photo.UpdatedAt,
	})
}

// Delete godoc
// @Summary      Delete an comment
// @Description  delete an comment by ID
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Comment ID"
// @Success      200  {string}  string
// @Failure      400  {object}	ErrorResponse
// @Failure      401  {object}	ErrorResponse
// @Failure      404  {object}	ErrorResponse
// @Security     Bearer
// @Router       /comments/{id} [delete]
func (route *commentRoutes) Delete(c *gin.Context) {
	commentIDStr := c.Param("id")
	commentIDInt, err := strconv.Atoi(commentIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Failed to cast comment id to int",
		})
		return
	}

	commentID := uint(commentIDInt)

	err = route.cuc.Delete(c.Request.Context(), commentID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "Your comment has been successfully deleted",
		},
	)
}
