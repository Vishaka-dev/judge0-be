package handlers

import (
	"net/http"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/logger"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/repositories"
	"github.com/gin-gonic/gin"
)

func GetDSASubmissionResultsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")

	submissions, currentPage, totalPages, err := repositories.GetDSASubmissionResults(ctx, page, pageSize)
	if err != nil {
		logger.Log.Error("GetDSASubmissionResultsHandler: failed to get submission results", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"submissions": submissions,
		"currentPage": currentPage,
		"totalPages":  totalPages,
	})
}
