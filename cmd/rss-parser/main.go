package main

import app "rss-feed/internal/infrastructure/kernel"

func main() {
	kernel := app.NewBuilder().
		WithPrettySlogTracingLogger().
		WithGoCache().
		WithHandlers().
		WithEndpoints().
		Build()

	kernel.Run()
}
