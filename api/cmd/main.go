package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/igp-iw/config"
	"github.com/igp-iw/database"
	"github.com/igp-iw/notifications"
	"github.com/igp-iw/user"
	migrate "github.com/rubenv/sql-migrate"

	jwtware "github.com/gofiber/contrib/jwt"

	_ "github.com/lib/pq"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		panic(1)
	}

	app := fiber.New()
	validate := validator.New()
	db, err := sql.Open("postgres", config.DBSource)
	if err != nil {
		log.Panicf("Error opening connection to the database")
	}

	waitForDatabase(db)
	migrateDb(db)

	queries := database.New()
	kafkaNotifications, err := notifications.NewKafkaNotificationPublisher(config.KafkaBrokers, config.KafkaTopic)
	if err != nil {
		log.Println(err.Error())
		return
	}

	userService := user.NewUserService(queries, db, kafkaNotifications, config.JwtSecret)

	authHandler := user.NewAuthHandler(userService, validate)

	app.Post("/register", authHandler.Register())

	app.Get("verify/:email/:code", authHandler.Verify())
	app.Post("/login", authHandler.Login())

	protected := app.Group("/protected", jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.JwtSecret)},
	}))

	protected.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, from protected World!")
	})

	app.Listen(":3000")
}

func waitForDatabase(db *sql.DB) *sql.DB {
	timeout := 10
	retry := 0
	for {
		log.Println("Trying to establish connection to the database")
		err := db.Ping()
		if err == nil {
			log.Println("Established connection to the database")
			return db
		}
		if retry == timeout {
			log.Panicf("Error connecting to database after: %v retries, panicing", timeout)
		}
		retry++
		time.Sleep(2 * time.Second)
	}
}

func migrateDb(db *sql.DB) {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Panicln("Could not initial migration" + err.Error())
	}
	log.Printf("Applied %d migrations!\n", n)
}
