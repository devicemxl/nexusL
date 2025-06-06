// cdm/nexusl/main.go

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	//"github.com/devicemxl/nexusl/pkg/environment"
	"github.com/devicemxl/nexusl/pkg/evaluator"
	"github.com/devicemxl/nexusl/pkg/evaluator/environment"
	"github.com/devicemxl/nexusl/pkg/evaluator/store"
	"github.com/devicemxl/nexusl/pkg/lexer"
	"github.com/devicemxl/nexusl/pkg/object" // Importa el paquete object
	"github.com/devicemxl/nexusl/pkg/parser"
	//"github.com/devicemxl/nexusl/pkg/store"
)

const PROMPT = "nexusl > "

func main() {
	// 1. Inicializar el entorno global y el knowledge store
	knowledgeStore := store.NewKnowledgeStore()
	env := environment.NewEnvironment()
	// X Shitposting
	multiLineString := `
                                                                                    

███╗   ██╗███████╗██╗  ██╗██╗   ██╗███████╗██╗     
████╗  ██║██╔════╝╚██╗██╔╝██║   ██║██╔════╝██║     
██╔██╗ ██║█████╗   ╚███╔╝ ██║   ██║███████╗██║     
██║╚██╗██║██╔══╝   ██╔██╗ ██║   ██║╚════██║██║     
██║ ╚████║███████╗██╔╝ ██╗╚██████╔╝███████║███████╗
╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚══════╝


`

	fmt.Printf(multiLineString)
	// 2. Inicializar símbolos built-in (como "cli")
	// Crear el SymbolObject para "cli"
	cliSymbol := object.NewSymbolObject("cli")

	// Asignar la función BuiltinPrint al miembro "print" del objeto "cli"
	// object.Builtins["print"] proviene de pkg/object/builtin.go
	cliSymbol.SetMember("print", object.Builtins["print"])

	// Añadir el símbolo "cli" pre-configurado al entorno global
	env.Set("cli", cliSymbol)

	// ... puedes inicializar otros built-ins aquí si los tienes ...
	// Ej: mathSymbol := object.NewSymbolObject("math")
	// mathSymbol.SetMember("add", object.Builtins["add"])
	// env.Set("math", mathSymbol)

	// Iniciar el REPL
	startREPL(os.Stdin, os.Stdout, env, knowledgeStore)
}

func startREPL(in io.Reader, out io.Writer, env *environment.Environment, knowledgeStore *store.KnowledgeStore) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.NewLexer(line)
		p := parser.NewParser(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env, knowledgeStore)
		if evaluated != nil {
			io.WriteString(out, evaluated.String())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
