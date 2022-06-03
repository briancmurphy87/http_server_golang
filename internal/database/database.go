package database

import (
	"log"
	"errors"
	"encoding/json"
	"os"
	"time"
	"github.com/google/uuid"
) 

type Client struct {
	dbFilePath string
}

func NewClient(_dbFilePath string) Client {
	return Client{dbFilePath: _dbFilePath}
}

// db schema 
type databaseSchema struct {
	Users map[string]User `json:"users"`
	Posts map[string]Post `json:"posts"`
}

// User -
type User struct {
	CreatedAt time.Time `json:"createdAt"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
}

// Post -
type Post struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UserEmail string    `json:"userEmail"`
	Text      string    `json:"text"`
}

func (c Client) createDb() error {
	db := databaseSchema{
		Users: make(map[string]User), 
		Posts: make(map[string]Post), 
	}
	dat, _ := json.Marshal(db)
	err := os.WriteFile(c.dbFilePath, dat, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (c Client) EnsureDB() error {
	_, err := os.ReadFile(c.dbFilePath)
	if err != nil {
		// db file path DNE -> create db
		return c.createDb()
	}
	return nil
}

func (c Client) updateDB(db databaseSchema) error {
	dat, _ := json.Marshal(db)
	err := os.WriteFile(c.dbFilePath, dat, 066)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (c Client) readDB() (databaseSchema, error) {
	
	db := databaseSchema{}
	
	// read db file into json
	dat, err := os.ReadFile(c.dbFilePath)
	if err != nil {
		log.Fatal(err)
		return db, err
	}

	// decode json into struct = db schema
	err = json.Unmarshal(dat, &db)
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}
 
func (c Client) CreateUser(
	email, password, name string, 
	age int) (User, error) {

	// read current state of db
	db, err := c.readDB()
	if err != nil {
		log.Fatal(err)
		return User{}, err
	}
	
	// create a new User struct
	newUser := User{
		CreatedAt: time.Now().UTC(),
		Email: email, 
		Password: password,
		Name: name, 
		Age: age,
	}
	
	// add it to schema map
	if len(db.Users) == 0 {
		db.Users = make(map[string]User)
	}
	db.Users[newUser.Email] = newUser
	
	// update data on disk
	err = c.updateDB(db)
	return newUser, err
}


func (c Client) UpdateUser(email, password, name string, age int) (User, error) {
	// load db
	db, err := c.readDB()
	if err != nil {
		return User{}, err
	}
	
	// check if user exists in db
	user, ok := db.Users[email]
	if !ok {
		return user, errors.New("user doesn't exist")
	}

	// update fields of specified user
	user.Password = password
	user.Name = name
	user.Age = age
	
	// update data on disk
	err = c.updateDB(db)
	if err != nil {
		return user, err
	}

	// return updated user
	return user, nil
}

func (c Client) GetUser(email string) (User, error) {
		// load db
		db, err := c.readDB()
		if err != nil {
			return User{}, err
		}
		
		// check if user exists in db
		user, ok := db.Users[email]
		if !ok {
			return User{}, errors.New("user doesn't exist")
		}
		return user, nil 
}

func (c Client) DeleteUser(email string) error {
	// load db
	db, err := c.readDB()
	if err != nil {
		return err
	}
			
	// check if user exists in db
	_, ok := db.Users[email]
	if !ok {
		return errors.New("user doesn't exist")
	}
		
	// user exists, delete from db
	delete(db.Users, email)

	// update data on disk
	err = c.updateDB(db)

	return err
}


/*
Read the current database state
Make sure the user exists in the database using the userEmail
Create a new post struct. 
	Set the CreatedAt to now, 
	and create a new uuid using the google package: id := uuid.New().String()
Update the database on disk
*/ 
func (c Client) CreatePost(userEmail, text string) (Post, error) {
	// read current state of db
	db, err := c.readDB()
	if err != nil {
		return Post{}, err
	}

	// check that user exists in db
	_, ok := db.Users[userEmail]
	if !ok {
		return Post{}, errors.New("user doesn't exist")
	}

	// create a new post struct
	newPost := Post{
		ID: uuid.New().String(), 
		CreatedAt: time.Now().UTC(),
		UserEmail: userEmail,
		Text: text, 
	}

	// add new post to db schema
	if len(db.Posts) == 0 {
		db.Posts = make(map[string]Post)
	}
	db.Posts[newPost.ID] = newPost

	// update data on disk
	err = c.updateDB(db)
	return newPost, err
}

/**
Read the current database state
Iterate through all the posts 
	and add each post from the given user to a new slice
Return the matching posts
*/ 
func (c Client) GetPosts(userEmail string) ([]Post, error) {
	// read current state of db
	db, err := c.readDB()
	if err != nil {
		return []Post{}, err
	}
	// iterate through all posts in db 
	// -> match posts of specified user
	var matchingPosts []Post 
	for _, post := range db.Posts {
		if post.UserEmail == userEmail {
			matchingPosts = append(matchingPosts, post)
		}
	}

	// return posts matched to specified user
	return matchingPosts, nil 
}

func (c Client) DeletePost(id string) error {
	// load db
	db, err := c.readDB()
	if err != nil {
		return err
	}
	
	// check if post exists in db
	_, ok := db.Posts[id]
	if !ok {
		return errors.New("post doesn't exist")
	}
	
	// post exists, delete from db
	delete(db.Posts, id)
	
	// update data on disk
	err = c.updateDB(db)

	return err
}






