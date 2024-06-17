package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type ImageResponse struct {
	Id     primitive.ObjectID `json:"id"`
	Images []string           `json:"images"`
}
