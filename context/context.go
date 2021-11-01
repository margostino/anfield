package context

var configuration = GetConfig("./configuration/configuration.yml")

const (
	REALTIME = "realtime"
	BATCH    = "batch"
)

func Config() *Configuration {
	return configuration
}

func ShouldUpdateData() bool {
	return configuration.Data.Update
}
