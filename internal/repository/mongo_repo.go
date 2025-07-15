package repository

import (
	"context"
	"todocli/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepo struct {
	collection *mongo.Collection
}

func NewMongoRepo(collection *mongo.Collection) *MongoRepo {
	return &MongoRepo{collection: collection}
}

func (r *MongoRepo) Save(task *model.Task) error {
	_, err := r.collection.InsertOne(context.TODO(), task)
	return err
}

func (r *MongoRepo) Update(task *model.Task) error {
	filter := bson.M{"id": task.ID}
	update := bson.M{"$set": task}
	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *MongoRepo) GetAll() []*model.Task {
	cur, err := r.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil
	}
	defer cur.Close(context.TODO())

	var tasks []*model.Task
	for cur.Next(context.TODO()) {
		var t model.Task
		if err := cur.Decode(&t); err == nil {
			tasks = append(tasks, &t)
		}
	}
	return tasks
}

func (r *MongoRepo) FindByID(id int) *model.Task {
	var task model.Task
	err := r.collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&task)
	if err != nil {
		return nil
	}
	return &task
}

func (r *MongoRepo) Delete(id int) error {
	_, err := r.collection.DeleteOne(context.TODO(), bson.M{"id": id})
	return err
}
