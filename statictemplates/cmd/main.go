// CLI command to execute statictemplates.
package main

import (
	"fmt"
	"os"

	"github.com/erdnaxeli/wishlister/statictemplates"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Printf(
			"Usage: %s TEMPLATES_DIRECTORY GENERATED_PACKAGE_NAME OUT_DIRECTORY\n",
			os.Args[0],
		)
		os.Exit(1)
	}

	templatesDirectory := os.Args[1]
	packageName := os.Args[2]
	outputDirectory := os.Args[3]

	app, err := statictemplates.NewWithDefaults(templatesDirectory, outputDirectory, packageName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = app.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
