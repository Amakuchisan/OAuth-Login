package main

import (
	route "github.com/Amakuchisan/OAuth-Login/route"
)

func main() {
	route.Echo.Logger.Fatal(route.Echo.Start(":1323"))
}
