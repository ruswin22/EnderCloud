package main

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
)


type ServerRequest struct {
	ServerName string `json:"serverName"`
	Seed       string `json:"seed"`
	Version    string `json:"version"`
	RAM        string `json:"ram"`
}

type TerraformOutputs struct {
	PublicIp struct {
		Value string `json:"value"`
	} `json:"public_ip"`
	PrivateKeyPath struct {
		Value string `json:"value"`
	} `json:"private_key_path"`
}

func ramToInstance(ram string) string {
	switch ram {
	case "1G":
		return "t2.micro"
	case "2G":
		return "t2.small"
	case "4G":
		return "t2.medium"
	default:
		return "t2.micro"
	}
}

func runCmd(ctx context.Context, dir string, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return out.String() + "\n" + stderr.String(), fmt.Errorf("command failed: %w", err)
	}
	return out.String(), nil
}

func createServerHandler(c *fiber.Ctx) error {
	var req ServerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	terraformDir := "./internal/terraform"
	ansibleDir, _ := filepath.Abs("./internal/ansible")


	if _, err := os.Stat(terraformDir); os.IsNotExist(err) {
		return c.Status(500).JSON(fiber.Map{"error": "terraform directory missing"})
	}

	tfvars := map[string]interface{}{
		"instance_type": ramToInstance(req.RAM),
		"minecraft_tag": req.ServerName,
	}

	tfvarsJSON, _ := json.MarshalIndent(tfvars, "", "  ")
	if err := os.WriteFile(filepath.Join(terraformDir, "terraform.tfvars.json"), tfvarsJSON, 0644); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "cannot write tfvars", "details": err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	// Terraform init
	out, err := runCmd(ctx, terraformDir, "terraform", "init", "-input=false")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "terraform init failed", "details": out})
	}
	log.Println("Terraform initialized.")


	// Terraform apply
	out, err = runCmd(ctx, terraformDir, "terraform", "apply", "-auto-approve", "-input=false")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "terraform apply failed", "details": out})
	}
	log.Println("Terraform applied.")

	// Terraform output
	out, err = runCmd(ctx, terraformDir, "terraform", "output", "-json")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "terraform output failed", "details": out})
	}

	var tfOut TerraformOutputs
	if err := json.Unmarshal([]byte(out), &tfOut); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "parse terraform output failed", "details": err.Error()})
	}

	publicIP := tfOut.PublicIp.Value
	privateKeyPath := tfOut.PrivateKeyPath.Value

	

	inventoryContent := fmt.Sprintf("[mc]\n%s ansible_user=ubuntu ansible_private_key_file=%s\n", publicIP, privateKeyPath)
	inventoryPath := filepath.Join(ansibleDir, "inventory.ini")

	log.Printf("📁 Using inventory at: %s", inventoryPath)

	if err := os.WriteFile(inventoryPath, []byte(inventoryContent), 0600); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "cannot write inventory", "details": err.Error()})
	}

	jarURL := "https://piston-data.mojang.com/v1/objects/95495a7f485eedd84ce928cef5e223b757d2f764/server.jar"
	extraVars := map[string]string{
		"mc_server_jar_url": jarURL,
		"world_name":        req.ServerName,
		"seed":              req.Seed,
		"mc_ram":            req.RAM,
	}
	extraVarsBytes, _ := json.Marshal(extraVars)
	extraVarsArg := fmt.Sprintf("--extra-vars=%s", string(extraVarsBytes))

	ansibleCtx, cancelAnsible := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancelAnsible()

	ansibleCmd := exec.CommandContext(ansibleCtx,
    "ansible-playbook",
    "-i", inventoryPath,
    "--ssh-extra-args", "-o StrictHostKeyChecking=no",
    "playbook.yml",
    extraVarsArg,
)

	ansibleCmd.Dir = ansibleDir
	ansibleCmd.Stdout = os.Stdout
	ansibleCmd.Stderr = os.Stderr

	log.Println("⏳ Waiting 60 seconds for instance to initialize SSH...")
time.Sleep(60 * time.Second)

	log.Printf("▶ Running Ansible playbook for %s...", req.ServerName)
	if err := ansibleCmd.Run(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "ansible failed", "details": err.Error()})
	}

	log.Printf("✅ Minecraft server deployed: %s", publicIP)
	return c.JSON(fiber.Map{
		"message":   "Server created and configured successfully",
		"public_ip": publicIP,
	})
}

func destroyServerHandler(c *fiber.Ctx) error {
	terraformDir := "./internal/terraform"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	out, err := runCmd(ctx, terraformDir, "terraform", "destroy", "-auto-approve", "-input=false")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "terraform destroy failed", "details": out})
	}

	log.Println("Terraform destroyed.")

	return c.JSON(fiber.Map{
		"message": "Server destroyed successfully",
		"details": out,
	})
	
}

func main() {
    app := fiber.New()

    // ✅ Enable CORS for frontend requests
    app.Use(cors.New(cors.Config{
        AllowOrigins: "*", // or "http://localhost:5173" for security
        AllowMethods: "GET,POST,DELETE,OPTIONS",
        AllowHeaders: "Origin, Content-Type, Accept",
    }))

    app.Post("/api/create-server", createServerHandler)
    app.Delete("/api/delete-server", destroyServerHandler)

    log.Println("🚀 EnderCloud backend listening on :8080")
    log.Fatal(app.Listen(":8080"))
}

