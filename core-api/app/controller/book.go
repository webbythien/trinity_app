package controller

import (
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	repo "github.com/hrshadhin/fiber-go-boilerplate/app/repository"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/validator"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
)

func GetBooks(c *fiber.Ctx) error {
	pageNo, pageSize := GetPagination(c)
	bookRepo := repo.NewBookRepo(database.GetDB())
	books, err := bookRepo.All(pageSize, uint(pageSize*(pageNo-1)))

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "books were not found",
		})
	}

	return c.JSON(fiber.Map{
		"page":      pageNo,
		"page_size": pageSize,
		"count":     len(books),
		"books":     books,
	})
}

func GetBook(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	bookRepo := repo.NewBookRepo(database.GetDB())
	book, err := bookRepo.Get(ID)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "book were not found",
		})
	}

	return c.JSON(fiber.Map{
		"book": book,
	})
}

func CreateBook(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"]
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "can't extract user info from request",
		})
	}

	// Create new Book struct
	book := &model.Book{}
	if err := c.BodyParser(book); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	book.ID = uuid.New()
	book.UserID = int(userID.(float64))
	book.Status = 1 // Active

	// Create a new validator for a Book model.
	validate := validator.NewValidator()
	if err := validate.Struct(book); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	bookRepo := repo.NewBookRepo(database.GetDB())
	if err := bookRepo.Create(book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"book": book,
	})
}

func UpdateBook(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	bookRepo := repo.NewBookRepo(database.GetDB())
	_, err = bookRepo.Get(ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "book were not found",
		})
	}

	book := &model.Book{}
	if err := c.BodyParser(book); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	book.ID = ID

	// Create a new validator for a Book model.
	validate := validator.NewValidator()
	if err := validate.Struct(book); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	if err := bookRepo.Update(ID, book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	dbBook, err := bookRepo.Get(ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"book": dbBook,
	})
}

func DeleteBook(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	bookRepo := repo.NewBookRepo(database.GetDB())
	_, err = bookRepo.Get(ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "book were not found",
		})
	}

	err = bookRepo.Delete(ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{})
}
