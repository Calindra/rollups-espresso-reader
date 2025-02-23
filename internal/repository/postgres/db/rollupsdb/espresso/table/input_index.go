//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var InputIndex = newInputIndexTable("espresso", "input_index", "")

type inputIndexTable struct {
	postgres.Table

	// Columns
	ApplicationAddress postgres.ColumnString
	Index              postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type InputIndexTable struct {
	inputIndexTable

	EXCLUDED inputIndexTable
}

// AS creates new InputIndexTable with assigned alias
func (a InputIndexTable) AS(alias string) *InputIndexTable {
	return newInputIndexTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new InputIndexTable with assigned schema name
func (a InputIndexTable) FromSchema(schemaName string) *InputIndexTable {
	return newInputIndexTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new InputIndexTable with assigned table prefix
func (a InputIndexTable) WithPrefix(prefix string) *InputIndexTable {
	return newInputIndexTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new InputIndexTable with assigned table suffix
func (a InputIndexTable) WithSuffix(suffix string) *InputIndexTable {
	return newInputIndexTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newInputIndexTable(schemaName, tableName, alias string) *InputIndexTable {
	return &InputIndexTable{
		inputIndexTable: newInputIndexTableImpl(schemaName, tableName, alias),
		EXCLUDED:        newInputIndexTableImpl("", "excluded", ""),
	}
}

func newInputIndexTableImpl(schemaName, tableName, alias string) inputIndexTable {
	var (
		ApplicationAddressColumn = postgres.StringColumn("application_address")
		IndexColumn              = postgres.IntegerColumn("index")
		allColumns               = postgres.ColumnList{ApplicationAddressColumn, IndexColumn}
		mutableColumns           = postgres.ColumnList{IndexColumn}
	)

	return inputIndexTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ApplicationAddress: ApplicationAddressColumn,
		Index:              IndexColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
