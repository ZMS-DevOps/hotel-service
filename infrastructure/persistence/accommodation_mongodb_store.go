package persistence

import (
	"context"
	"github.com/mmmajder/zms-devops-hotel-service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE   = "accommodationdb"
	COLLECTION = "accommodation"
)

type AccommodationMongoDBStore struct {
	accommodations *mongo.Collection
}

func NewAccommodationMongoDBStore(client *mongo.Client) domain.AccommodationStore {
	accommodations := client.Database(DATABASE).Collection(COLLECTION)
	return &AccommodationMongoDBStore{
		accommodations: accommodations,
	}
}

func (store *AccommodationMongoDBStore) Get(id primitive.ObjectID) (*domain.Accommodation, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *AccommodationMongoDBStore) GetAll() ([]*domain.Accommodation, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *AccommodationMongoDBStore) Insert(accommodation *domain.Accommodation) error {
	accommodation.Id = primitive.NewObjectID()
	result, err := store.accommodations.InsertOne(context.TODO(), accommodation)
	if err != nil {
		return err
	}
	accommodation.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (store *AccommodationMongoDBStore) DeleteAll() {
	store.accommodations.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *AccommodationMongoDBStore) Delete(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := store.accommodations.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (store *AccommodationMongoDBStore) Update(id primitive.ObjectID, accommodation *domain.Accommodation) error {
	filter := bson.M{"_id": id}

	//updateFields := bson.D{}
	//updateFields = append(updateFields, bson.E{"name", accommodation.Name})
	//updateFields = append(updateFields, bson.E{"location", accommodation.Location})
	//updateFields = append(updateFields, bson.E{"benefits", accommodation.Benefits})
	//updateFields = append(updateFields, bson.E{"photos", accommodation.Photos})
	//updateFields = append(updateFields, bson.E{"guest_number", accommodation.GuestNumber})
	//updateFields = append(updateFields, bson.E{"default_price", accommodation.DefaultPrice})
	updateFields := bson.D{
		{"name", accommodation.Name},
		{"location", accommodation.Location},
		{"benefits", accommodation.Benefits},
		{"photos", accommodation.Photos},
		{"guest_number", accommodation.GuestNumber},
		{"default_price", accommodation.DefaultPrice},
	}
	update := bson.D{{"$set", updateFields}}

	_, err := store.accommodations.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (store *AccommodationMongoDBStore) filter(filter interface{}) ([]*domain.Accommodation, error) {
	cursor, err := store.accommodations.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func (store *AccommodationMongoDBStore) filterOne(filter interface{}) (accommodation *domain.Accommodation, err error) {
	result := store.accommodations.FindOne(context.TODO(), filter)
	err = result.Decode(&accommodation)
	return
}

func decode(cursor *mongo.Cursor) (accommodations []*domain.Accommodation, err error) {
	for cursor.Next(context.TODO()) {
		var accommodation domain.Accommodation
		err = cursor.Decode(&accommodation)
		if err != nil {
			return
		}
		accommodations = append(accommodations, &accommodation)
	}
	err = cursor.Err()
	return
}
