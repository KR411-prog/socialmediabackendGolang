package main

import (
	"fmt"
	"log"

	"github.com/KR411-prog/socialmedia/internal/database"
)

func main() {
	c := database.NewClient("../internal/database/db.json")
	err := c.EnsureDB()
	if err != nil {
		log.Fatal(err)
	}

	 user, err := c.CreateUser("test@example.com", "password", "john doe", 18)
	 if err != nil {
		log.Fatal(err)
		}
	fmt.Println("user created", user)
	user, err = c.CreateUser("test1@example.com", "password", "john doe", 18)
	if err != nil {
	   log.Fatal(err)
	   }
   fmt.Println("user created", user)
   user, err = c.CreateUser("test2@example.com", "password", "john doe", 18)
   if err != nil {
	  log.Fatal(err)
	  }
  fmt.Println("user created", user)

	gotUser, err := c.GetUser("test@example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("user got", gotUser)

	// err = c.DeleteUser("test@example.com")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("user deleted")

	// _, err = c.GetUser("test@example.com")
	// if err == nil {
	// 	log.Fatal("shouldn't be able to get user that was deleted")
	// }
	// fmt.Println("user confirmed deleted")

	secondPost, err := c.CreatePost("test@example.com", "my cat is getting skinny now")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("another post created", secondPost)
	posts, err := c.GetPosts("test@example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("got posts", posts)


	posts, err = c.GetPosts("test@example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("got posts", posts)

	// err = c.DeletePost(secondPost.ID)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("deleted second post", posts)

	posts, err = c.GetPosts("test@example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("got posts", posts)

	err = c.DeleteUser("test@example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("user redeleted")
}
