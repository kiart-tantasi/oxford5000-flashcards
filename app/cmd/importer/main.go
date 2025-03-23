package main

import (
	"bufio"
	"database/sql"
	"flashcards/internal/env"
	"fmt"
	"log"
	"os"

	// https://stackoverflow.com/a/43810534
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	file := env.GetEnv("IMPORT_FILE", "../inputs/oxford_5000_sorted_exclude_3000.txt")
	importOxford(file)
}

func importOxford(file string) {

	// connect db
	db := connectDB()

	// read file
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error when opening file: %s", err)
		os.Exit(1)
	}

	// read line by line
	reader := bufio.NewReader(f)
	count := 0
	countError := 0
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			break
		}

		// iterate until all chars in a line is consumed
		for isPrefix {
			nextPart, nextIsPrefix, err := reader.ReadLine()
			if err != nil {
				break
			}
			line = append(line, nextPart...)
			isPrefix = nextIsPrefix
		}

		// Insert data into the flashcards table
		_, err = db.Exec("INSERT INTO flashcards_oxford (word) VALUES (?) ON DUPLICATE KEY UPDATE word = VALUES(word)", string(line))
		if err != nil {
			log.Fatalf("Error inserting word %s :%v", string(line), err)
			countError++
		} else {
			count++
		}
	}
	fmt.Println("finsihed importing")
	fmt.Printf("count: %d, countError: %d\n", count, countError)
}

func connectDB() *sql.DB {
	// Database connection details
	mysqlUser := env.GetEnv("MYSQL_USER", "root")
	mysqlPassword := env.GetEnv("MYSQL_PASSWORD", "changeme")
	mysqlDatabase := env.GetEnv("MYSQL_PASSWORD", "flashcards")
	mysqlHost := env.GetEnv("MYSQL_HOST", "localhost")
	mysqlPort := env.GetEnv("MYSQL_PORT", "3306")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)

	// Open a connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	fmt.Printf("Connected to MySQL (%s:%s) successfully!\n", mysqlHost, mysqlPort)
	return db
}
