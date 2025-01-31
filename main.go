package main

import (
	"habbit-tracker/cmd/server"
	_ "habbit-tracker/docs"
)

//	@title			gin-boilerplate API
//	@version		1.0
//	@description

// @securityDefinitions.apikey BearerAuth
// @in							header
// @name						Authorization
// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	server.StartServer()
}
