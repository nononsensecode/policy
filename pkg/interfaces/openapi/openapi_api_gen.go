// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.2 DO NOT EDIT.
package openapi

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/go-chi/chi/v5"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /v1/policies)
	SavePolicy(w http.ResponseWriter, r *http.Request)

	// (GET /v1/policies/find/{match_string})
	FindPolicies(w http.ResponseWriter, r *http.Request, matchString string)

	// (GET /v1/policies/{policy_name})
	GetPolicy(w http.ResponseWriter, r *http.Request, policyName string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// SavePolicy operation middleware
func (siw *ServerInterfaceWrapper) SavePolicy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.SavePolicy(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// FindPolicies operation middleware
func (siw *ServerInterfaceWrapper) FindPolicies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "match_string" -------------
	var matchString string

	err = runtime.BindStyledParameter("simple", false, "match_string", chi.URLParam(r, "match_string"), &matchString)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter match_string: %s", err), http.StatusBadRequest)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.FindPolicies(w, r, matchString)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetPolicy operation middleware
func (siw *ServerInterfaceWrapper) GetPolicy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "policy_name" -------------
	var policyName string

	err = runtime.BindStyledParameter("simple", false, "policy_name", chi.URLParam(r, "policy_name"), &policyName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter policy_name: %s", err), http.StatusBadRequest)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetPolicy(w, r, policyName)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL     string
	BaseRouter  chi.Router
	Middlewares []MiddlewareFunc
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/policies", wrapper.SavePolicy)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1/policies/find/{match_string}", wrapper.FindPolicies)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1/policies/{policy_name}", wrapper.GetPolicy)
	})

	return r
}
