package database

import (
	"FOLLOWER/internal/config"
	"context"
	"log"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func Connect(cfg *config.Config) (neo4j.DriverWithContext, error) {
	var driver neo4j.DriverWithContext
	var err error

	for i := 0; i < 10; i++ {
		driver, err = neo4j.NewDriverWithContext(
			cfg.Neo4jURI,
			neo4j.BasicAuth(cfg.Neo4jUser, cfg.Neo4jPassword, ""),
		)

		if err == nil {
			err = driver.VerifyConnectivity(context.Background())
			if err == nil {
				log.Println("Connected to Neo4j successfully.")
				return driver, nil
			}
		}
		log.Printf("Database not ready, retrying in 2s... (Attempt %d/10)", i+1)
		time.Sleep(2 * time.Second)
	}

	return nil, err
}
