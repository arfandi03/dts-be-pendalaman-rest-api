package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arfandidts/dts-be-pendalaman-microservice/menu-service/config"
	"github.com/arfandidts/dts-be-pendalaman-microservice/menu-service/database"
	"github.com/arfandidts/dts-be-pendalaman-microservice/menu-service/handler"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Panic(err)
		return
	}

	db, err := initDB(cfg.Database)
	if err != nil {
		log.Panic(err)
		return
	}

	router := mux.NewRouter()

	authMiddleware := handler.AuthMiddleware{
		AuthService: cfg.AuthService,
	}

	menuHandler := handler.Menu{
		Db: db,
	}

	router.Handle("/add-menu", authMiddleware.ValidateAuth(http.HandlerFunc(menuHandler.AddMenu)))
	router.Handle("/menu", http.HandlerFunc(menuHandler.GetAllMenu))

	fmt.Println("Menu service listen on port :5000")
	log.Panic(http.ListenAndServe(":5000", router))
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DbName, dbConfig.Config)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Membuat tabel jika belum ada
	err = db.AutoMigrate(&database.Menu{})
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to DB")

	return db, nil
}
