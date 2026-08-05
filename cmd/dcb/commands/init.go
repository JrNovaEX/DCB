package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactively scaffold a new dcb.yaml",
	Long: `Init walks you through a quick Q&A and writes a ready-to-use dcb.yaml
in the current directory (or the directory given by --dir).

If dcb.yaml already exists, you will be asked whether to overwrite it.`,
	Example: `  dcb init
  dcb init --dir ./myproject`,
	RunE: runInit,
}

var initFlagDir string

func init() {
	initCmd.Flags().StringVarP(&initFlagDir, "dir", "d", ".", "Directory to write dcb.yaml into")
	rootCmd.AddCommand(initCmd)
}

func prompt(r *bufio.Reader, label, def string) string {
	if def != "" {
		fmt.Printf("  %s [%s]: ", label, def)
	} else {
		fmt.Printf("  %s: ", label)
	}
	line, _ := r.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return def
	}
	return line
}

// initService holds the gathered data for one service during init.
type initService struct {
	name    string
	image   string
	build   string
	portStr string
}

func runInit(_ *cobra.Command, _ []string) error {
	outPath := filepath.Join(initFlagDir, "dcb.yaml")

	// Guard: existing file
	if _, err := os.Stat(outPath); err == nil {
		fmt.Printf("\033[33m! \033[0m%s already exists. Overwrite? [y/N]: ", outPath)
		r := bufio.NewReader(os.Stdin)
		ans, _ := r.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(ans)) != "y" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	r := bufio.NewReader(os.Stdin)

	fmt.Println("\n  🐳  DCB — Project Initialiser")
	fmt.Println("  ─────────────────────────────")

	projectName := prompt(r, "Project name", filepath.Base(mustGetwd()))
	if projectName == "" {
		projectName = "myapp"
	}

	numServicesStr := prompt(r, "Number of services", "1")
	numServices, _ := strconv.Atoi(numServicesStr)
	if numServices < 1 {
		numServices = 1
	}

	services := make([]initService, 0, numServices)
	for i := 0; i < numServices; i++ {
		fmt.Printf("\n  — Service %d —\n", i+1)
		svc := initService{}
		svc.name = prompt(r, "  Name", fmt.Sprintf("service%d", i+1))

		sourceType := prompt(r, "  Source (image/build)", "image")
		if strings.HasPrefix(sourceType, "b") {
			svc.build = prompt(r, "  Build context path", "./")
		} else {
			svc.image = prompt(r, "  Docker image", "nginx:alpine")
		}
		svc.portStr = prompt(r, "  Exposed port (leave empty to skip)", "")
		services = append(services, svc)
	}

	// Build the YAML document in a structured way so it's always valid.
	doc := buildInitDoc(projectName, services)
	data, err := yaml.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshalling dcb.yaml: %w", err)
	}

	if err := os.MkdirAll(initFlagDir, 0750); err != nil {
		return fmt.Errorf("creating directory %q: %w", initFlagDir, err)
	}

	if err := os.WriteFile(outPath, data, 0640); err != nil {
		return fmt.Errorf("writing %q: %w", outPath, err)
	}

	printSuccess("Created %s", outPath)
	fmt.Printf("\n  Next steps:\n")
	fmt.Printf("    dcb check        → validate your config\n")
	fmt.Printf("    dcb build        → generate docker-compose.yml\n")
	fmt.Printf("    dcb up --detach  → start services in the background\n\n")
	return nil
}

// buildInitDoc constructs a yaml-marshallable map from the gathered answers.
func buildInitDoc(project string, services []initService) map[string]interface{} {
	svcMap := make(map[string]interface{}, len(services))
	for _, s := range services {
		entry := map[string]interface{}{}
		if s.image != "" {
			entry["image"] = s.image
		} else {
			entry["build"] = s.build
		}
		if s.portStr != "" {
			if p, err := strconv.Atoi(s.portStr); err == nil && p > 0 {
				entry["port"] = p
			}
		}
		svcMap[s.name] = entry
	}
	return map[string]interface{}{
		"project":  project,
		"version":  "3.8",
		"services": svcMap,
	}
}

func mustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		return "myapp"
	}
	return wd
}
