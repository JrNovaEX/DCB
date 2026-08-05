// DCB — Docker Compose Builder
// Build-time variables are injected by goreleaser via -ldflags.
package main

import "github.com/JrNovaEX/DCB/cmd/dcb/commands"

func main() {
	commands.Execute()
}
