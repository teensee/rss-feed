package kernel

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
	transport "rss-feed/internal/infrastructure/http"
	"rss-feed/internal/infrastructure/logger"
	"rss-feed/internal/infrastructure/processor"
	rest "rss-feed/internal/interfaces/rest/handler"
	mw "rss-feed/internal/interfaces/rest/middleware"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// AppHandler http-handler для обработки запросов
type AppHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

// kernel структура представляющая само ядро проекта
type kernel struct {
	log      logging.Logger
	router   *chi.Mux
	cache    domainCache.AppCache
	handlers map[string]AppHandler
}

// Builder строитель ядра проекта.
type Builder struct {
	kernel *kernel
}

// NewBuilder - конструктор билдера
func NewBuilder() *Builder {
	k := &kernel{
		cache:    cache.NewDummyCache(),
		handlers: make(map[string]AppHandler),
		log: logger.NewSlogAdapter(
			slog.New(
				slog.NewTextHandler(
					os.Stdout, &slog.HandlerOptions{
						AddSource:   true,
						Level:       slog.LevelDebug,
						ReplaceAttr: nil,
					},
				),
			),
		),
	}

	return &Builder{
		kernel: k,
	}
}

// WithPrettySlogTracingLogger Инициализирует pretty-slog логгер в проект
func (b *Builder) WithPrettySlogTracingLogger() *Builder {
	b.WithCustomLogger(
		logger.NewSlogAdapter(
			slog.New(
				logger.NewTraceIdSlogHandler(
					logger.NewSlogPrettyHandler(
						os.Stdout,
						&slog.HandlerOptions{
							Level: slog.LevelDebug,
						},
						true,
					),
				),
			),
		),
	)

	return b
}

// WithCustomLogger Устанавливает любой логгер
func (b *Builder) WithCustomLogger(logger logging.Logger) *Builder {
	b.kernel.log = logger

	return b
}

// WithGoCache Инициализирует in-memory cache
func (b *Builder) WithGoCache() *Builder {
	b.kernel.cache = cache.NewGoCache(
		5*time.Minute,
		30*time.Minute,
		b.kernel.log,
	)

	return b
}

// WithHandlers Инициализирует хендлеры и зависимости
func (b *Builder) WithHandlers() *Builder {
	l := b.kernel.log
	appCache := b.kernel.cache

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
		transport.NewFeedFetcher(
			transport.NewTracingClient(
				transport.NewClient(&url.URL{}, 30*time.Second, l),
				l,
			),
			appCache,
			domainCache.NewHashGenerator(hasher.NewSha256Hasher()),
			l,
		),
		registry,
		l,
	)

	b.kernel.handlers = map[string]AppHandler{
		"feed":       rest.NewFeedHandler(feedAggr, l),
		"processors": rest.NewProcessorListHandler(registry),
		"rawDoc":     rest.NewRawDocHandler(),
	}

	return b
}

// WithEndpoints инициализирует роутер
func (b *Builder) WithEndpoints() *Builder {
	if len(b.kernel.handlers) == 0 {
		panic("call WithHandlers first")
	}

	r := chi.NewRouter()

	if os.Getenv("APP_ENV") != "" {
		r.Use(middleware.BasicAuth(
			"feed-api",
			map[string]string{
				os.Getenv("BASIC_USER"): os.Getenv("BASIC_PASSWORD"),
			}),
		)
	}

	r.Use(
		mw.TraceId,

		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,

		middleware.Heartbeat("/ping"),
	)

	r.Route("/raw-doc", func(r chi.Router) {
		r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
			b.kernel.handlers["rawDoc"].Handle(writer, request)
		})
		r.Get("/{type}", b.kernel.handlers["rawDoc"].Handle)
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/feed", b.kernel.handlers["feed"].Handle)
		r.Get("/feed/processors", b.kernel.handlers["processors"].Handle)
	})

	b.kernel.router = r

	return b
}

// Build Возвращает готовый	образ ядра
func (b *Builder) Build() *kernel {
	return b.kernel
}

// Run запуск сервера
func (a *kernel) Run() {
	a.log.Info(context.Background(), "Startup server")

	err := http.ListenAndServe(":3003", a.router) // nolint:gosec // не использую таймаут?
	if err != nil {
		log.Fatal(err)
	}
}
