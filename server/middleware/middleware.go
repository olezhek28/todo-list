package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/olezhek28/todo-list/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DbUri            = "DB_URI"
	DbName           = "DB_NAME"
	DbCollectionName = "DB_COLLECTION_NAME"
)

var collection *mongo.Collection

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading the .env file")
	}
}

func createDBInstance(ctx context.Context) {
	connectionString := os.Getenv(DbUri)
	dbName := os.Getenv(DbName)
	collName := os.Getenv(DbCollectionName)

	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to mongodb!")

	collection = client.Database(dbName).Collection(collName)
	fmt.Println("collection instance created!")
}

func GetList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	payload, err := getList(r.Context())
	if err != nil {
		http.Error(w, "failed to get all tasks", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(payload)
}

func getList(ctx context.Context) ([]primitive.M, error) {
	cursor, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var res []primitive.M
	for cursor.Next(ctx) {
		var t bson.M
		errDecode := cursor.Decode(&t)
		if errDecode != nil {
			return nil, err
		}

		res = append(res, t)
	}

	err = cursor.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", http.MethodPost)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var task model.Task
	json.NewDecoder(r.Body).Decode(&task)
	err := create(r.Context(), task)
	if err != nil {
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func create(ctx context.Context, task model.Task) error {
	res, err := collection.InsertOne(ctx, task)
	if err != nil {
		return err
	}

	fmt.Println("Created a single task with id=", res.InsertedID)
	return nil
}

func Done(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", http.MethodPut)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	err := done(r.Context(), params["id"])
	if err != nil {
		http.Error(w, "failed to done task", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(params["id"])
}

func done(ctx context.Context, taskID string) error {
	id, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": true}}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Println("modified count:", res.ModifiedCount)
	return nil
}

func Undone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", http.MethodPut)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	err := undone(r.Context(), params["id"])
	if err != nil {
		http.Error(w, "failed to undone task", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(params["id"])
}

func undone(ctx context.Context, taskID string) error {
	id, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": false}}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Println("modified count:", res.ModifiedCount)
	return nil
}

func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", http.MethodDelete)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	err := delete(r.Context(), params["id"])
	if err != nil {
		http.Error(w, "failed to delete task", http.StatusInternalServerError)
		return
	}
}

func delete(ctx context.Context, taskID string) error {
	id, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}

	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	fmt.Println("deleted count:", res.DeletedCount)
	return nil
}

func DeleteAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	count, err := deleteAll(r.Context())
	if err != nil {
		http.Error(w, "failed to delete all tasks", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(count)
}

func deleteAll(ctx context.Context) (int64, error) {
	res, err := collection.DeleteMany(ctx, bson.D{{}})
	if err != nil {
		return 0, err
	}

	fmt.Println("deleted count:", res.DeletedCount)
	return res.DeletedCount, nil
}
