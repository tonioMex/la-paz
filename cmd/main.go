package main

import (
	"context"
	"log"
	"pg-intro/app"
	"pg-intro/website"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// db, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/postgres")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// websiteRepository := website.NewPostgresRepository(db)

	// use pgx
	// dbPool, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer dbPool.Close()

	// websiteRepository := website.NewPGXRepository(dbPool)

	// use gorm
	db, err := gorm.Open(postgres.Open("postgres://postgres:postgres@localhost:5432/postgres"))
	if err != nil {
		log.Fatal(err)
	}

	websiteRepository := website.NewGormRepository(db)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.RunDemoRepository(ctx, websiteRepository)
}
