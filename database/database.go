package database

import (
	"github.com/gocql/gocql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"

	"orderworker/utils"
)

var (
	PgDB          *gorm.DB
	ScyllaSession *gocql.Session
)

func Connect() {
	var err error

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	PgDB, err = gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	utils.FailOnError(err, "Failed to connect to database")
	log.Println("Successfully connected to database PostgreSQL")

	scyllaHosts := os.Getenv("SCYLLA_HOSTS")
	if scyllaHosts == "" {
		log.Fatal("SCYLLA_HOSTS environment variable not set")
	}
	cluster := gocql.NewCluster(strings.Split(scyllaHosts, ",")...)
	cluster.Keyspace = "sabordarondonia"
	cluster.Consistency = gocql.Quorum
	ScyllaSession, err = cluster.CreateSession()
	utils.FailOnError(err, "Failed to create a session")
	log.Println("Successfully created a session ScyllaDB")
}
