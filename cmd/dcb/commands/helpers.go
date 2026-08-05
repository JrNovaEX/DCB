package commands

import (
	"fmt"
	"os"

	"github.com/JrNovaEX/DCB/internal/adapters/generator"
	"github.com/JrNovaEX/DCB/internal/adapters/parser"
	"github.com/JrNovaEX/DCB/internal/adapters/runner"
	"github.com/JrNovaEX/DCB/internal/usecases"
)

// green, red, yellow ANSI helpers for terminal output.
const (
	colorGreen = "\033[32m"
	colorRed   = "\033[31m"
	colorReset = "\033[0m"
)

func printSuccess(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, colorGreen+"✓ "+colorReset+format+"\n", a...)
}

func printError(err error) {
	fmt.Fprintf(os.Stderr, colorRed+"✗ "+colorReset+"%v\n", err)
}

// newBuildUseCase wires up the BuildUseCase with concrete adapters.
func newBuildUseCase() *usecases.BuildUseCase {
	p := parser.New()
	g := generator.New()
	r := runner.New()
	return usecases.NewBuildUseCase(p, g, r)
}

// newCheckUseCase wires up the CheckUseCase with concrete adapters.
func newCheckUseCase() *usecases.CheckUseCase {
	p := parser.New()
	pc := runner.NewPortChecker()
	return usecases.NewCheckUseCase(p, pc)
}

// newUpUseCase wires up the UpUseCase with concrete adapters.
func newUpUseCase() *usecases.UpUseCase {
	r := runner.New()
	build := newBuildUseCase()
	return usecases.NewUpUseCase(build, r)
}

// newDownUseCase wires up the DownUseCase.
func newDownUseCase() *usecases.DownUseCase {
	return usecases.NewDownUseCase(runner.New())
}

// newLogsUseCase wires up the LogsUseCase.
func newLogsUseCase() *usecases.LogsUseCase {
	return usecases.NewLogsUseCase(runner.New())
}
