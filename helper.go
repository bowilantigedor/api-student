package main

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ok(c *fiber.Ctx, message string, data interface{}, meta *Meta) error {
	return c.Status(fiber.StatusOK).JSON(WebResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func created(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(WebResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func noContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func fail(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(WebResponse{
		Success: false,
		Message: message,
	})
}

func failValidation(c *fiber.Ctx, errors interface{}) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(WebResponse{
		Success: false,
		Message: "Validation Error",
		Errors:  errors,
	})
}

type ListParams struct {
	Search string
	Major  string
	SortBy string
	Order  string
	Page   int
	Limit  int
}

func parseListQuery(c *fiber.Ctx) ListParams {
	search := strings.TrimSpace(c.Query("search", ""))
	major := strings.TrimSpace(c.Query("major", ""))

	sortByInput := strings.ToLower(strings.TrimSpace(c.Query("sort_by", "id")))
	orderInput := strings.ToUpper(strings.TrimSpace(c.Query("order", "ASC")))

	// Whitelist kolom pengurutan untuk keamanan
	allowedSortCols := map[string]string{
		"id":         "id",
		"nim":        "nim",
		"name":       "name",
		"gpa":        "gpa",
		"created_at": "created_at",
	}

	sortBy, exists := allowedSortCols[sortByInput]
	if !exists {
		sortBy = "id"
	}

	order := "ASC"
	if orderInput == "DESC" {
		order = "DESC"
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return ListParams{
		Search: search,
		Major:  major,
		SortBy: sortBy,
		Order:  order,
		Page:   page,
		Limit:  limit,
	}
}
