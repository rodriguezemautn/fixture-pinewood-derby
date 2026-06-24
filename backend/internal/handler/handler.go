// Package handler define interfaces y helpers para handlers HTTP.
package handler

import "net/http"

// Handler define el contrato para un handler manejable por el router.
type Handler interface {
	Register(r Router)
}

// Router abstrae el registro de rutas HTTP.
type Router interface {
	Get(pattern string, h http.HandlerFunc)
	Post(pattern string, h http.HandlerFunc)
	Put(pattern string, h http.HandlerFunc)
	Delete(pattern string, h http.HandlerFunc)
}
