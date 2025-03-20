package main

import (
	"context"

	"github.com/sanchey92/jwt-example/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		panic(err)
	}

	a.Run()
}
