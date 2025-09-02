// Package main contains application entry point
// Initializes and runs the application
package main

import "currencyhub/internal/app"

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
