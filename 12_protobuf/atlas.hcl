// Atlas configuration for 12_protobuf.
//
// Run atlas commands from the 12_protobuf directory, e.g.:
//
//   atlas migrate status --env local

env "local" {
  // Local development database started by compose.yaml.
  // These credentials are for local development only and are not secret.
  url = "postgres://goworld:goworld@localhost:5432/goworld?sslmode=disable"

  // Versioned migration files.
  migration {
    dir = "file://db/migrations"
  }

  // Temporary database Atlas uses to compute diffs and lint migrations
  // (atlas migrate diff / lint). Atlas starts and removes it in Docker.
  dev = "docker://postgres/17/dev?search_path=public"
}
