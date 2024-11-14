package migrator

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	var storagePath, migrationsPath string
	storagePath = os.Getenv("STORAGE_PATH")
	migrationsPath = os.Getenv("MIGRATIONS_PATH")

	if storagePath == "" {
		panic("storagePath is not required")
	}
	if migrationsPath == "" {
		panic("migartionsPath is not required")
	}
	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s", storagePath),
	)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
		}
		panic(err)
	}
	fmt.Println("migrations are applied")
}
