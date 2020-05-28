package wordfilter

import (
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

type Config struct {
		Addr      string
		MysqlConn *gorm.DB
	}

func NewConfig() *Config {
	addr := getEnv("ADDR", ":8080")
	mysql := getEnv("MYSQL", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open("mysql", mysql)
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		Addr:      addr,
		MysqlConn: db,
	}
}
func getEnv(name string, defaultValue string) string {
	if len(os.Getenv(name)) > 0 {
		return os.Getenv(name)
	}
	return defaultValue
}
