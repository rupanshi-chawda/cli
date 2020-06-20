package cgapp

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/create-go-app/cli/internal/embed"
	"github.com/urfave/cli/v2"
)

// CreateCLIAction actions for `create` CLI command
func CreateCLIAction(c *cli.Context) error {
	// START message
	SendMessage("[*] Create Go App v"+version, "yellow")
	SendMessage("\n[START] Creating a new app...", "green")

	// Create main folder for app
	SendMessage("\n[PROCESS] App folder and config files", "cyan")
	ErrChecker(os.Mkdir(appPath, 0750))
	SendMessage("[OK] App folder was created!", "")

	// Create config files for app
	ErrChecker(File(".editorconfig", embed.Get("/dotfiles/.editorconfig")))
	ErrChecker(File(".gitignore", embed.Get("/dotfiles/.gitignore")))
	ErrChecker(File("Makefile", embed.Get("/dotfiles/Makefile")))

	// Create backend files
	SendMessage("\n[PROCESS] App backend", "cyan")
	ErrChecker(
		Create(&Config{
			name:   strings.ToLower(appBackend),
			match:  "^(net/http|fiber|echo)$",
			view:   "backend",
			folder: appPath,
		},
			registry,
		),
	)

	// Create frontend files
	if appFrontend != "none" {
		SendMessage("\n[PROCESS] App frontend", "cyan")
		ErrChecker(
			Create(&Config{
				name:   strings.ToLower(appFrontend),
				match:  "^(preact|react-js|react-ts)$",
				view:   "frontend",
				folder: appPath,
			},
				registry,
			),
		)

		// Install dependencies
		SendMessage("\n[PROCESS] Frontend dependencies", "cyan")
		SendMessage("[WAIT] Installing frontend dependencies (may take some time)!", "yellow")

		// Go to ./frontend folder and run npm install
		cmd := exec.Command("npm", "install")
		cmd.Dir = filepath.Join(appPath, "frontend")
		ErrChecker(cmd.Run())

		SendMessage("[OK] Frontend dependencies was installed!", "green")
	}

	// Docker containers
	SendMessage("\n[START] Configuring Docker containers...", "green")
	SendMessage("\n[PROCESS] File docker-compose.yml", "cyan")

	// Check frontend
	if appFrontend != "none" {
		// If `-f` argument exists, create fullstack app docker-compose override file
		ErrChecker(File("docker-compose.yml", embed.Get("/docker/docker-compose.fullstack.yml")))
	} else {
		// Default docker-compose.yml
		ErrChecker(File("docker-compose.yml", embed.Get("/docker/docker-compose.backend.yml")))
	}

	// Production settings docker-compose.prod.yml
	ErrChecker(File("docker-compose.prod.yml", embed.Get("/docker/docker-compose.prod.yml")))

	// Create container files
	if appWebServer != "none" {
		SendMessage("\n[PROCESS] Web/proxy server", "cyan")
		ErrChecker(
			Create(&Config{
				name:   strings.ToLower(appWebServer),
				match:  "^(nginx)$",
				view:   "webserver",
				folder: appPath,
			},
				registry,
			),
		)
	}

	// END message
	SendMessage("\n[DONE] Run `make` from `"+appPath+"` folder!", "yellow")

	return nil
}
