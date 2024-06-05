package route

import (
	"bisabelajar/api/v1/handler"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

type Route struct {
	seriesHandle handler.SeriesHandler
	port         string
}

func NewV1Route(seriesHandler handler.SeriesHandler, port string) Route {
	return Route{seriesHandle: seriesHandler, port: port}
}

func (b *Route) Intialize() {
	router := chi.NewRouter()

	router.Use(
		render.SetContentType(render.ContentTypeJSON), //forces Content-type
		middleware.RedirectSlashes,
		middleware.Recoverer,
		middleware.Logger, //middleware to recover from panics
		cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"https://*", "http://*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}),
	)
	//Sets context for all requests
	router.Use(middleware.Timeout(30 * time.Second))

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/series", b.seriesHandle.Routes())
	})
	log.Printf("runnning on port %s", b.port)
	log.Fatal(http.ListenAndServe(b.port, router))
}
