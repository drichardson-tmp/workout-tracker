// atlasloader is invoked by Atlas via `atlas.hcl` to produce the desired
// database schema from GORM models. Run via:
//
//	atlas migrate diff --env local
package main

import (
	"fmt"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"workout-tracker/backend/models"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&models.User{},
		&models.Workout{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(stmts)
}
