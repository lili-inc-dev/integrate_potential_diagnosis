package generator

import (
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"

	"github.com/iancoleman/strcase"
)

func GenerateProtoFile(inputFilePath string, model string) error {
	structFields := paramCreatorForProto(inputFilePath)
	fmt.Printf("message %sParams {\n%s}\n", strcase.ToCamel(model), structFields)
	return nil
}

func paramCreatorForProto(filePath string) string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		log.Fatal(err)
	}
	typeInfo := make([]string, 0)
	dbTag := make([]string, 0)
	// Make sure that the db tag is attached.
	ast.Inspect(f, func(n ast.Node) bool {
		ident, ok := n.(*ast.Field)
		if !ok {
			return true
		}
		if ident.Tag != nil {
			t := fmt.Sprintf("%s", ident.Type)
			tag := ident.Tag.Value
			_, right, ok := strings.Cut(tag, "db:")
			if ok {
				typeInfo = append(typeInfo, t)
				dbTag = append(dbTag, right[1:len(right)-2])
			}
		}
		return true
	})
	var ans string
	cnt := 1
	for i := range typeInfo {
		typ := typeInfo[i]
		tag := dbTag[i]
		if strings.Contains(typ, "time") {
			typ = "google.protobuf.Timestamp"
		}
		if strings.Contains(typ, "NullString") {
			typ = "string"
		}
		ans += fmt.Sprintf("	%s %s = %d;\n", typ, tag, cnt)
		cnt++
	}
	return ans
}
