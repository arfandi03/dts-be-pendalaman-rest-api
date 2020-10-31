package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arfandidts/dts-be-pendalaman-microservice/auth-service/config"
	"github.com/arfandidts/dts-be-pendalaman-microservice/auth-service/database"
	"github.com/arfandidts/dts-be-pendalaman-microservice/auth-service/handler"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println(cfg)
	}

	db, err := initDB(cfg.Database)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("DB connection success")
	}

	authHandler := handler.AuthDB{DB: db}

	router := mux.NewRouter()

	router.Handle("/auth/validate", http.HandlerFunc(authHandler.ValidateAuth))
	router.Handle("/auth/signup", http.HandlerFunc(authHandler.SignUp))
	router.Handle("/auth/login", http.HandlerFunc(authHandler.Login))

	fmt.Printf("Auth service listen on :5001")
	log.Panic(http.ListenAndServe(":5001", router))
}

func getConfig() (config.Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	viper.SetConfigName("config.yml")

	if err := viper.ReadInConfig(); err != nil {
		return config.Config{}, err
	}

	var cfg config.Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return config.Config{}, err
	}

	return cfg, nil
}

func initDB(dbConfig config.Database) (*gorm.DB, error) {
	// root:password@tcp(localhost:3306)/dts_microservice_auth?charset=utf8&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DbName, dbConfig.Config)
	log.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&database.Auth{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
