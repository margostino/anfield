package context

var configuration = GetConfig("./configuration/configuration.yml")

func Config() *Configuration {
	return configuration
}

func ShouldUpdateData() bool {
	return configuration.Data.Update
}
