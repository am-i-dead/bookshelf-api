package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
	Genre  string `json:"genre"`
}

type FilterBody struct {
	KeyFilter string      `json:"keyFilter"`
	Filter    interface{} `json:"filter"`
}

type UpdateBody struct {
	KeyFilter string      `json:"keyFilter"`
	Filter    interface{} `json:"filter"`
	UpdateKey string      `json:"updateKey"`
	Update    interface{} `json:"update"`
}

func PrintJSON(data interface{}) (string, error) {
	value, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func AddBook(b *Book) {
	clientOptions := options.Client().ApplyURI(MONGO)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	client.Database(DATABASE).Collection(COLLECTION).InsertOne(context.TODO(), b)
}

func GetBook(keyFilter string, filter interface{}) *Book {
	clientOptions := options.Client().ApplyURI(MONGO)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	var b Book
	client.Database(DATABASE).Collection(COLLECTION).FindOne(context.TODO(), bson.D{{keyFilter, filter}}).Decode(&b)
	return &b
}

func EditBook(keyFilter string, filter interface{}, keyUpdate string, update interface{}) {
	clientOptions := options.Client().ApplyURI(MONGO)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	client.Database(DATABASE).Collection(COLLECTION).
		UpdateOne(context.TODO(), bson.D{{keyFilter, filter}}, bson.D{{"$set", bson.D{{keyUpdate, update}}}})
}

func DeleteBook(keyFilter string, filter interface{}) {
	clientOptions := options.Client().ApplyURI(MONGO)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	client.Database(DATABASE).Collection(COLLECTION).DeleteOne(context.TODO(), bson.D{{keyFilter, filter}})
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		fmt.Fprint(w, "Hello, server!")
		fmt.Println("Hello, server!")
	}
}

func BookHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/book/" {
		switch r.Method {
		case "POST":
			var b Book
			err := json.NewDecoder(r.Body).Decode(&b)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			AddBook(&b)
			w.WriteHeader(201)
			fmt.Println("New book was added!")
		case "GET":
			var fb FilterBody
			err := json.NewDecoder(r.Body).Decode(&fb)
			if err != nil {
				http.Error(w, "Can't decode JSON", http.StatusBadRequest)
				return
			}
			fmt.Println(fb)
			res, err := PrintJSON(GetBook(fb.KeyFilter, fb.Filter))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Fprint(w, res)
			fmt.Println("Book was sent!")
		case "PUT":
			var ub UpdateBody
			err := json.NewDecoder(r.Body).Decode(&ub)
			if err != nil {
				http.Error(w, "Can't decode JSON", http.StatusBadRequest)
				return
			}
			fmt.Println(ub)
			EditBook(ub.KeyFilter, ub.Filter, ub.UpdateKey, ub.Update)
			w.WriteHeader(201)
			fmt.Println("Some book was edited!")
		case "DELETE":
			var fb FilterBody
			err := json.NewDecoder(r.Body).Decode(&fb)
			if err != nil {
				http.Error(w, "Can't decode JSON", http.StatusBadRequest)
				return
			}
			fmt.Println(fb)
			DeleteBook(fb.KeyFilter, fb.Filter)
			w.WriteHeader(201)
			fmt.Println("Some book was deleted!")
		}
	}
}

const MONGO = "mongodb://localhost:27017/"
const DATABASE = "bookshelf"
const COLLECTION = "books"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HelloServer)
	mux.HandleFunc("/book/", BookHandler)

	fmt.Println("Server started!")

	log.Fatal(http.ListenAndServe("localhost:4040", mux))
}
