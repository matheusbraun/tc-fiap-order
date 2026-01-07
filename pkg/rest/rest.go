package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Controller is an interface that defines a method for registering routes with a chi router.
// It is used to create a consistent way of defining and registering routes in the application.
type Controller interface {
	RegisterRoutes(r chi.Router)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
