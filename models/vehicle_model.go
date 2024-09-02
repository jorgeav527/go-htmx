package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StringOrInt string

type VehicleModel struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Year      int                `json:"year,omitempty" validate:"required"`
	Make      string             `json:"make,omitempty" validate:"required"`
	Model     string             `json:"model,omitempty" validate:"required"`
	BodyStyle *string            `json:"bodyStyle,omitempty" bson:"bodyStyle,omitempty" validate:"omitempty"`
}
