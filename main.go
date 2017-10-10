package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/davrodpin/ovsdbviz/graphviz"
	"github.com/davrodpin/ovsdbviz/ovsdb"
)

const (
	tableAttrRow         = `<tr><td port="f%d" border="1" bgcolor="%s">%s</td></tr>`
	tableBackgroundColor = "turquoise1"
)

func CreateLabel(columns []string) string {
	var label []string
	for index, columnName := range columns {
		tableBgColor := "transparent"
		if index == 0 {
			tableBgColor = tableBackgroundColor
		}

		label = append(label, fmt.Sprintf(tableAttrRow, index, tableBgColor, columnName))
	}

	return strings.Join(label, "")
}

func GetPortIndex(columns []string, column string) int {
	portIndex := 0 // pointing to the table name by default
	for i, columnName := range columns {
		if columnName == column {
			portIndex = i
			break
		}
	}

	return portIndex
}

var schemaPath schemaValue
var outputPath outputValue

func init() {
	flag.Var(&schemaPath, "schema", "ovsdb schema file path")
	flag.Var(&outputPath, "out", "graphviz dot output file path")
	flag.Parse()

	if schemaPath == "" || outputPath == "" {
		fmt.Printf("usage:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

// TODO: Filter graph by source table (show part of the diagram)
// TODO: Show only attributes that have references to another table
// TODO: Add support to retrieve schema through TCP connection (tcp:127.0.0.1:6644)
// TODO: Show only tables with no column names
func main() {

	schema, err := ovsdb.NewDatabaseSchema(schemaPath.String())
	if err != nil {
		panic(err)
	}

	// Need to always iterate all column for a given table following the same order
	// in order to build and reference graphviz node ports
	tableColumnOrder := schema.OrderedColumns()

	graph := graphviz.NewGraph()

	// NODES
	for tableName, columnOrder := range tableColumnOrder {
		label := CreateLabel(columnOrder)
		nodeAttrs := make(map[string]string)
		nodeAttrs["shape"] = "none"
		nodeAttrs["label"] = fmt.Sprintf(`<<table border="0" cellspacing="0">%s</table>>`, label)
		graph.AddNode(tableName, nodeAttrs)

	}

	// EDGES
	for tableName, table := range schema.Tables {
		for cn, column := range table.Columns {
			references := column.RefersTo()
			if len(references) > 0 {

				portIndex := GetPortIndex(tableColumnOrder[tableName], cn)

				for refAttribute, reference := range references {
					src := tableName
					srcPort := fmt.Sprintf(":f%d", portIndex)
					dst := reference
					dstPort := ":f0"

					edgeAttrs := make(map[string]string)
					edgeAttrs["label"] = refAttribute
					edgeAttrs["splines"] = "polyline"
					switch refAttribute {
					case "key":
						edgeAttrs["color"] = "red"
					case "value":
						edgeAttrs["color"] = "blue"
					}

					graph.AddEdge(src, srcPort, dst, dstPort, edgeAttrs)
				}
			}
		}
	}

	output, err := os.Create(outputPath.String())
	if err != nil {
		panic(err)
	}
	defer output.Close()

	_, err = output.WriteString(graph.String())
	if err != nil {
		panic(fmt.Sprintf("Error while writing output to %s: %v", outputPath.String(), err))
	}

}
