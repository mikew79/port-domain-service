package entities

// Configuration - The global entity representing the necessary configuration options for the application
type Configuration struct {
	MongoURI   string
	DbName     string
	PortNumber int
}
