package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"os/exec"
	"github.com/spekary/got/got"
)


func processFile(file string) string {
	buf, err := ioutil.ReadFile(file)

	s := fmt.Sprintf("%s", buf)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	/*
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered ", r)
		}
	}()*/

	s = ProcessString(s, file)
	if s != "" {
		s = "//** This file was code generated by got. ***\n\n\n" + s
	}
	return s
}

func writeFile (s string, file string, outDir string, runImports bool) {

	i := strings.LastIndex(file, ".")

	dir := filepath.Dir(file)
	dir,_ = filepath.Abs(dir)
	file = filepath.Base(file)

	if i < 0 {
		file = file + ".go"
	} else {
		file = file[:i] + ".go"
	}

	if outDir != "" {
		dir = outDir
	}

	if dir != "/" {
		dir = dir + "/"
	}
	file = dir + file

	ioutil.WriteFile(file, []byte(s), os.ModePerm)

	if runImports {
		execCommand("goimports -w " + file)
	} else {
		execCommand("go fmt " + file)	// at least format it if we are not going to run imports on it
	}
}

//Process a string that is a got template, and return the go code
func ProcessString(input string, fileName string) string {
	l := got.Lex(input, fileName)

	s := got.Parse(l)

	return s
}


// execCommand wraps exec.Command
func execCommand(command string) {
	parts := strings.Split(command, " ")
	if len(parts) == 0 {
		return
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	var  outDir string
	var typ string
	var runImports bool
	var includes string

	flag.StringVar(&outDir, "o", "", "Output directory")
	flag.StringVar(&typ, "t", "", "Will process all files with this suffix in current directory")
	flag.BoolVar(&runImports, "i", false, "Run goimports on the file to automatically add your imports to the file. You will need to install goimports to do this.")
	flag.StringVar(&includes, "I", "", "The list of directories to look in to find template include files. Separate with semicolons.")
	flag.Parse()
	files := flag.Args()

	if len(os.Args[1:]) == 0 {
		fmt.Println("got processes got template files, turning them into go code to use in your application.")
		fmt.Println("Usage: got [-o outDir] [-t fileType] [-i] [-I includeDirs] file1 [file2 ...] ")
		fmt.Println("-o: send processed files to the given directory. Otherwise sends to the same directory that the template is in.")
		fmt.Println("-t: process all files with this suffix in the current directory. Otherwise, specify specific files at the end.")
		fmt.Println("-i: run goimports on the result files to automatically fix up the import statement and format the file. You will need goimports installed.")
		fmt.Println("-I: the list of directories to search for include files. They are searched in the order given, and first one found will be used.")
	}


	got.IncludePaths = []string{"."}
	if includes != "" {
		i := strings.Split(includes, ";")
		for _, i2 := range i {
			f, _ := filepath.Abs(i2)
			got.IncludePaths = append(got.IncludePaths, f)
		}
	}

	if outDir != "" {
		dir,err := filepath.Abs(outDir)
		if err != nil {
			fmt.Println(err)
			return
		}
		outDir = dir
	}


	//var err error

	if typ != "" {
		files, _ = filepath.Glob("*." + typ)
	}


	for _, file := range files {
		s := processFile(file)
		if s != "" {
			writeFile(s, file, outDir, runImports)
		}
	}
}
