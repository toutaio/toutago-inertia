package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	output := flag.String("o", "views/types.ts", "Output file path")
	flag.Parse()

	types := make(map[string]*TypeDef)

	// Parse Go files in handlers and models
	dirs := []string{"handlers", "models"}
	for _, dir := range dirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
				return err
			}

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return err
			}

			ast.Inspect(file, func(n ast.Node) bool {
				typeSpec, ok := n.(*ast.TypeSpec)
				if !ok {
					return true
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					return true
				}

				types[typeSpec.Name.Name] = parseStruct(structType)
				return true
			})

			return nil
		})
	}

	// Generate TypeScript
	content := "// Auto-generated TypeScript types from Go structs\n"
	content += "// Do not edit manually\n\n"

	for name, typeDef := range types {
		content += fmt.Sprintf("export interface %s {\n", name)
		for _, field := range typeDef.Fields {
			optional := ""
			if field.Optional {
				optional = "?"
			}
			content += fmt.Sprintf("  %s%s: %s\n", field.Name, optional, field.Type)
		}
		content += "}\n\n"
	}

	os.WriteFile(*output, []byte(content), 0644)
	fmt.Printf("Generated TypeScript types: %s\n", *output)
}

type TypeDef struct {
	Fields []Field
}

type Field struct {
	Name     string
	Type     string
	Optional bool
}

func parseStruct(s *ast.StructType) *TypeDef {
	typeDef := &TypeDef{}

	for _, field := range s.Fields.List {
		if len(field.Names) == 0 {
			continue
		}

		jsonTag := ""
		if field.Tag != nil {
			tag := field.Tag.Value
			if strings.Contains(tag, "json:") {
				parts := strings.Split(tag, "json:")
				if len(parts) > 1 {
					jsonTag = strings.Trim(strings.Split(parts[1], "\"")[1], ",")
				}
			}
		}

		if jsonTag == "-" {
			continue
		}

		name := field.Names[0].Name
		if jsonTag != "" {
			name = strings.Split(jsonTag, ",")[0]
		}

		optional := strings.Contains(field.Tag.Value, ",omitempty")

		tsType := goTypeToTS(field.Type)

		typeDef.Fields = append(typeDef.Fields, Field{
			Name:     name,
			Type:     tsType,
			Optional: optional,
		})
	}

	return typeDef
}

func goTypeToTS(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		switch t.Name {
		case "string":
			return "string"
		case "int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"float32", "float64":
			return "number"
		case "bool":
			return "boolean"
		default:
			return t.Name
		}
	case *ast.ArrayType:
		return goTypeToTS(t.Elt) + "[]"
	case *ast.MapType:
		return fmt.Sprintf("Record<%s, %s>", goTypeToTS(t.Key), goTypeToTS(t.Value))
	case *ast.StarExpr:
		return goTypeToTS(t.X)
	}
	return "any"
}
