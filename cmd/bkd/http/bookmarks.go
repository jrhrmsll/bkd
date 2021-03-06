package http

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"bookmarks/internal"
)

type BookmarksController struct {
	ctx                context.Context
	bookmarkRepository internal.BookmarksRepository
	logger             *log.Logger
}

func NewBookmarkController(
	ctx context.Context,
	bookmarkRepository internal.BookmarksRepository,
	logger *log.Logger) *BookmarksController {
	return &BookmarksController{
		ctx:                ctx,
		bookmarkRepository: bookmarkRepository,
		logger:             logger,
	}
}

func (controller *BookmarksController) Index(c echo.Context) error {
	bookmarks, err := controller.bookmarkRepository.Find()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, bookmarks)
}

func (controller *BookmarksController) Create(c echo.Context) error {
	var bookmark = new(internal.Bookmark)

	if err := c.Bind(bookmark); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	_, err := controller.bookmarkRepository.Store(bookmark)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (controller *BookmarksController) Delete(c echo.Context) error {
	var version = new(struct {
		Version string `json:"version"`
	})

	if err := c.Bind(version); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	_, err := controller.bookmarkRepository.Delete(version.Version)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
