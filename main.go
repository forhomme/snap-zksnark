package main

import (
	"flag"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	r1cs2 "github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"smart-contract-service/app/handlers/web"
	"smart-contract-service/app/repo"
	"smart-contract-service/app/usecase"
	"smart-contract-service/configuration"
	"smart-contract-service/internal"
	"smart-contract-service/models"
	models2 "smart-contract-service/models/circuit"
	"time"
)

func configAndStartServer() {
	configMain := configuration.ServiceApp{
		EnvVariable: "DEV",
		Path:        setPath(),
	}
	path := setPath()
	fmt.Println(path)

	configMain.Load()

	htmlEcho := setWebRouter(configMain.Config)
	start(configMain, htmlEcho)
}

func start(config configuration.ServiceApp, htmlEcho *echo.Echo) {
	var err error

	if err = htmlEcho.Start(config.Config.ListenPort); err != nil {
		log.WithField("error", err).Error("Unable to start the server")
		os.Exit(1)
	}
}

func setPath() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	basePath := filepath.Dir(d)
	return basePath
}

func makeLogEntry(c echo.Context) *log.Entry {
	if c == nil {
		return log.WithFields(log.Fields{
			"at": time.Now().Format("2006-01-02 15:04:05"),
		})
	}
	return log.WithFields(log.Fields{
		"at":        time.Now().Format("2006-01-02 15:04:05"),
		"method":    c.Request().Method,
		"uri":       c.Request().URL.String(),
		"ip":        c.Request().RemoteAddr,
		"requestID": c.Response().Header().Get(echo.HeaderXRequestID),
	})
}

func middlewareLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("logger", makeLogEntry(c))
		makeLogEntry(c).Info("incoming request")
		return next(c)
	}
}

func setWebRouter(config configuration.ConfigApp) *echo.Echo {
	e := echo.New()
	e.Use(middleware.RequestID())
	e.Use(middlewareLogging)

	dbConn := configuration.InitSingleDB(config.PostgreConnection, config.LogMode)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisConnection,
		Password: "",
		DB:       0,
	})
	repoDb := repo.NewDatabaseConnection(dbConn)
	repoRedis := repo.NewRedisConnection(redisClient)

	uc := usecase.NewUsecase(repoRedis, repoDb, config)

	handler := web.NewHTTP(config, uc)
	var (
		migrate bool
		init    bool
	)
	flag.BoolVar(&migrate, "migrate", true, "If migrate true")
	flag.BoolVar(&init, "init", false, "set to true to run circuit Setup and export solidity Verifier")
	flag.Parse()

	if migrate {
		dbConn.AutoMigrate(
			&models.Partners{}, // create table partners
			&models.Customer{}, // create table customers
			&models.Payment{},  // create table payment
		)
	}

	if init {
		initElliptic()
		initHash()
		initEddsa()
	}

	web.NewRoutes(config).RegisterServices(e, handler)

	return e
}

func initHash() {
	var circuit models2.Circuit

	// compile circuit
	log.Println("compiling circuit")
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs2.NewBuilder, &circuit)
	assertNoError(err)

	// run groth16 trusted setup
	log.Println("running groth16.Setup")
	pk, vk, err := groth16.Setup(r1cs)
	assertNoError(err)

	// serialize R1CS, proving & verifying key
	log.Println("serialize R1CS (circuit)", internal.R1csPath)
	internal.Serialize(r1cs, internal.R1csPath)

	log.Println("serialize proving key", internal.PkPath)
	internal.Serialize(pk, internal.PkPath)

	log.Println("serialize verifying key", internal.VkPath)
	internal.Serialize(vk, internal.VkPath)
}

func initElliptic() {
	var elliptic models2.EllipticCurve

	// compile circuit
	log.Println("compiling elliptic circuit")
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs2.NewBuilder, &elliptic)
	assertNoError(err)

	// run groth16 trusted setup
	log.Println("running groth16.Setup")
	pk, vk, err := groth16.Setup(r1cs)
	assertNoError(err)

	// serialize R1CS, proving & verifying key
	log.Println("serialize R1CS (circuit)", internal.R1csEllipticPath)
	internal.Serialize(r1cs, internal.R1csEllipticPath)

	log.Println("serialize proving key", internal.PkEllipticPath)
	internal.Serialize(pk, internal.PkEllipticPath)

	log.Println("serialize verifying key", internal.VkEllipticPath)
	internal.Serialize(vk, internal.VkEllipticPath)
}

func initEddsa() {
	var eddsa models2.EddsaCircuit

	// compile circuit
	log.Println("compiling eddsa circuit")
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs2.NewBuilder, &eddsa)
	assertNoError(err)

	// run groth16 trusted setup
	log.Println("running groth16.Setup")
	pk, vk, err := groth16.Setup(r1cs)
	assertNoError(err)

	// serialize R1CS, proving & verifying key
	log.Println("serialize R1CS (circuit)", internal.R1csEddsaPath)
	internal.Serialize(r1cs, internal.R1csEddsaPath)

	log.Println("serialize proving key", internal.PkEddsaPath)
	internal.Serialize(pk, internal.PkEddsaPath)

	log.Println("serialize verifying key", internal.VkEddsaPath)
	internal.Serialize(vk, internal.VkEddsaPath)
}

func assertNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	configAndStartServer()
}
