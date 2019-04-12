package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/hathbanger/microsvc-base/pkg/microsvc"
	"github.com/hathbanger/microsvc-base/pkg/microsvc/models"
)

const (
	name = "microsvc-base"
)

var (
	logger     log.Logger
	errChannel = make(chan error)

	fconfig = flag.String(
		"config",
		"config.json",
		"the path to the service configuration file",
	)
	faddress = flag.String(
		"addr",
		os.Getenv("SERVICE_ADDRESS"),
		"the address of the service (default env var \"SERVICE_ADDRESS\")",
	)
	fport = flag.String(
		"port",
		os.Getenv("PORT"),
		"the service port (default env var \"PORT\")",
	)
	fdebug = flag.Bool(
		"debug", false, "print debug information for the service",
	)
)

func main() {
	flag.Parse()

	logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "service", name)
	logger = log.With(logger, "timestamp", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	config, err := loadConfig()
	if err != nil {
		fmt.Println("error", err)
		errChannel <- err
	}

	svc := microsvc.New(config, logger)

	// metrics:
	// svc = foo.InstrumentingMiddleware(
	// 	kitprometheus.NewSummaryFrom(
	// 		prometheus.SummaryOpts{
	// 			Namespace: "microservices",
	// 			Subsystem: "microsvc_base",
	// 			Name:      "request_duration_seconds",
	// 			Help:      "Request duration in seconds",
	// 		},
	// 		[]string{"result", "mtype", "unit", "method"},
	// 	),
	// 	kitprometheus.NewCounterFrom(
	// 		prometheus.CounterOpts{
	// 			Namespace: "microservices",
	// 			Subsystem: "microsvc_base",
	// 			Name:      "request_count",
	// 			Help:      "Number of requests received.",
	// 		},
	// 		[]string{"result", "mtype", "unit", "method"},
	// 	),
	// 	svc,
	// )

	// sd, reg, err := svc.ServiceDiscovery()
	// if err != nil {
	// 	errChannel <- err
	// }

	// logger.Log("registration", reg)

	// trap any panics and deregister
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		logger.Log(
	// 			"trapped_panic", r,
	// 			"stack_trace", debug.Stack(),
	// 			"deregistering", reg.ID,
	// 		)
	// 		err := sd.Agent().ServiceDeregister(reg.ID)
	// 		if err != nil {
	// 			logger.Log("err", err)
	// 			errChannel <- err
	// 		}
	// 	}
	// 	logger.Log("degistering", reg.ID)
	// 	err := sd.Agent().ServiceDeregister(reg.ID)
	// 	if err != nil {
	// 		errChannel <- err
	// 	}
	// }()

	server := initializeServer(
		microsvc.MakeRoutes(svc, logger, initializeMiddleware(config)),
		config.ServicePort,
	)

	go func() {
		// err := sd.Agent().ServiceRegister(reg)
		// if err != nil {
		// 	errChannel <- err
		// }
		logger.Log("transport", "HTTP", "port", config.ServicePort)
		errChannel <- server.ListenAndServe()
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(
			c,
			syscall.SIGINT,
			syscall.SIGKILL,
			syscall.SIGTERM,
			syscall.SIGHUP,
			syscall.SIGQUIT,
		)
		errChannel <- fmt.Errorf("%s", <-c)
	}()
	logger.Log("exit", <-errChannel)
}

func initializeMiddleware(config *models.Config) []endpoint.Middleware {
	return []endpoint.Middleware{
		microsvc.NewAuth(
			log.With(logger, "event", "authorization"),
			config,
		).ValidateAuth(),
	}
}

func initializeServer(router *mux.Router, port string) http.Server {
	return http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: handlers.CORS(
			handlers.AllowedOrigins([]string{
				"*",
			}),
			handlers.AllowedMethods([]string{
				"GET", "HEAD", "POST", "PUT", "OPTIONS",
			}),
			handlers.AllowedHeaders([]string{
				"Origin", "X-Requested-With", "Content-Type", "Authorization",
			}),
		)(router),
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
}

func loadConfig() (*models.Config, error) {
	fmt.Println("loading")
	var c *models.Config
	_, err := os.Stat(*fconfig)
	if err != nil {
		return nil, err
	}
	f, err := ioutil.ReadFile(*fconfig)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(f, &c)
	if err != nil {
		return nil, err
	}
	logger.Log("loading", "config")
	logger.Log("addr", *faddress)
	if len(*fport) > 0 {
		c.ServicePort = *fport
		logger.Log("service_port", c.ServicePort)
	}
	return c, nil
}
