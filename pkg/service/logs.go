package service

// Logging messages.
const (
	Init         = "Initializing"
	Running      = "Running"
	ShuttingDown = "Shutting down"

	ReadingConfig          = "Reading configuration"
	ReadingConfigFailed    = "Reading configuration failed"
	ValidatingConfigFailed = "Validating configuration failed"

	InitDatabase       = "Initializing database"
	InitDatabaseFailed = "Initializing database failed"

	InitRESTServer    = "Initializing REST server"
	RunningRESTServer = "Running REST server"
)
