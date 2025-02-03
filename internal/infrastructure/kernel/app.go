package app

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"rss-feed/internal/application"
	domainCache "rss-feed/internal/domain/cache"
	"rss-feed/internal/domain/logging"
	"rss-feed/internal/domain/rss"
	"rss-feed/internal/infrastructure/cache"
	"rss-feed/internal/infrastructure/cache/hasher"
	http2 "rss-feed/internal/infrastructure/http"
	"rss-feed/internal/infrastructure/logger"
	"rss-feed/internal/infrastructure/processor"
	rest "rss-feed/internal/interface/rest/handler"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type AppHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type Kernel struct {
	log      logging.Logger
	router   *chi.Mux
	cache    domainCache.AppCache
	handlers map[string]AppHandler
}

type Builder struct {
	kernel *Kernel
}

func NewBuilder() *Builder {
	return &Builder{
		kernel: &Kernel{},
	}
}

func (b *Builder) WithLogger() *Builder {
	b.kernel.log = logger.NewSlogAdapter(
		slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				}),
		),
	)

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
	l := b.kernel.log
	appCache := b.kernel.cache
	if appCache == nil {
		l.Info(context.TODO(), "Startup with dummy cache")
		appCache = cache.NewDummyCache()
	}

	processors := []rss.Processor{
		processor.NewHtmlSanitizer(),

		processor.NewSizeOfProcessor(8),
		processor.NewSizeOfProcessor(16),
		processor.NewSizeOfProcessor(24),
		processor.NewSizeOfProcessor(32),

		processor.NewMaxLengthProcessor(),
	}
	processors = append(processors, processor.NewTimeFormatters()...)

	registry := processor.NewProcessorRegistry(processors)

	feedAggr := application.NewFeedService(
		http2.NewFeedFetcher(http2.NewClient(&url.URL{}, l), appCache, domainCache.NewHashGenerator(hasher.NewSha256Hasher()), l),
		registry,
		l,
	)

	handlers["feed"] = rest.NewFeedHandler(feedAggr, l)
	handlers["processors"] = rest.NewProcessorListHandler(registry)

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
		r.Post("/feed", b.kernel.handlers["feed"].Handle)
		r.Get("/filter/processor", b.kernel.handlers["processors"].Handle)
	})

	b.kernel.router = r

	return b
}

func (b *Builder) Build() *Kernel {
	return b.kernel
}

func (a *Kernel) Run() {
	err := http.ListenAndServe(":3003", a.router) // nolint:gosec // не использую таймаут?
	if err != nil {
		log.Fatal(err)
	}
}
