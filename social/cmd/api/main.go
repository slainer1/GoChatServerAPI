package main

import (
	"expvar"
	_ "main/docs"
	"main/internal/auth"
	"main/internal/db"
	"main/internal/env"
	"main/internal/mailer"
	"main/internal/ratelimiter"
	"main/internal/store"
	cache2 "main/internal/store/cache"
	"runtime"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////
// // To run main,
// // Use cmd 'air' in terminal to run while in 'social' dir
// // or start by using go run /cmd/api/*.go
// //////////////////////////////////////////////

const version = ""

//	@title			GopherSocial API
//	@description	API for GopherSocial, a social network for gophers
//	@termsOfService	http://swagger.io/terms/

//	@host	localhost:8080

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@desciription				JWT Authorization header using the Bearer scheme

func main() {
	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "localhost:4000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", ""),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisCfg: redisConfig{
			addr:     env.GetString("REDIS_ADDR", "localhost:6379"),
			password: env.GetString("REDIS_PASSWORD", ""),
			db:       env.GetInt("REDIS_DB", 0),
			enabled:  env.GetBool("REDIS_ENABLED", true),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3,
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				//email
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("BASIC_AUTH_USER", "admin"),
				pass: env.GetString("BASIC_AUTH_PASS", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "example"),
				exp:    time.Hour * 24 * 3, //3 days
				iss:    "gophersocial",
			},
		},
		rateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.GetInt("RATELIMITER_REQUESTS_COUNT", 20),
			TimeFrame:            time.Second * 5,
			Enabled:              env.GetBool("RATE_LIMITER_ENABLED", true),
		},
	}
	//Logger
	//flushes buffered log entries
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	//Database
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("DB Conntection Pool Established")

	// Redis Cache
	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache2.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.password, cfg.redisCfg.db)
		logger.Info("Redis Conntection Pool Established")
		defer rdb.Close()
	}
	//rate limiter
	ratelimiter := ratelimiter.NewFixedWindowRateLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	store := store.NewStorage(db)
	cacheStorage := cache2.NewRedisStorage(rdb)
	mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)

	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)

	app := &application{
		config:        cfg,
		store:         store,
		cacheStorage:  cacheStorage,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
		rateLimiter:   ratelimiter,
	}

	//metrics collected
	//version -- adds to the /v1/debug/vas
	expvar.NewString("version").Set(version)
	//declares a named export variable
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()
	logger.Fatal(app.run(mux))

}

// Package Validator
//func main() {
//	post := Post{}
//	t := reflect.TypeOf(post)
//	field, _ := t.FieldByName("Post")
//	validateTag := field.Tag.Get("validate")
//	validateRules := strings.Split(validateTag, ",")
//	fmt.Println("Validate Tags", validateTag)
//	fmt.Println("Validate Rules", validateRules)
//}
