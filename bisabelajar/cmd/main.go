package main

import (
	"bisabelajar/api/v1/handler"
	"bisabelajar/api/v1/route"
	"bisabelajar/config"
	"bisabelajar/repository"
	"bisabelajar/service"
	"context"
	"encoding/json"
	"fmt"
	"internal/pkg/db"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()
	envFilePath := "../config/app.env"
	envMarshal, err := godotenv.Read(envFilePath)
	if err != nil {
		log.Fatal("can't read env")
	}

	marshalByte, err := json.Marshal(envMarshal)
	if err != nil {
		log.Fatal("can't marshal env")
	}

	var envStruct config.Env
	if err := json.Unmarshal(marshalByte, &envStruct); err != nil {
		log.Fatal("can't unmarshal env")
	}
	postgresUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		envStruct.PostgresUser, envStruct.PostgresPassword, envStruct.PostgresHost, envStruct.PostgresPort, envStruct.PostgresDB)
	redisAddr := fmt.Sprintf("%s:%s", envStruct.RedisHost, envStruct.RedisPort)
	sqlCon := db.NewPostgresConnection(ctx, postgresUrl)

	redisCon, err := db.NewRedisConnection(redisAddr, envStruct.RedisPassword)
	if err != nil {
		log.Fatal(err)
	}

	seriesRepo := repository.NewSeriesRepository(sqlCon, redisCon)
	seriesService := service.NewSeriesService(seriesRepo)
	seriesHandler := handler.NewSeriesHandler(seriesService)
	router := route.NewV1Route(seriesHandler, envStruct.Port)

	router.Intialize()
}
