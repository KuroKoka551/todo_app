package main

import (
	"flag"
	"os"
)

type config struct {
	Addr   string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

func newConfig() config {
	host := flag.String("host", "", "host")
	port := flag.String("port", "8080", "port")
	flag.Parse()

	dbHost, exists := os.LookupEnv("DB_HOST")
	if !exists {
		dbHost = "localhost"
	}
	dbPort, exists := os.LookupEnv("DB_PORT")
	if !exists {
		dbPort = "5432"
	}
	dbUser, exists := os.LookupEnv("DB_USER")
	if !exists {
		dbUser = "postgres"
	}
	dbPass, exists := os.LookupEnv("DB_PASS")
	if !exists {
		dbPass = "postgres"
	}
	dbName, exists := os.LookupEnv("DB_NAME")
	if !exists {
		dbName = "todos"
	}
	return config{
		Addr:   *host + ":" + *port,
		DBHost: dbHost,
		DBPort: dbPort,
		DBUser: dbUser,
		DBPass: dbPass,
		DBName: dbName,
	}
}
