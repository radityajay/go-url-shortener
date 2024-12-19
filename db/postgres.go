package db

import (
	"fmt"
	"log"
	"os"

	"github.com/pandeptwidyaop/golog"
	"github.com/radityajay/go-url-shortener/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type PostgresStruct struct {
	DB *gorm.DB
}

var Postgres PostgresStruct

func NewPostgres() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Asia/Jakarta", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE"), os.Getenv("DB_PORT"))

	useResolver := os.Getenv("DB_HOST_READ_SOURCES")

	if Postgres.DB == nil {
		var db *gorm.DB
		var err error
		if os.Getenv("USE_SIMPLE_PROTOCOL") == "true" {
			db, err = gorm.Open(postgres.New(postgres.Config{
				DSN:                  dsn,
				PreferSimpleProtocol: true,
			}), &gorm.Config{SkipDefaultTransaction: true})

			if useResolver != "" {
				resolvers := []gorm.Dialector{}
				dsnSources := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Asia/Jakarta", os.Getenv("DB_HOST_READ_SOURCES"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE"), os.Getenv("DB_PORT"))

				for {
					resvDsn := os.Getenv(fmt.Sprintf("DB_HOST_READ_%d", (len(resolvers) + 1)))
					if resvDsn == "" {
						break
					}

					dsnRes := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Asia/Jakarta", resvDsn, os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE"), os.Getenv("DB_PORT"))
					resolvers = append(resolvers, postgres.Open(dsnRes))
				}

				db.Use(dbresolver.Register(dbresolver.Config{
					Sources:           []gorm.Dialector{postgres.Open(dsnSources)},
					Replicas:          resolvers,
					Policy:            dbresolver.RandomPolicy{},
					TraceResolverMode: true,
				}))

				fmt.Printf("Initialized %d replicas\n", len(resolvers))
			}

			if err != nil {
				log.Fatal(err)
			}
		} else {
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatal(err)
			}
		}

		Postgres = PostgresStruct{
			DB: db,
		}

		if os.Getenv("ENABLE_AUTO_MIGRATE") == "true" {
			RegisterTableToMigrate(db)
		}
	}

}

func RegisterTableToMigrate(db *gorm.DB) {
	e := db.AutoMigrate(
		&models.URL{},
	)
	if e != nil {
		golog.Slack.Error("Failed to run migration", e)
		log.Fatal(e)
	}
}
