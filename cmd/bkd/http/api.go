package http

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type API struct {
	ctx    context.Context
	srv    *http.Server
	logger *log.Logger
}

func NewAPI(
	ctx context.Context,
	logger *log.Logger,
	bookmarkController *BookmarkController,
) *API {
	handler := echo.New()

	// Middleware
	handler.Use(middleware.Logger())
	handler.Use(middleware.Recover())
	handler.Use(middleware.CORS())

	// Routes
	handler.GET("/bookmarks", bookmarkController.Index).Name = "bookmarks.index"
	handler.POST("/bookmarks", bookmarkController.Create).Name = "bookmarks.create"
	handler.DELETE("/bookmarks", bookmarkController.Delete).Name = "bookmarks.delete"

	return &API{
		ctx: ctx,
		srv: &http.Server{
			Addr:    ":8080",
			Handler: handler,
		},
		logger: logger,
	}
}

func (api *API) Execute() error {
	api.logger.Printf("Bookmark service is ready to receive web requests")

	err := api.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (api *API) Interrupt(err error) {
	api.logger.Printf("Stopping API service")

	e := api.srv.Shutdown(api.ctx)
	if e != nil {
		log.Println(err)
	}
}
