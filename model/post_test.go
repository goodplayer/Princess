package model

import (
	"log"
	"math/rand"
	"strconv"
	"testing"
)

import ()

func ExamplePostSave() {
	initDb()

	user := NewUser()
	user.Id = 8
	err := user.GetUserById()
	if err != nil {
		log.Fatal(err)
	}

	post := NewPost()
	post.Title = "postTitle"
	post.Abstract = "postAbstract"
	post.Content = "postContent"
	post.PostUser = user

	err = post.Save()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(post)
	log.Println(post.PostUser)
}

func ExamplePostGetById() {
	initDb()

	post := NewPost()
	post.Id = 1

	err := post.GetPostById()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(post)
	log.Println(post.PostUser)
}

func TestPostSaveForTestData(t *testing.T) {
	initDb()

	user := NewUser()
	user.Id = 8
	err := user.GetUserById()
	if err != nil {
		t.Fatal(err)
	}

	l := rand.Intn(100)

	for i := 0; i < l; i++ {
		post := NewPost()
		post.Title = "postTitle" + strconv.Itoa(i)
		post.Abstract = "postAbstract" + strconv.Itoa(i)
		post.Content = "postContent" + strconv.Itoa(i)
		post.PostUser = user

		err = post.Save()
		if err != nil {
			t.Fatal(err)
		}

		log.Println(post)
		log.Println(post.PostUser)
	}

}
