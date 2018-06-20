package common

// DbConfigResult represents a row from DB query result
type DbConfigResult struct {
	// Property indicates DB column named "Property"
	Property string `json:"property"`
	// Value
	Value string `json:"value"`
}

const (
	// DbConfigTblName is a config table name
	DbConfigTblName = "Config"
	// DbConfigPropName is a property colum name
	DbConfigPropName = "Property"
	// DbConfigValName is a value colum name
	DbConfigValName = "Value"
)
