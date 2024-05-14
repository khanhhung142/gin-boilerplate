package main

import (
	"gin-boilerplate/cmd/server"
	_ "gin-boilerplate/docs"
)

//	@title			gin-boilerplate API
//	@version		1.0
//	@description	My boilerplate for personal projects

// @securityDefinitions.apikey BearerAuth
// @in							header
// @name						Authorization
// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	server.StartServer()
}
