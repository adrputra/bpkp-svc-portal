package app

import (
	"bpkp-svc-portal/app/config"
	"bpkp-svc-portal/app/connection"
	"bpkp-svc-portal/app/model"
	"bpkp-svc-portal/app/router"
	"bpkp-svc-portal/app/utils"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func Start() {
	config.InitConfig()
	cfg := config.GetConfig()

	utils.InitTimeLocation()

	tracer, closer, err := utils.InitJaeger(cfg)
	if err != nil {
		logrus.Fatalf("Failed to initialize Jaeger tracer: %v", err)
	}
	defer closer.Close()

	// Set global tracer
	opentracing.SetGlobalTracer(tracer)

	connection.InitConnection(*cfg)
	router.InitFactory(cfg, connection.Db, connection.Storage, connection.Redis, connection.Mq)

	host := cfg.Listener.Host
	port := cfg.Listener.Port

	e := echo.New()

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: utils.LogError,
	}))

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	auth := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(model.JwtCustomClaims)
		},
		SigningKey: []byte(cfg.Auth.AccessSecret),
	}

	public := e.Group("api")
	api := public.Group("/service")

	api.Use(echojwt.WithConfig(auth))
	api.Use(utils.IsAuthorized())

	e.Use(middleware.Logger())
	router.InitPublicRoute("", public)
	router.InitUserRoute("/user", api)
	router.InitRoleRoute("/role", api)
	router.InitParamRoute("/param", api)
	router.InitAttendanceRoute("/attendance", api)
	router.InitInstitutionRoute("/institution", api)

	e.Logger.Fatal(e.Start(host + ":" + strconv.Itoa(port)))
}
