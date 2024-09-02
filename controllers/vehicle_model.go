package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/jorgeav527/vehicle-model/configs"
	"github.com/jorgeav527/vehicle-model/helpers"
	"github.com/jorgeav527/vehicle-model/models"
	"github.com/jorgeav527/vehicle-model/responses"
	views "github.com/jorgeav527/vehicle-model/views/vehicleModel"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var vehicleModelCollection *mongo.Collection = configs.GetCollection(configs.DB, "vehicleModel")

func CreateVehicleModel(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Extract and convert form fields
	year, err := strconv.Atoi(c.FormValue("year"))
	if err != nil || year == 0 {
		return c.JSON(http.StatusBadRequest, responses.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &echo.Map{"data": "Valid 'year' is required"},
		})
	}
	make := c.FormValue("make")
	model := c.FormValue("model")
	bodyStyle := c.FormValue("bodyStyle")

	// Custom validation
	if make == "" || model == "" || bodyStyle == "" {
		return c.JSON(http.StatusBadRequest, responses.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &echo.Map{"data": "All fields are required"},
		})
	}

	newVehicleModel := models.VehicleModel{
		Id:        primitive.NewObjectID(),
		Year:      year,
		Make:      make,
		Model:     model,
		BodyStyle: &bodyStyle,
	}

	// Insert the new vehicle model into the collection
	_, err = vehicleModelCollection.InsertOne(ctx, newVehicleModel)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.Response{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	return helpers.Render(c, views.NewVehicleRow(newVehicleModel))
}

func GetAVehicleModel(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	vehicleModelId := c.Param("vehicleModelId")
	var vehicleModel models.VehicleModel
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(vehicleModelId)

	response := vehicleModelCollection.FindOne(ctx, bson.M{"_id": objId})
	err := response.Decode(&vehicleModel)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	return helpers.Render(c, views.NewVehicleRow(vehicleModel))
}

func EditAVehicleModel(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	vehicleModelId := c.Param("vehicleModelId")
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(vehicleModelId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if c.Request().Method == http.MethodGet {
		var vehicleModel models.VehicleModel
		err := vehicleModelCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&vehicleModel)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
		}
		return helpers.Render(c, views.EditableVehicleRow(vehicleModel))
	}

	if c.Request().Method == http.MethodPost {

		// Extract and convert form fields
		year, err := strconv.Atoi(c.FormValue("year"))
		if err != nil || year == 0 {
			return c.JSON(http.StatusBadRequest, responses.Response{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data:    &echo.Map{"data": "Valid 'year' is required"},
			})
		}
		make := c.FormValue("make")
		model := c.FormValue("model")
		bodyStyle := c.FormValue("bodyStyle")

		// Custom validation
		if make == "" || model == "" || bodyStyle == "" {
			return c.JSON(http.StatusBadRequest, responses.Response{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data:    &echo.Map{"data": "All fields are required"},
			})
		}

		// Populate vehicleModel struct
		vehicleModel := models.VehicleModel{
			Year:      year,
			Make:      make,
			Model:     model,
			BodyStyle: &bodyStyle,
		}

		// Update in database
		update := bson.M{
			"year":      vehicleModel.Year,
			"make":      vehicleModel.Make,
			"model":     vehicleModel.Model,
			"bodyStyle": vehicleModel.BodyStyle,
		}

		result, err := vehicleModelCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
		}

		// Fetch updated vehicle model
		var updatedVehicleModel models.VehicleModel
		if result.MatchedCount == 1 {
			err := vehicleModelCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedVehicleModel)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
			}
		}
		return helpers.Render(c, views.NewVehicleRow(updatedVehicleModel))
	}
	return c.JSON(http.StatusMethodNotAllowed, responses.Response{Status: http.StatusMethodNotAllowed, Message: "error", Data: &echo.Map{"data": "Method not allowed"}})
}

func DeleteAVehicleModel(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	vehicleModelId := c.Param("vehicleModelId")
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(vehicleModelId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	result, err := vehicleModelCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.JSON(http.StatusNotFound, responses.Response{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": "VehicleModel with specified ID not found!"}})
	}

	return c.NoContent(http.StatusOK)
}

func GetAllVehicleModel(c echo.Context) error {
	// Get pagination parameters from query string
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}

	// Calculate the offset
	skip := int64((page - 1) * limit)

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var vehicleModels []models.VehicleModel

	options := options.Find().SetLimit(int64(limit)).SetSkip(skip)
	results, err := vehicleModelCollection.Find(ctx, bson.M{}, options)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	defer results.Close(ctx)

	// Decode all results into the vehicleModels slice
	for results.Next(ctx) {
		var singleVehicleModel models.VehicleModel
		err := results.Decode(&singleVehicleModel)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
		}
		// Handle nil BodyStyle
		if singleVehicleModel.BodyStyle == nil {
			na := "N/A"
			singleVehicleModel.BodyStyle = &na
		}
		vehicleModels = append(vehicleModels, singleVehicleModel)
	}
	// Check if there are more results to load
	hasMore := len(vehicleModels) == limit

	return helpers.Render(c, views.VehicleModelList(vehicleModels, hasMore, page, limit))
}

func SearchVehicleModel(c echo.Context) error {
	// Get the search query from form data
	searchQuery := c.FormValue("search")

	// Get pagination parameters from query string
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}

	// Calculate the offset
	skip := (page - 1) * limit

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create the aggregation pipeline
	pipeline := mongo.Pipeline{
		// Match stage: search in multiple fields using $regex for partial matching
		{
			{"$match", bson.M{
				"$or": []bson.M{
					{"year": bson.M{"$regex": searchQuery, "$options": "i"}},
					{"make": bson.M{"$regex": searchQuery, "$options": "i"}},
					{"model": bson.M{"$regex": searchQuery, "$options": "i"}},
					{"bodyStyle": bson.M{"$regex": searchQuery, "$options": "i"}},
				},
			}},
		},
		// Skip stage: skip the number of documents based on pagination
		{
			{"$skip", skip},
		},
		// Limit stage: limit the number of documents returned
		{
			{"$limit", limit},
		},
	}

	// Execute the aggregation pipeline
	cursor, err := vehicleModelCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.Response{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}
	defer cursor.Close(ctx)

	var vehicleModels []models.VehicleModel

	// Decode all results into the vehicleModels slice
	for cursor.Next(ctx) {
		var singleVehicleModel models.VehicleModel
		err := cursor.Decode(&singleVehicleModel)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    &echo.Map{"data": err.Error()},
			})
		}
		// Handle nil BodyStyle
		if singleVehicleModel.BodyStyle == nil {
			na := "N/A"
			singleVehicleModel.BodyStyle = &na
		}
		vehicleModels = append(vehicleModels, singleVehicleModel)
	}

	// Check if there are more results to load
	hasMore := len(vehicleModels) == limit

	// Return the search results as HTML (you can adjust this to return JSON if needed)
	return helpers.Render(c, views.VehicleModelList(vehicleModels, hasMore, page, limit))
}
