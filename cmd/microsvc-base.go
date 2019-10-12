package main

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"

	stdjwt "github.com/dgrijalva/jwt-go"

	"github.com/hathbanger/microsvc-base/pkg/microsvc"
	"github.com/hathbanger/microsvc-base/pkg/microsvc/models"
)

var (
	errChannel = make(chan error)

	fconfig = flag.String(
		"config",
		"config.json",
		"the path to the service configuration file",
	)
	flogPath = flag.String(
		"log-path",
		os.Getenv("LOG_PATH"),
		"the absolute path to the log file",
	)
	faddress = flag.String(
		"address",
		os.Getenv("SERVICE_ADDRESS"),
		"the advertise address of the service (default env var \"SERVICE_ADDRESS\")",
	)
	fport = flag.String(
		"port",
		os.Getenv("SERVICE_PORT"),
		"the service advertise port (default env variable \"SERVICE_PORT\")",
	)
	fbindPort = flag.String(
		"bind-port",
		os.Getenv("PORT"),
		"the service port (default env variable \"PORT\")",
	)
	fdebug = flag.Bool(
		"debug", false, "print debug information for the service",
	)
	fversion = flag.Bool(
		"version", false, "set to true for version info",
	)
)

func main() {
	// parse flags
	flag.Parse()
	if *fversion {
		fmt.Printf(
			`{ "name": "%s", "version": "%s", "commit": "%s", "arch": "%s", "build_time": "%s", "api_version": "%s" }`,
			microsvc.Name,
			microsvc.Ver,
			microsvc.GitCommit,
			microsvc.Arch,
			microsvc.BuildTime,
			microsvc.APIVersion,
		)
		os.Exit(0)
	}

	// load configuration
	config, err := loadConfig()

	if err != nil {
		fmt.Println("WOO!!!", err)
		errChannel <- err
	}

	// set up logger
	var logger log.Logger
	{
		if len(*flogPath) <= 0 {
			*flogPath = fmt.Sprintf("/var/log/%s.log", microsvc.Name)
		}

		file, err := os.OpenFile(
			*flogPath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
		)
		if err != nil {
			fmt.Println("YO", err)

			errChannel <- err
		}

		logger = log.NewJSONLogger(io.MultiWriter(os.Stdout, file))
		logger = log.With(logger, "service", microsvc.Name)
		logger = log.With(logger, "timestamp", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)

		if !*fdebug {
			logger = level.NewFilter(logger, level.AllowError(), level.AllowInfo())
		}
	}

	var svc microsvc.Service
	{
		svc = microsvc.New(config, logger)
		svc = microsvc.InstrumentingMiddleware(
			kitprometheus.NewSummaryFrom(
				prometheus.SummaryOpts{
					Namespace: "microservices",
					Subsystem: "microsvc",
					Name:      "request_duration_seconds",
					Help:      "Request duration in seconds.",
				},
				[]string{"result", "mtype", "unit", "method"},
			),
			kitprometheus.NewCounterFrom(
				prometheus.CounterOpts{
					Namespace: "microservices",
					Subsystem: "microsvc",
					Name:      "request_count",
					Help:      "Number of requests received",
				},
				[]string{"result", "mtype", "unit", "method"},
			),
			svc,
		)
	}

	sd, reg, err := svc.ServiceDiscovery(config.ServiceAddr, config.ServicePort)
	if err != nil {

		errChannel <- err
	}

	logger.Log("registration", reg)
	if err != nil {

		errChannel <- err
	}

	// trap any unmanaged panaics and deregister the service from CONSUL
	defer func() {
		if r := recover(); r != nil {
			err = logger.Log(
				"trapped_panic", r,
				"stack_trace", debug.Stack(),
				"deregistering", reg.ID,
			)
			if err != nil {
				errChannel <- err
			}
			err := sd.Agent().ServiceDeregister(reg.ID)
			if err != nil {
				errChannel <- err
			}
		}
		logger.Log("deregistering", reg.ID)
		err = sd.Agent().ServiceDeregister(reg.ID)
		if err != nil {
			errChannel <- err
		}
	}()

	server := initializeServer(
		microsvc.MakeRoutes(svc, logger, config),
		*fbindPort,
		config,
	)

	go func() {
		err := sd.Agent().ServiceRegister(reg)
		if err != nil {
			errChannel <- err
		}
		logger.Log("transport", "HTTP", "port", config.ServicePort)
		errChannel <- server.ListenAndServe()
	}()

	// trap interrupts.
	go func() {
		c := make(chan os.Signal)
		signal.Notify(
			c,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGHUP,
			syscall.SIGQUIT,
		)
		errChannel <- fmt.Errorf("%s", <-c)
	}()
	logger.Log("exit", <-errChannel)
}

func initializeServer(router *mux.Router, port string, config *models.Config) http.Server {
	httpServerReadTimeout, _ := time.ParseDuration(config.HTTPServerReadTimeout)
	httpServerWriteTimeout, _ := time.ParseDuration(config.HTTPServerWriteTimeout)
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
		ReadTimeout:  httpServerReadTimeout,
		WriteTimeout: httpServerWriteTimeout,
	}
}
func loadConfig() (*models.Config, error) {
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
	if len(*flogPath) <= 0 {
		*flogPath = c.LogPath
	}
	if len(*faddress) > 0 {
		c.ServiceAddr = *faddress
	}
	if len(*fport) > 0 {
		c.ServicePort = *fport
	}
	if len(*fbindPort) <= 0 {
		return nil, errors.New("bind-port not set. please set $PORT or pass the bind-port flag")
	}

	var pubPEM = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCtnZ+HA5hmAcEW
SG/NrT9bddJHyF2w3SdEH/TYSvVl63Mxn9FcVkyL29AKdAUp65SmEOLpRK6hqoLq
8UdC3NXQ/X21OTcV05VAJuMm2KqKeSuEcsUIkhqpfLxNoIN9AIrGEMKNHJXUMNvf
QMr1hTIupZ6Jpy4AAlsKQ6D90kHp9UcN8lJ69gcmGejcvmI7+WlGv3HVbjulBcqK
UJdtFcWgZr0vFrihWjJgA0MPlknMgQf/jerMjGpQIKLZ2MxHdHngSc2NMma0dP4G
oh0GUsp5MIUg5hM/nD4v3iXM7tNPEQFj4V1o5fbLGSyPHqZYpw230ByxPw6x7Xn4
Gyl2crL1AgMBAAECggEAWn3Qu3+tPGXnrWSeGbcWUeaMbuvJobjzkXeSl/YiCDh7
tz7U0esNRMySmBA27M2kkhY1H160IwGL8UdHXFtceuzVS9MBmjfJEEH0nbfK1Bgq
DYQAnOICUZr5TwC96DaTHn932DMxCQNaZvgPkX8WU+fxRVBFEq4no6byT7n6ryVU
PCJsxs1kjXUJsCrE2MqZJCM19j+CNKa0JLoWAafPMb82znYMQ6f1q9Z4T+MZ1ajH
J+ljKa0WySUJexaFnTwSo/UXq3m+vSB9fDM87QPA/bqfhHnOp4RmT6pTfgC1iKUZ
dHVnUvi1Wao7KcnY00EyMbazmhg1fUbLUldQHcyE4QKBgQDZvHJRyqdyNP6Oab5g
X9brwCfriFTCf06K4Cy5ocaYNZowpuANsDNl5wikkHLFneJbQtJlhvUwtwTXfNQt
jbMuyj5zwieveAlPrheXGW/VjX4CmoVwv7HixenQscwKWfzS8Vk6GHCumqrLYyGs
OAMRHNXQ3pE8ZMgRK3rj239P3QKBgQDMIEeSH2UIQjyFtIUTvLzMwCsgGNpOrPHf
A5k+wC55Hv95/GhCkA2KTmzHl9JzrmqoYqVrTaryWtqtvLiXiGgxmPkPT1RLOH/W
Fz6qC9khrlBYyb3sc53mnwhAVEOPsp8eY+lN07bVMV8Any0vB4qyiCFTEGWBBXXd
SIDOn7hJ+QKBgQCAhyr8eSIK2pmBO45zmV9m3pEyCdHu1fNpKxd7pLF0W//exELy
EZblilGhwtrdKGvb7z//SoEl9oNXKIqfMUwaTKw87Nk8TSFB9cRbH1rStqkxpEEs
4xuAf8+br7iAS8pgQrOnBZJOn2I+mQ/hd1boHRtiJl+ZROyMphvusT0fyQKBgQCB
JldCK5zn2cijK/Ea6MpnnZprh23wY1nxKTy3SC7fMW6gxsNMggofHLmUmwlracpP
2YIh3xUum69KR2Jfdc2+u7OxLRb/NLMlSLW8LxzlQ33Qf2wsA4a/GJXG5cmNTI2C
U+KT/ETspH0gTpXu8I2foaO8A17FgCfvpuTgVovqOQKBgHVslwTHZzF1W86gghlQ
Pu/R6nmGCn3eoExdDTUb95UsaJARUI2Btah2kG4hb+n/tv7CgU72+ZLlkbGxPXnv
I4qNruA1uDyzzf1WRck72bRHezbJaiihQxPSg8gexS4KN2x6z0jp8U1D2Dv24+XI
brsn3pvJjLnMcQeLtDRe8t7z
-----END PRIVATE KEY-----`

	k, _ := stdjwt.ParseRSAPublicKeyFromPEM([]byte(pubPEM))
	if err != nil {
		fmt.Println("READIN", k)
		return nil, err
	}
	var tempKey = rsa.PublicKey{}
	c.PublicKey = &tempKey
	return c, nil
}
