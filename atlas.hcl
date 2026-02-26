data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "./cmd/atlasloader",
  ]
}

env "local" {
  src = data.external_schema.gorm.url
  # A throw-away dev DB Atlas uses to validate SQL before writing migrations.
  # Must exist; Atlas will NOT apply migrations to it automatically.
  dev = "postgres://postgres:postgres@localhost:5432/workout_tracker_dev?sslmode=disable"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
