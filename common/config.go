package common

type DbConfigResult struct {
	Property string `json: "Property"`
	Value    string `json: "Value"`
}

const (
	DbConfigTblName  = "Config"
	DbConfigPropName = "Property"
	DbConfigValName  = "Value"
)
