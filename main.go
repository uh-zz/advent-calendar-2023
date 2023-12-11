package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string
	Age  int64
}

// FYI: https://pkg.go.dev/golang.org/x/exp/slog
func Infof(logger *slog.Logger, format string, args ...any) {
	if !logger.Enabled(context.Background(), slog.LevelInfo) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelInfo, fmt.Sprintf(format, args...), pcs[0])
	_ = logger.Handler().Handle(context.Background(), r)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=postgres user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Tokyo",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&User{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})

	// Create
	db.Create(&User{Name: "L1212", Age: 100})

	// Read
	var user User
	db.First(&user, 1)
	Infof(logger, "find user with id: %+v", user)

	db.First(&user, "name = ?", "L1212")
	Infof(logger, "find user with name: %+v", user)

	// Update - update user's age to 101
	db.Model(&user).Update("Age", 101)
	Infof(logger, "udpated: %+v", user)
}
