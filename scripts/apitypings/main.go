package main

import (
	"fmt"
	"log"

	"github.com/coder/guts"
	"github.com/coder/guts/config"

	// Must import the packages we are trying to convert
	// And include the ones we are referencing
	_ "github.com/Emyrk/strava/api/modelsdk/sdktype"
	_ "github.com/Emyrk/strava/api/superlative"
)

func main() {
	gen, err := guts.NewGolangParser()
	if err != nil {
		log.Fatalf("new convert: %v", err)
	}

	generateDirectories := map[string]string{
		"github.com/Emyrk/strava/api/modelsdk":         "",
		"github.com/Emyrk/strava/api/modelsdk/sdktype": "",
		"github.com/Emyrk/strava/api/superlative":      "Superlative",
	}
	for dir, prefix := range generateDirectories {
		err = gen.IncludeGenerateWithPrefix(dir, prefix)
		if err != nil {
			log.Fatalf("include generate package %q: %v", dir, err)
		}
	}

	referencePackages := map[string]string{}
	for pkg, prefix := range referencePackages {
		err = gen.IncludeReference(pkg, prefix)
		if err != nil {
			log.Fatalf("include reference package %q: %v", pkg, err)
		}
	}

	err = TypeMappings(gen)
	if err != nil {
		log.Fatalf("type mappings: %v", err)
	}

	// We are configured
	ts, err := gen.ToTypescript()
	if err != nil {
		log.Fatalf("to typescript: %v", err)
	}

	// We have a TS AST
	TsMutations(ts)

	output, err := ts.Serialize()
	if err != nil {
		log.Fatalf("serialize: %v", err)
	}
	fmt.Println(output)
}

func TsMutations(ts *guts.Typescript) {
	ts.ApplyMutations(
		// Enum list generator
		config.EnumLists,
		// Export all top level types
		config.ExportTypes,
		// Readonly interface fields
		//config.ReadOnly,
		// Add ignore linter comments
		config.BiomeLintIgnoreAnyTypeParameters,
		// Omitempty + null is just '?' in golang json marshal
		// number?: number | null --> number?: number
		config.SimplifyOmitEmpty,
	)
}

func TypeMappings(gen *guts.GoParser) error {
	gen.IncludeCustomDeclaration(config.StandardMappings())

	err := gen.IncludeCustom(map[string]string{
		// Serpent fields should be converted to their primitive types
		"github.com/coder/serpent.Regexp":                        "string",
		"github.com/coder/serpent.StringArray":                   "string",
		"github.com/coder/serpent.String":                        "string",
		"github.com/coder/serpent.YAMLConfigPath":                "string",
		"github.com/coder/serpent.Strings":                       "[]string",
		"github.com/coder/serpent.Int64":                         "int64",
		"github.com/coder/serpent.Bool":                          "bool",
		"github.com/coder/serpent.Duration":                      "int64",
		"github.com/coder/serpent.URL":                           "string",
		"github.com/coder/serpent.HostPort":                      "string",
		"github.com/Emyrk/strava/api/modelsdk/sdktype.StringInt": "string",
		"github.com/Emyrk/strava/api/modelsdk.StringInt":         "string",
		"encoding/json.RawMessage":                               "map[string]string",
	})
	if err != nil {
		return fmt.Errorf("include custom: %w", err)
	}

	return nil
}
