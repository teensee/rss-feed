package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/patrickmn/go-cache"
	"log/slog"
	"net/http"
	"os"
	"rss-feed/internal/domain/rss/habr"
	"rss-feed/internal/handler"
	app_cache "rss-feed/pkg/cache"
	"time"
)

type AppHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type App struct {
	log      *slog.Logger
	router   *chi.Mux
	cache    app_cache.AppCache
	handlers map[string]AppHandler
}

type Builder struct {
	kernel *App
}

func NewBuilder() *Builder {
	return &Builder{
		kernel: &App{},
	}
}

func (b *Builder) WithLogger() *Builder {
	b.kernel.log = slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))

	return b
}

func (b *Builder) WithCache() *Builder {
	b.kernel.cache = app_cache.NewGoCache(cache.New(5*time.Minute, 30*time.Minute))

	return b
}

func (b *Builder) WithHandlers() *Builder {
	handlers := make(map[string]AppHandler)
	log := b.kernel.log
	appCache := b.kernel.cache

	habrHandler := habr.NewHabr(log, appCache)

	handlers["ping"] = handler.NewPingHandler(habrHandler)

	b.kernel.handlers = handlers
	return b
}

func (b *Builder) WithEndpoints() *Builder {
	if len(b.kernel.handlers) == 0 {
		panic("call WithHandlers first")
	}

	r := chi.NewRouter()

	r.Use(
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	)

	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", b.kernel.handlers["ping"].Handle)
	})

	b.kernel.router = r

	return b
}

func (b *Builder) Build() *App {
	return b.kernel
}

func (a *App) Run() {
	http.ListenAndServe(":3003", a.router)
}
