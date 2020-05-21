// starting from https://github.com/yuin/gopher-lua/blob/master/cmd/glua/glua.go
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chzyer/readline"
	nlua "github.com/midnightfreddie/nbt-go-lua"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

func main() {
	os.Exit(mainAux())
}

func mainAux() int {
	var opt_e string
	var opt_i, opt_v bool
	flag.StringVar(&opt_e, "e", "", "")
	// flag.StringVar(&opt_l, "l", "", "")
	// flag.StringVar(&opt_p, "p", "", "")
	flag.BoolVar(&opt_i, "i", false, "")
	flag.BoolVar(&opt_v, "v", false, "")
	// flag.BoolVar(&opt_dt, "dt", false, "")
	// flag.BoolVar(&opt_dc, "dc", false, "")
	flag.Usage = func() {
		fmt.Println(`Usage: luanbt [options] [script [args]].
Available options are:
  -e stat  execute string 'stat'
  -i       enter interactive mode after executing 'script'
  -v       show version information`)
	}
	flag.Parse()
	if len(opt_e) == 0 && !opt_i && !opt_v && flag.NArg() == 0 {
		opt_i = true
	}

	status := 0

	// We'll default to Java encoding for this executable
	nlua.UseJavaEncoding()

	// Create gopher-lua environment
	L := nlua.NewState()
	defer L.Close()

	if opt_v || opt_i {
		fmt.Println("nbtlua early release Copyright (C) 2020 Jim Nelson")
		fmt.Println("  based on")
		fmt.Println(lua.PackageCopyRight)
	}

	// if len(opt_l) > 0 {
	// 	if err := L.DoFile(opt_l); err != nil {
	// 		fmt.Println(err.Error())
	// 	}
	// }

	if nargs := flag.NArg(); nargs > 0 {
		script := flag.Arg(0)
		argtb := L.NewTable()
		for i := 1; i < nargs; i++ {
			L.RawSet(argtb, lua.LNumber(i), lua.LString(flag.Arg(i)))
		}
		L.SetGlobal("arg", argtb)
		if err := L.DoFile(script); err != nil {
			fmt.Println(err.Error())
			status = 1
		}
	}

	if len(opt_e) > 0 {
		if err := L.DoString(opt_e); err != nil {
			fmt.Println(err.Error())
			status = 1
		}
	}

	if opt_i {
		fmt.Println("\nWARNING! Early release! Back up all files before modifying!")
		fmt.Print("Load an NBT file with loadnbt(path-to-nbt). ")
		fmt.Print(`Try print(nbt[1].name) or tagType or value. Try changing the name or value. `)
		fmt.Print("Save an NBT file with savenbt(path-to-modified-nbt, true), where the second parameter is whether to compress the output or not. ")
		fmt.Println("Press control-D to exit. ")
		doREPL(L)
	}

	return status
}

// do read/eval/print/loop
func doREPL(L *lua.LState) {
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()
	for {
		if str, err := loadline(rl, L); err == nil {
			if err := L.DoString(str); err != nil {
				fmt.Println(err)
			}
		} else { // error on loadline
			fmt.Println(err)
			return
		}
	}
}

func incomplete(err error) bool {
	if lerr, ok := err.(*lua.ApiError); ok {
		if perr, ok := lerr.Cause.(*parse.Error); ok {
			return perr.Pos.Line == parse.EOF
		}
	}
	return false
}

func loadline(rl *readline.Instance, L *lua.LState) (string, error) {
	rl.SetPrompt("> ")
	if line, err := rl.Readline(); err == nil {
		if _, err := L.LoadString("return " + line); err == nil { // try add return <...> then compile
			return line, nil
		} else {
			return multiline(line, rl, L)
		}
	} else {
		return "", err
	}
}

func multiline(ml string, rl *readline.Instance, L *lua.LState) (string, error) {
	for {
		if _, err := L.LoadString(ml); err == nil { // try compile
			return ml, nil
		} else if !incomplete(err) { // syntax error , but not EOF
			return ml, nil
		} else {
			rl.SetPrompt(">> ")
			if line, err := rl.Readline(); err == nil {
				ml = ml + "\n" + line
			} else {
				return "", err
			}
		}
	}
}
