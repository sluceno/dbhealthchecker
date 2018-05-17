package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/sluceno/dbhealthchecker"
)

const (
	host     = "DB_HOST"
	port     = 5432
	dbname   = "DB_NAME"
	user     = "DB_USER"
	password = "DB_PASS"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	healthchecker := dbhealthchecker.New(
		db,
		dbhealthchecker.SetWaitTimeBetweenChecks(1*time.Second),
		dbhealthchecker.SetRunSuiteEvery(15*time.Second))

	checkProductsStock := dbhealthchecker.HealthCheck{
		Name:          "check-total-products-with-low-stock",
		Query:         "SELECT count(*) from products where stock < 5",
		ConditionType: dbhealthchecker.Lower,
		Thereshold:    1,
	}

	checkProductsPrice := dbhealthchecker.HealthCheck{
		Name:          "check-total-products-wrong-price",
		Query:         "SELECT count(*) from products where price < 5",
		ConditionType: dbhealthchecker.Lower,
		Thereshold:    1,
	}

	healthchecker.AddHealthChecker(checkProductsStock)
	healthchecker.AddHealthChecker(checkProductsPrice)

	healthchecks := healthchecker.Run()

	for healthcheck := range healthchecks {
		if healthcheck.Error() != nil {
			log.Println("healthcheck: ", healthcheck.Name, " healthy: ", "; result: ERROR error:", healthcheck.Error().Error())
		} else {
			log.Println("healthcheck: ", healthcheck.Name, " healthy: ", healthcheck.Healthy(), "; result: ", healthcheck.String())
		}
	}
}
