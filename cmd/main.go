// Package main contains application entry point
// Initializes and runs the application
package main

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
