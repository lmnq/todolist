package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/lmnq/todolist/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoDBName    = "todolist"
	taskCollection = "tasks"
)

type ToDoListRepo struct {
	*mongo.Client
}

func New(mongoClient *mongo.Client) *ToDoListRepo {
	return &ToDoListRepo{mongoClient}
}

func (m *ToDoListRepo) CreateTask(ctx context.Context, task *entity.Task) (string, error) {
	collection := m.Database(mongoDBName).Collection(taskCollection)

	res, err := collection.InsertOne(ctx, task)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", fmt.Errorf("task already exists")
		}

		return "", fmt.Errorf("failed to create task: %w", err)
	}

	id := res.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}

func (m *ToDoListRepo) UpdateTask(ctx context.Context, id string, task *entity.Task) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	collection := m.Database(mongoDBName).Collection(taskCollection)

	res, err := collection.ReplaceOne(ctx, bson.M{"_id": objectID}, task)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("task not found. wrong id")
	}

	return nil
}

func (m *ToDoListRepo) DeleteTask(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	collection := m.Database(mongoDBName).Collection(taskCollection)

	res, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("task not found. wrong id")
	}

	return nil
}

func (m *ToDoListRepo) SetTaskStatusDone(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	collection := m.Database(mongoDBName).Collection(taskCollection)

	res, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": bson.M{"done": true}})
	if err != nil {
		return fmt.Errorf("failed to set task status: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("task not found. wrong id")
	}

	return nil
}

func (m *ToDoListRepo) GetTaskListByStatus(ctx context.Context, status string) ([]*entity.Task, error) {
	var filter bson.M
	now := time.Now().UTC()

	switch status {
	case "active":
		filter = bson.M{
			"$or": []bson.M{
				{"done": false},
				{"done": bson.M{"$exists": false}},
			},
			"activeAt": bson.M{"$lte": now},
		}
	case "done":
		filter = bson.M{
			"$or": []bson.M{
				{"done": true},
				{"done": bson.M{"$exists": true}},
			},
			"activeAt": bson.M{"$gt": now},
		}
	}

	sort := bson.D{{"activeAt", 1}}

	option := options.Find().SetSort(sort)

	collection := m.Database(mongoDBName).Collection(taskCollection)

	cursor, err := collection.Find(ctx, filter, option)
	if err != nil {
		return []*entity.Task{}, fmt.Errorf("failed to get task list: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []*entity.Task

	for cursor.Next(ctx) {
		var task *entity.Task

		if err := cursor.Decode(&task); err != nil {
			return []*entity.Task{}, fmt.Errorf("failed to decode task: %w", err)
		}

		tasks = append(tasks, task)

		if err := cursor.Err(); err != nil {
			return []*entity.Task{}, fmt.Errorf("failed to decode task: %w", err)
		}
	}

	if err := cursor.Err(); err != nil {
		return []*entity.Task{}, fmt.Errorf("failed to decode task: %w", err)
	}

	return tasks, nil
}
