package main

import "github.com/nikolajbotalov/itk_academy_test/internal/app"

func main() {
	application, err := app.NewApp()
	if err != nil {
		panic(err)
	}
	defer application.Close()

	application.Server.Run()
}
