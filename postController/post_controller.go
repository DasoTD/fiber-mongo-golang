package postController

import (
	"context"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/models"
	"fiber-mongo-api/responses"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var postCollection *mongo.Collection = configs.GetCollection(configs.DB, "posts")
var validate = validator.New()

func CreatePost(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var post models.Post
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&post); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PostResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	newPost := models.Post{
		ID:     primitive.NewObjectID(),
		Title:  post.Title,
		Body:   post.Body,
		Author: post.Author,
	}
	result, err := postCollection.InsertOne(ctx, newPost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PostResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.PostResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetPost(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	postID := c.Params("postID")
	var post models.Post
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(postID)

	err := postCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&post)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PostResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.PostResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": post}})
}

func EditPost(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	postID := c.Params("postID")
	var post models.Post
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(postID)

	//validate the request body
	if err := c.BodyParser(&post); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PostResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&post); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PostResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"title": post.Title, "body": post.Body, "author": post.Author}

	result, err := postCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PostResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated post details
	var updatedPost models.Post
	if result.MatchedCount == 1 {
		err := postCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedPost)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.PostResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.PostResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedPost}})

}

func DeletePost(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	postID := c.Params("postID")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(postID)
	result, err := postCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PostResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.PostResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "User with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.PostResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}},
	)

}

func GetAllPost(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var posts []models.Post
	defer cancel()

	//objId, _ := primitive.ObjectIDFromHex(post)
	//const query = string{}
	//if c.Body.endDate

	results, err := postCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PostResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singlePost models.Post
		if err = results.Decode(&singlePost); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.PostResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		posts = append(posts, singlePost)

	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": posts}},
	)
}
