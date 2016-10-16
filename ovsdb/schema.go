package ovsdb

import (
	"encoding/json"
	"io/ioutil"
)

type DatabaseSchema struct {
	Name     string                 `json:"name"`
	Version  string                 `json:"version"`
	Checksum string                 `json:"cksum,omitempty"`
	Tables   map[string]TableSchema `json:"tables"`
}

// OrderedColumns returns all column names ordered by a numeric index for each
// table in the schema
func (database DatabaseSchema) OrderedColumns() map[string][]string {
	tableColumnOrder := make(map[string][]string)

	for tableName, table := range database.Tables {
		var columnOrder []string
		columnOrder = append(columnOrder, tableName)
		for columnName := range table.Columns {
			columnOrder = append(columnOrder, columnName)
		}
		tableColumnOrder[tableName] = columnOrder
	}

	return tableColumnOrder

}

type TableSchema struct {
	Columns     map[string]ColumnSchema `json:"columns"`
	ColumnOrder []string
	MaxRows     int        `json:"maxrows,omitempty"`
	IsRoot      bool       `json:"isRoot,omitempty"`
	Indexes     [][]string `json:"indexes,omitempty"`
}

type ColumnSchema struct {
	Type interface{} `json:"type"`
}

func getRefTable(key string, typeColumn interface{}) string {
	typeColumnMap := typeColumn.(map[string]interface{})

	if valueInterface, exists := typeColumnMap[key]; exists {
		switch valueInterface.(type) {
		case map[string]interface{}:
			value := valueInterface.(map[string]interface{})
			if refTable, exists := value["refTable"]; exists {
				return refTable.(string)
			}
		}
	}

	return ""
}

func (column ColumnSchema) RefersTo() map[string]string {

	references := make(map[string]string)

	switch column.Type.(type) {
	case map[string]interface{}:
		keyRefTable := getRefTable("key", column.Type)
		if keyRefTable != "" {
			references["key"] = keyRefTable
		}

		valueRefTable := getRefTable("value", column.Type)
		if valueRefTable != "" {
			references["value"] = valueRefTable
		}
	}

	return references
}

func NewDatabaseSchema(schemaPath string) (*DatabaseSchema, error) {
	fp, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}

	database := DatabaseSchema{}

	if err = json.Unmarshal(fp, &database); err != nil {
		return nil, err
	}

	return &database, nil
}
