package handlers

import (
	"fmt"
	"os/exec"

	"github.com/gofiber/fiber/v2"
	"endercloud-backend/internal/models"
)

func CreateServer(c *fiber.Ctx) error {
	var req models.ServerRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Example: Run Terraform to create EC2 instance
	cmd := exec.Command("terraform", "apply", "-auto-approve")
	cmd.Dir = "./internal/terraform"
	output, err := cmd.CombinedOutput()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Terraform failed",
			"details": string(output),
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Server created for world '%s'", req.WorldName),
		"seed":    req.Seed,
	})
}
