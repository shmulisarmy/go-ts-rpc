package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"runtime"
	"strings"
)

const typescriptRPCFileName string = "rpc_functions.ts"

// Map of registered functions
var functions = make(map[string]interface{})

// TypeScript type mappings
var typescriptTypes = map[string]string{
	"int":          "number",
	"string":       "string",
	"float":        "number",
	"bool":         "boolean",
	"interface {}": "null",
	"any":          "any",
}

type function_info struct {
	name            string
	args            []string
	args_with_types []string
	go_types        []string
	file_from       string
	returns         string
}

// Helper to get the function name
func getFunctionName(fn interface{}) string {
	ptr := reflect.ValueOf(fn).Pointer()
	funcObj := runtime.FuncForPC(ptr)
	return strings.Split(funcObj.Name(), ".")[1]
}

// Register a function in the map
func addFunction(fn interface{}) {
	functionName := getFunctionName(fn)
	functions[functionName] = fn
}

func get_param_names_from_ast(funcName string, filename string) ([]string, error) {
	// Open the source file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Parse the source file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, file, parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %v", err)
	}

	// Traverse the AST to find the function
	var paramNames []string
	ast.Inspect(node, func(n ast.Node) bool {
		// Check for function declarations
		if funcDecl, ok := n.(*ast.FuncDecl); ok && funcDecl.Name.Name == funcName {
			for _, param := range funcDecl.Type.Params.List {
				for _, name := range param.Names {
					paramNames = append(paramNames, name.Name)
				}
			}
			return false // Stop traversing
		}
		return true
	})

	if len(paramNames) == 0 {
		return nil, fmt.Errorf("function %s not found in file %s", funcName, filename)
	}

	return paramNames, nil
}

func make_typescript_function_call_code(functionInfo *function_info) string {
	functionContent := fmt.Sprintf(`export function %s(%s): %s {
    return rpc_call("%s", %s);
}`, functionInfo.name, strings.Join(functionInfo.args_with_types, ", "), functionInfo.returns, functionInfo.name, strings.Join(functionInfo.args, ", "))
	return functionContent
}

// Generate TypeScript definitions for a function
func setup_rpc(fn interface{}) {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		panic("Not a function")
	}

	// Extract argument types
	filename := function_to_file[getFunctionName(fn)]
	if filename == "" {
		panic("Function's filename is not known wich means not all files were loaded not found")
	}
	argCount := fnType.NumIn()
	param_names, err := get_param_names_from_ast(getFunctionName(fn), filename)
	if err != nil {
		panic(err)
	}
	typescriptArgs := make([]string, argCount)
	go_types := make([]string, argCount)
	for i := 0; i < argCount; i++ {
		goType := fnType.In(i).String()
		go_types[i] = goType
		tsType, ok := typescriptTypes[goType]
		if !ok {
			panic("Unsupported argument type: " + goType)
		}
		// param_names[i] = fmt.Sprintf("arg%d", i)
		typescriptArgs[i] = fmt.Sprintf("%s: %s", param_names[i], tsType)
	}

	// Extract return type
	var returnType string
	if fnType.NumOut() == 0 {
		returnType = "void"
	} else {
		goType := fnType.Out(0).String()
		tsType, ok := typescriptTypes[goType]
		if !ok {
			panic("Unsupported return type: " + goType)
		}
		returnType = tsType
	}

	// Generate TypeScript function content
	functionName := getFunctionName(fn)
	functionInfo := functions_in_file[filename][functionName]
	functionInfo.go_types = go_types
	functionInfo.args = param_names
	functionInfo.args_with_types = typescriptArgs
	functionInfo.returns = returnType

	functionContent := make_typescript_function_call_code(functionInfo)

	// Write to TypeScript file
	file, err2 := os.OpenFile(typescriptRPCFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err2 != nil {
		panic(err2)
	}
	defer file.Close()

	_, err = file.WriteString(functionContent)
	if err != nil {
		panic(err)
	}

	// Register function
	addFunction(fn)
}

// Generic function caller
func callFunction(funcName string, args []interface{}) (interface{}, error) {
	fn := functions[funcName]
	if fn == nil {
		return nil, errors.New("function not found")
	}

	fnValue := reflect.ValueOf(fn)
	if len(args) != fnValue.Type().NumIn() {
		return nil, errors.New("argument count mismatch")
	}

	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	results := fnValue.Call(in)
	if len(results) == 0 {
		return nil, nil
	}
	return results[0].Interface(), nil
}

var functions_in_file map[string]map[string]*function_info = make(map[string]map[string]*function_info)
var function_to_file map[string]string = make(map[string]string)

func load_file(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Parse the source file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, file, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	ast.Inspect(node, func(n ast.Node) bool {

		func_node, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		if func_node.Name.Name == "main" {
			return false
		}

		param_names, err := get_param_names_from_ast(func_node.Name.Name, filename)
		if err != nil {
			panic(err)
		}

		if functions_in_file[filename] == nil {
			functions_in_file[filename] = make(map[string]*function_info)
		}

		functions_in_file[filename][func_node.Name.Name] = &function_info{
			name:            func_node.Name.Name,
			args:            param_names,
			args_with_types: make([]string, len(param_names)),
		}
		function_to_file[func_node.Name.Name] = filename
		return false
	})
}
