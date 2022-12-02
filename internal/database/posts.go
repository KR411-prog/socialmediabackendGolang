package database

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Post
type Post struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UserEmail string    `json:"userEmail"`
	Text      string    `json:"text"`
}

func (c Client) CreatePost(userEmail, text string) (Post, error) {
	db,err := c.readDB()
	if err != nil {
		return Post{},err
	}
	if  _,ok := db.Users[userEmail]; !ok {
		return Post{},errors.New("user does not exist")
	}
	id := uuid.New().String()
	p := Post{
		ID: id,
		CreatedAt: time.Now().UTC(),
		UserEmail: userEmail,
		Text: text,
	}
	db.Posts[id] = p
    err = c.updateDB(db)
	if err != nil {
		return Post{},err
	}
	return p,nil
}

func (c Client) GetPosts(userEmail string) ([]Post, error) {
	db,err := c.readDB()
	if err != nil {
		return []Post{},err
	}
	posts := []Post{}

//	for i:=0;i<len(db.Posts);i++ {
	for _,post := range db.Posts {
      if post.UserEmail == userEmail {
		posts = append(posts,post)
	  }
	}
    return posts,nil
}

func (c Client) DeletePost(id string) error {
	dat,err := c.readDB()
	if err != nil {
		return err
	}
	_,ok := dat.Posts[id]
    if !ok {
		return errors.New("Post does not exist")
	}
	delete(dat.Posts,id)
	err = c.updateDB(dat)
	if err != nil {
		return err
	}
	return nil

}