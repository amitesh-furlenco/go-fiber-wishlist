package main

import (
	"log"
	"net/http"
	"os"

	"time"

	"github.com/amitesh-furlenco/go-fiber-wishlist/models"
	"github.com/amitesh-furlenco/go-fiber-wishlist/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

type Wishlist struct {
	ID                    uint           `gorm:"primary key;autoIncrement" json:"id"`
	USER_ID               *uint          `gorm:"index" json:"user_id"`
	CATALOG_ID            *uint          `gorm:"index" json:"catalog_id"`
	CATALOG_NAME          *string        `json:"catalog_name"`
	CATALOG_TYPE          *string        `json:"catalog_type"`
	CATALOG_IMAGE_URL     *string        `json:"catalog_image_url"`
	CATALOG_CONDITION     *string        `json:"catalog_condition"`
	CATALOG_STRIKE_PRICE  *float64       `json:"catalog_strike_price"`
	CATALOG_SELLING_PRICE *float64       `json:"catalog_selling_price"`
	CREATED_AT            time.Time      `json:"created_at"`
	DELETED_AT            gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (r *Repository) AddToWishlist(context *fiber.Ctx) error {
	wishlist := Wishlist{}

	err := context.BodyParser(&wishlist)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Request Failed"},
		)
		return err
	}

	err = r.DB.Create(&wishlist).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not add to wishlist"},
		)
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Added to wishlist successfully"},
	)
	return nil
}

func (r *Repository) GetWishlist(context *fiber.Ctx) error {
	wishlistModels := &[]models.Wishlists{}

	err := r.DB.Find(wishlistModels).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not get wishlist"},
		)
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "Wishlist fetched successfully",
			"data":    wishlistModels,
		},
	)
	return nil
}

func (r *Repository) RemoveFromWishlist(context *fiber.Ctx) error {
	wishlistModel := models.Wishlists{}
	id := context.Params("id")

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "ID cannot be empty",
			},
		)
		return nil
	}

	err := r.DB.Delete(&wishlistModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "Could not remove from wishlist",
			},
		)
		return err.Error
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "Removed from wishlist successfully",
		},
	)
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/add_to_wishlist", r.AddToWishlist)
	api.Delete("/remove_from_wishlist/:id", r.RemoveFromWishlist)
	api.Get("/get_wishlist", r.GetWishlist)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("Could not load the database")
	}

	err = models.MigrateWishlists(db)

	if err != nil {
		log.Fatal("Could not migrate database")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
