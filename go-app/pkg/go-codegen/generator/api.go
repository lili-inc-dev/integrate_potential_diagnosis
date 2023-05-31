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

func GenerateApiStruct(inputFilePath string, model string) error {
	fields := CreateStructForApi(inputFilePath)
	upper_model := strcase.ToCamel(model)

	fmt.Printf("Create%sReq{\n%s}\nUpdate%sReq{\n%s}\nDelete%sReq{\n\tId uint64 `json:\"id\"`\n}\nFindById%sReq{\n\tId uint64 `json:\"id\"`\n}\nFindById%sRes{\n%s}\n", upper_model, fields, upper_model, fields, upper_model, upper_model, upper_model, fields)

	return nil
}

func CreateStructForApi(filePath string) string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	typeInfo := make([]string, 0)
	fieldInfo := make([]string, 0)
	dbTag := make([]string, 0)

	// Make sure that the db tag is attached.
	ast.Inspect(f, func(n ast.Node) bool {
		ident, ok := n.(*ast.Field)
		if !ok {
			return true
		}
		if ident.Tag != nil {
			t := fmt.Sprintf("%s", ident.Type)
			f := ident.Names[0].String()
			tag := ident.Tag.Value
			_, right, ok := strings.Cut(tag, "db:")
			if ok {
				typeInfo = append(typeInfo, t)
				fieldInfo = append(fieldInfo, f)
				dbTag = append(dbTag, right[1:len(right)-2])
			}
		}
		return true
	})
	var ans string
	for i := range typeInfo {
		typ := typeInfo[i]
		field := fieldInfo[i]
		tag := dbTag[i]
		if strings.Contains(typ, "time") {
			typ = "string"
		}
		if strings.Contains(typ, "NullString") {
			typ = "string"
			tag = fmt.Sprintf("%s,%s,%s", tag, "optional", "omitempty")
		}
		ans += fmt.Sprintf("	%s %s `json:\"%s\"`\n", field, typ, tag)
	}
	return ans
}
