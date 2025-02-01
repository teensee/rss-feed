package app

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"rss-feed/internal/domain/rss"
	"rss-feed/internal/domain/rss/composite"
	"rss-feed/internal/domain/rss/habr"
	"rss-feed/internal/handler"
	"rss-feed/pkg/cache"
	http2 "rss-feed/pkg/http"
	"time"
)

type AppHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type App struct {
	log      *slog.Logger
	router   *chi.Mux
	cache    cache.AppCache
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
	b.kernel.cache = cache.NewGoCache(
		5*time.Minute,
		30*time.Minute,
		b.kernel.log,
	)

	return b
}

func (b *Builder) WithHandlers() *Builder {
	handlers := make(map[string]AppHandler)
	log := b.kernel.log
	appCache := b.kernel.cache
	if appCache == nil {
		log.InfoContext(context.TODO(), "Startup with dummy cache")
		appCache = cache.NewDummyCache()
	}

	habrHandler := habr.NewHabr(log, appCache)
	compositeHandler := composite.NewCompositeRss(
		composite.NewCompositeFetcher(http2.NewClient(url.URL{}, log), appCache, log),
		[]rss.Processor{
			rss.NewHtmlSanitizer(),

			// must be last in processor list
			rss.NewMaxLengthProcessor(),
		},
		log,
	)

	handlers["ping"] = handler.NewPingHandler(habrHandler, compositeHandler)
	handlers["feed"] = handler.NewFeedHandler(compositeHandler, log)

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
		r.Post("/feed", b.kernel.handlers["feed"].Handle)
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
