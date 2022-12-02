package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type Client struct {
	FilePath string
}

type databaseSchema struct {
	Users map[string]User `json:"users"`
	Posts map[string]Post `json:"posts"`
}

func NewClient(path string) *Client {
	return &Client{
		FilePath: path,
	}
}

func (c *Client) createDB() error {
	// convert struct to json
	dat, err := json.Marshal(databaseSchema{
		Users: make(map[string]User),
		Posts: make(map[string]Post),
	})
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(c.FilePath, dat, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (c *Client) EnsureDB() error {
	_, err := os.ReadFile(c.FilePath)
	if errors.Is(err, os.ErrNotExist) {
		c.createDB()
	}
	return err
}

func (c *Client) updateDB(db databaseSchema) error {
	dat, err := json.Marshal(db)

	if err != nil {
		return err
	}
	err = os.WriteFile(c.FilePath, dat, 0666)
	if err != nil {
		return err
	}
	return nil
}
func (c *Client) readDB() (databaseSchema, error) {
	dat, err := os.ReadFile(c.FilePath)
	if err != nil {
		log.Fatal(err)
	}
	db := databaseSchema{}
	err = json.Unmarshal(dat, &db)
	return db, err
}
