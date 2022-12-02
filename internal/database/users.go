package database

import (
	"errors"
	"time"
)

// User
type User struct {
	CreatedAt time.Time `json:"createdAt"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
}

func (c *Client) CreateUser(email, password, name string, age int) (User, error) {
	dat, err := c.readDB()
	if err != nil {
		return User{}, err
	}
	user := User{
		CreatedAt: time.Now().UTC(),
		Email:     email,
		Password:  password,
		Name:      name,
		Age:       age,
	}
	dat.Users[email] = user
	err = c.updateDB(dat)
	if err != nil {
		return User{}, err
	}
	return user, err
}

func (c *Client) GetUser(email string) (User, error) {
	//dat is struct of databaseschema
	dat, err := c.readDB()
	if err != nil {
		return User{}, err
	}
	user, ok := dat.Users[email]
	if !ok {
		return User{}, errors.New("user doesnt exit")
	}
	return user, nil

}

func (c Client) DeleteUser(email string) error {
	dat, err := c.readDB()
	if err != nil {
		return err
	}
	_, ok := dat.Users[email]
	if !ok {
		return errors.New("User does not exist")
	}
	delete(dat.Users, email)
	err = c.updateDB(dat)
	if err != nil {
		return err
	}
	return nil

}
