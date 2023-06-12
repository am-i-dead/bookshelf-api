package main

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAddBook(t *testing.T) {
	book := &Book{
		ID:     111,
		Name:   "Test Name",
		Author: "Test Author",
		Genre:  "Test Genre",
	}

	AddBook(book)

	clientOptions := options.Client().ApplyURI(MONGO)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		t.Errorf("Error connecting to database: %s", err)
	}
	defer client.Disconnect(context.TODO())

	coll := client.Database(DATABASE).Collection(COLLECTION)
	result := coll.FindOne(context.TODO(), bson.D{{"name", "Test Name"}})
	var testBook Book
	err = result.Decode(&testBook)
	if err != nil {
		t.Errorf("Error decoding inserted book: %s", err)
	}

	if testBook.ID != 111 {
		t.Errorf("Expected id to be 111 but got '%d'", testBook.ID)
	}
	if testBook.Name != "Test Name" {
		t.Errorf("Expected name to be 'Test Name' but got '%s'", testBook.Name)
	}
	if testBook.Author != "Test Author" {
		t.Errorf("Expected author to be 'Test Author' but got '%s'", testBook.Author)
	}
	if testBook.Genre != "Test Genre" {
		t.Errorf("Expected genre to be 'Test Genre' but got '%s'", testBook.Genre)
	}

	coll.DeleteOne(context.TODO(), bson.D{{"id", 111}})
}

func TestGetBook(t *testing.T) {
	book := &Book{
		ID:     111,
		Name:   "Test Name",
		Author: "Test Author",
		Genre:  "Test Genre",
	}

	clientOptions := options.Client().ApplyURI(MONGO)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		t.Errorf("Error connecting to database: %s", err)
	}
	defer client.Disconnect(context.TODO())

	coll := client.Database(DATABASE).Collection(COLLECTION)
	coll.InsertOne(context.TODO(), book)

	result := GetBook("name", "Test Name")
	if result.ID != book.ID {
		t.Errorf("Expected ID to be '%d' but got '%d'", book.ID, result.ID)
	}
	if result.Name != book.Name {
		t.Errorf("Expected name to be '%s' but got '%s'", book.Name, result.Name)
	}
	if result.Author != book.Author {
		t.Errorf("Expected author to be '%s' but got '%s'", book.Author, result.Author)
	}
	if result.Genre != book.Genre {
		t.Errorf("Expected genre to be '%s' but got '%s'", book.Genre, result.Genre)
	}

	coll.DeleteOne(context.TODO(), bson.D{{"id", 111}})
}

func TestEditBook(t *testing.T) {
	book := &Book{
		ID:     111,
		Name:   "Test Name",
		Author: "Test Author",
		Genre:  "Test Genre",
	}

	expected := &Book{
		ID:     111,
		Name:   "Edited Name",
		Author: "Test Author",
		Genre:  "Test Genre",
	}

	clientOptions := options.Client().ApplyURI(MONGO)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		t.Errorf("Error connecting to database: %s", err)
	}
	defer client.Disconnect(context.TODO())

	coll := client.Database(DATABASE).Collection(COLLECTION)
	coll.InsertOne(context.TODO(), book)

	EditBook("id", 111, "name", "Edited Name")

	result := GetBook("id", 111)
	if result.Name != expected.Name {
		t.Errorf("Expected name to be '%s' but got '%s'", expected.Name, result.Name)
	}

	coll.DeleteOne(context.TODO(), bson.D{{"id", 111}})
}

func TestDeleteBook(t *testing.T) {
	book := &Book{
		ID:     111,
		Name:   "Test Name",
		Author: "Test Author",
		Genre:  "Test Genre",
	}

	clientOptions := options.Client().ApplyURI(MONGO)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		t.Errorf("Error connecting to database: %s", err)
	}
	defer client.Disconnect(context.TODO())

	coll := client.Database(DATABASE).Collection(COLLECTION)
	coll.InsertOne(context.TODO(), book)

	DeleteBook("id", 111)

	result := GetBook("id", 111)
	if result.Name != "" {
		t.Errorf("Expected to be deleted but still exist")
	}
	if result.Author != "" {
		t.Errorf("Expected to be deleted but still exist")
	}
	if result.Genre != "" {
		t.Errorf("Expected to be deleted but still exist")
	}
}
