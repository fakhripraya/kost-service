package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/fakhripraya/kost-service/config"
	"github.com/fakhripraya/kost-service/data"
	"github.com/fakhripraya/kost-service/entities"
	"github.com/fakhripraya/kost-service/handlers"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"
	"github.com/srinathgs/mysqlstore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var err error

// Session Store based on MYSQL database
var sessionStore *mysqlstore.MySQLStore

// Adapter is an alias
type Adapter func(http.Handler) http.Handler

// Adapt takes Handler funcs and chains them to the main handler.
func Adapt(handler http.Handler, adapters ...Adapter) http.Handler {
	// The loop is reversed so the adapters/middleware gets executed in the same
	// order as provided in the array.
	for i := len(adapters); i > 0; i-- {
		handler = adapters[i-1](handler)
	}
	return handler
}

func main() {

	// creates a structured logger for logging the entire program
	logger := hclog.Default()

	// load configuration from env file
	err = godotenv.Load(".env")

	if err != nil {
		// log the fatal error if load env failed
		log.Fatal(err)
	}

	// Initialize app configuration
	var appConfig entities.Configuration
	err = data.ConfigInit(&appConfig)

	if err != nil {
		// log the fatal error if config init failed
		log.Fatal(err)
	}

	// initialize db session based on dialector
	logger.Info("Establishing database connection on " + appConfig.Database.Host + ":" + strconv.Itoa(appConfig.Database.Port))
	config.DB, err = gorm.Open(mysql.Open(config.DbURL(config.BuildDBConfig(&appConfig.Database))), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Open the database connection based on the initialized db session
	mySQLDB, err := config.DB.DB()
	if err != nil {
		log.Fatal(err)
	}

	defer mySQLDB.Close()

	// Creates a session store based on MYSQL database
	// If table doesn't exist, creates a new one
	logger.Info("Building session store based on " + appConfig.Database.Host + ":" + strconv.Itoa(appConfig.Database.Port))
	sessionStore, err = mysqlstore.NewMySQLStore(config.DbURL(config.BuildDBConfig(&appConfig.Database)), "dbMasterSession", "/", 3600*24*7, []byte(appConfig.MySQLStore.Secret))
	if err != nil {
		log.Fatal(err)
	}

	defer sessionStore.Close()

	// creates a kost instance
	kost := data.NewKost(logger)

	// creates the kost handler
	kostHandler := handlers.NewKostHandler(logger, kost, sessionStore)

	// creates a new serve mux
	serveMux := mux.NewRouter()

	// handlers for the API
	logger.Info("Setting handlers for the API")

	// get handlers
	getRequest := serveMux.Methods(http.MethodGet).Subrouter()
	getKostRequest := serveMux.Methods(http.MethodGet).Subrouter()

	// get specific kost handlers
	getKostRequest.HandleFunc("/{id:[0-9]+}", kostHandler.GetKost)
	getKostRequest.HandleFunc("/{id:[0-9]+}/picts", kostHandler.GetKostPicts)
	getKostRequest.HandleFunc("/{id:[0-9]+}/facilities", kostHandler.GetKostFacilities)
	getKostRequest.HandleFunc("/{id:[0-9]+}/facilities/room/{roomId:[0-9]+}", kostHandler.GetKostFacilities)
	getKostRequest.HandleFunc("/{id:[0-9]+}/benchmark", kostHandler.GetKostBenchmark)
	getKostRequest.HandleFunc("/{id:[0-9]+}/access", kostHandler.GetKostAccessibility)
	getKostRequest.HandleFunc("/{id:[0-9]+}/period", kostHandler.GetKostPeriod)
	getKostRequest.HandleFunc("/{id:[0-9]+}/around", kostHandler.GetKostAround)
	getKostRequest.HandleFunc("/{id:[0-9]+}/review", kostHandler.GetKostReviewList)
	getKostRequest.HandleFunc("/{id:[0-9]+}/owner", kostHandler.GetKostOwner)
	getKostRequest.HandleFunc("/{id:[0-9]+}/rooms", kostHandler.GetKostRoomList)
	getKostRequest.HandleFunc("/{id:[0-9]+}/rooms/{roomId:[0-9]+}/details", kostHandler.GetKostRoomInfo)
	getKostRequest.HandleFunc("/{id:[0-9]+}/rooms/all/details", kostHandler.GetKostRoomInfoAll)

	// get kost handlers
	getRequest.HandleFunc("/all/{category:[0-9]+}/{page:[0-9]+}", Adapt(
		http.HandlerFunc(kostHandler.GetKostList),
		kostHandler.MiddlewareParseUserRequest,
	).ServeHTTP)
	getRequest.HandleFunc("/all/near", Adapt(
		http.HandlerFunc(kostHandler.GetNearYouList),
		kostHandler.MiddlewareParseUserRequest,
	).ServeHTTP)
	getRequest.HandleFunc("/event/all", kostHandler.GetEventList)

	// get global middleware
	getRequest.Use(kostHandler.MiddlewareValidateAuth)
	getKostRequest.Use(kostHandler.MiddlewareParseKostGetRequest)

	// post handlers
	postRequest := serveMux.Methods(http.MethodPost).Subrouter()

	// post add new kost
	postRequest.HandleFunc("/add", kostHandler.AddKost)

	// post global middleware
	postRequest.Use(
		kostHandler.MiddlewareValidateAuth,
		kostHandler.MiddlewareParseKostPostRequest,
	)

	// CORS
	corsHandler := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// creates a new server
	server := http.Server{
		Addr:         appConfig.API.Host + ":" + strconv.Itoa(appConfig.API.Port), // configure the bind address
		Handler:      corsHandler(serveMux),                                       // set the default handler
		ErrorLog:     logger.StandardLogger(&hclog.StandardLoggerOptions{}),       // set the logger for the server
		ReadTimeout:  5 * time.Second,                                             // max time to read request from the client
		WriteTimeout: 10 * time.Second,                                            // max time to write response to the client
		IdleTimeout:  120 * time.Second,                                           // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		logger.Info("Starting server on port " + appConfig.API.Host + ":" + strconv.Itoa(appConfig.API.Port))

		err = server.ListenAndServe()
		if err != nil {

			if strings.Contains(err.Error(), "http: Server closed") == true {
				os.Exit(0)
			} else {
				logger.Error("Error starting server", "error", err.Error())
				os.Exit(1)
			}
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	signal.Notify(channel, os.Kill)

	// Block until a signal is received.
	sig := <-channel
	logger.Info("Got signal", "info", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}
