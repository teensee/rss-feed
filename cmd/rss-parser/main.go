package main

import "rss-feed/internal/app"

func main() {
	kernel := app.NewBuilder().
		WithLogger().
		WithCache().
		WithHandlers().
		WithEndpoints().
		Build()

	kernel.Run()
}
