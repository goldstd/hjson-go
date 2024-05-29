package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/goldstd/hjson-go"
)

// Version Can be set when building for example like this:
// go build -ldflags "-X main.Version=v3.0"
var Version string

func fixJSON(data []byte) []byte {
	data = bytes.Replace(data, []byte("\\u003c"), []byte("<"), -1)
	data = bytes.Replace(data, []byte("\\u003e"), []byte(">"), -1)
	data = bytes.Replace(data, []byte("\\u0026"), []byte("&"), -1)
	data = bytes.Replace(data, []byte("\\u0008"), []byte("\\b"), -1)
	data = bytes.Replace(data, []byte("\\u000c"), []byte("\\f"), -1)
	return data
}

func main() {
	flag.Usage = func() {
		fmt.Println("usage: hjson [OPTIONS] [INPUT]")
		fmt.Println("hjson can be used to convert JSON from/to Hjson.")
		fmt.Println("")
		fmt.Println("hjson will read the given JSON/Hjson input file or read from stdin.")
		fmt.Println("")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	var (
		help        = flag.Bool("h", false, "Show this screen.")
		showJSON    = flag.Bool("j", false, "Output as formatted JSON.")
		showCompact = flag.Bool("c", false, "Output as JSON.")

		indentBy = flag.String("indent", "  ", "The indent string.")
		eol      = flag.String("eol", "\n", `End of line, should be either \n or \r\n`)

		bracesSameLine   = flag.Bool("bracesSameLine", false, "Print braces on the same line.")
		omitRootBraces   = flag.Bool("omitRootBraces", false, "Omit braces at the root.")
		quoteAlways      = flag.Bool("quoteAlways", false, "Always quote string values.")
		showVersion      = flag.Bool("v", false, "Show version.")
		preserveKeyOrder = flag.Bool("preserveKeyOrder", false, "Preserve key order in objects/maps.")
	)

	flag.Parse()
	if *help || flag.NArg() > 1 {
		flag.Usage()
		os.Exit(1)
	}

	if *showVersion {
		if Version != "" {
			fmt.Println(Version)
		} else if bi, ok := debug.ReadBuildInfo(); ok {
			fmt.Println(bi.Main.Version)
		} else {
			fmt.Println("Unknown version")
		}
		os.Exit(0)
	}

	var err error
	var data []byte
	if flag.NArg() == 1 {
		data, err = os.ReadFile(flag.Arg(0))
	} else {
		data, err = io.ReadAll(os.Stdin)
	}
	if err != nil {
		panic(err)
	}

	var value interface{}

	if *preserveKeyOrder {
		var node *hjson.Node
		err = hjson.Unmarshal(data, &node)
		value = node
	} else {
		err = hjson.Unmarshal(data, &value)
	}
	if err != nil {
		panic(err)
	}

	var out []byte
	if *showCompact {
		out, err = json.Marshal(value)
		if err != nil {
			panic(err)
		}
		out = fixJSON(out)
	} else if *showJSON {
		out, err = json.MarshalIndent(value, "", *indentBy)
		if err != nil {
			panic(err)
		}
		out = fixJSON(out)
	} else {
		opt := hjson.DefaultOptions()
		opt.IndentBy = *indentBy
		opt.BracesSameLine = *bracesSameLine
		opt.EmitRootBraces = !*omitRootBraces
		opt.QuoteAlways = *quoteAlways
		opt.Comments = false
		opt.Eol = *eol
		out, err = hjson.MarshalWithOptions(value, opt)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(string(out))
}