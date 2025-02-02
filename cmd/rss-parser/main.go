package main

import app "rss-feed/internal/infrastructure/kernel"

func main() {
	kernel := app.NewBuilder().
		WithLogger().
		WithCache().
		WithHandlers().
		WithEndpoints().
		Build()

	kernel.Run()
}
