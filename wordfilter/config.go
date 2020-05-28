package wordfilter

import (
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"sync"
)

type Config struct {
	Addr      string
	MysqlConn *gorm.DB
}

var once = &sync.Once{}

func NewConfig() *Config {
	addr := getEnv("ADDR", ":8080")
	return &Config{
		Addr:      addr,
		MysqlConn: getConn(),
	}
}
func getConn() *gorm.DB {
	var conn *gorm.DB
	once.Do(func() {
		mysql := getEnv("MYSQL", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
		db, err := gorm.Open("mysql", mysql)
		if err != nil {
			log.Fatal(err)
		}
		conn = db
	})
	return conn
}
func getEnv(name string, defaultValue string) string {
	if len(os.Getenv(name)) > 0 {
		return os.Getenv(name)
	}
	return defaultValue
}
