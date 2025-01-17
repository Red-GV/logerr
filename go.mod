module github.com/ViaQ/logerr

go 1.17

require (
	github.com/go-logr/logr v0.4.0
	github.com/stretchr/testify v1.4.0
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
)

retract [v1.1.0, v1.1.1] // Improper versioning of breaking changes
