package main

import (
	"flag"
	"fmt"
	"github.com/davrodpin/ovsdbviz/graphviz"
	"github.com/davrodpin/ovsdbviz/ovsdb"
	"os"
	"strings"
)

const (
	tableAttrRow         = `<tr><td port="f%d" border="1" bgcolor="%s">%s</td></tr>`
	tableBackgroundColor = "turquoise1"
)

func CreateLabel(table ovsdb.TableSchema, columns []string) string {
	var labels []string

	for index, columnName := range columns {
		tableBgColor := "transparent"
		label := columnName
		if index == 0 {
			tableBgColor = tableBackgroundColor
		}

		if table.IsIndex(columnName) {
			label = fmt.Sprintf("%s (index)", columnName)
		}

		labels = append(labels, fmt.Sprintf(tableAttrRow, index, tableBgColor, label))
	}

	return strings.Join(labels, "")
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
		label := CreateLabel(schema.Tables[tableName], columnOrder)
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
