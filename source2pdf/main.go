package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

var flagNoGen = flag.Bool("f", false, "fast mode; don't run XeLaTeX")

type TemplateData struct {
	Sources []Source
}

type Source struct {
	Language string
	Path     string
}

func hasExt(ext string) ignoreFunc {
	return func(path string) bool {
		return strings.HasSuffix(path, ext)
	}
}

func hasNoExt(except ...string) ignoreFunc {
	return func(path string) bool {
		for _, exception := range except {
			if filepath.Base(path) == exception {
				return false
			}
		}
		return filepath.Ext(path) == ""
	}
}

func exact(path string) ignoreFunc {
	return func(path2 string) bool {
		return path == path2
	}
}

func lit(filename string) ignoreFunc {
	return func(path string) bool {
		return filepath.Base(path) == filename
	}
}

func main() {
	flag.Parse()

	err := os.MkdirAll("build", 0o755)
	if err != nil {
		fmt.Printf("Error: could not create build directory: %s\n", err)
		os.Exit(1)
	}

	tpl := template.Must(
		template.
			New("main.tex.tpl").
			Funcs(map[string]interface{}{
				"texEscapeFull": func(s string) string {
					r := strings.NewReplacer(`%`, `\%`, `_`, `\_`, `&`, `\&`, `#`, `\#`, `\`, `\\`, `{`, `\{`)
					return r.Replace(s)
				},
				"texEscapeMintedPath": func(s string) string {
					r := strings.NewReplacer(`%`, `\%`, `\`, `\\`, `{`, `\{`)
					return r.Replace(s)
				},
			}).
			ParseFiles("main.tex.tpl"))
	sources, err := getSources(ignore{
		dirs: []ignoreFunc{
			exact("../mfrc522/docs/api"),
			exact("../mfrc522/docs/ext"),
			exact("../mfrc522/docs/doxyoutput"),
			exact("../mfrc522/docs/_build"),
			exact("../source-pdf"),
			exact("../docs/themes"),
			exact("../docs/public"),
			prepareWildcard("../os/build-*"),
			exact("../os/sstate-cache"),
			exact("../os/downloads"),
			func(s string) bool {
				if strings.HasPrefix(s, "../os/sources") {
					return !strings.HasPrefix(s, "../os/sources/meta-timeterm")
				}
				return false
			},
			prepareWildcard("*/*build-*"),
			lit("3rdparty"),
			lit("assets"),
			lit("build"),
			lit("node_modules"),
		},
		files: []ignoreFunc{
			exact("../design/js/CCapture.js"),
			exact("../design/js/download.js"),
			exact("../design/js/p5.js"),
			exact("../design/js/webm-writer-0.2.0.js"),
			exact("../mfrc522/src/mfrc522.cpp"),
			exact("../mfrc522/include/mfrc522.h"),
			exact("../os/setup-environment.sh"),
			hasExt(".code-workspace"),
			hasExt(".ico"),
			hasExt(".json"),
			hasExt(".lock"),
			hasExt(".log"),
			hasExt(".mtl"),
			hasExt(".obj"),
			hasExt(".pb.go"),
			hasExt(".png"),
			hasExt(".sum"),
			hasExt(".jpg"),
			hasExt(".jpeg"),
			hasExt(".svg"),
			hasExt(".ttf"),
			hasExt(".xml"),
			lit(".eslintcache"),
			lit(".env"),
			lit(".gitignore"),
			lit("CMakeLists.txt.user"),
			hasNoExt("Dockerfile"),
		},
	})
	if err != nil {
		fmt.Printf("Error: could not get sources: %s\n", err)
		os.Exit(1)
	}

	out, err := os.Create("build/main.tex")
	if err != nil {
		fmt.Printf("Error: could not create main.tex: %s\n", err)
		os.Exit(1)
	}
	defer out.Close()

	err = tpl.Execute(out, TemplateData{
		Sources: sources,
	})
	if err != nil {
		fmt.Printf("Error: could not execute template: %s\n", err)
		os.Exit(1)
	}

	if !*flagNoGen {
		cmd := exec.Command("xelatex", "-shell-escape", "main.tex")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = "build"
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error: failed to create PDF: %s\n", err)
			os.Exit(1)
		}
	}
}

type ignore struct {
	dirs  []ignoreFunc
	files []ignoreFunc
}

type ignoreFunc func(s string) bool

func getSources(ignore ignore) ([]Source, error) {
	var sources []Source
	if err := filepath.Walk("..", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if strings.HasPrefix(filepath.Base(info.Name()), ".") &&
				strings.Trim(info.Name(), ".") != "" {
				return filepath.SkipDir
			}
			for _, f := range ignore.dirs {
				if f(path) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		if info.Mode()&fs.ModeSymlink != 0 {
			return nil
		}
		for _, f := range ignore.files {
			if f(path) {
				return nil
			}
		}
		path = path[3:]
		fmt.Println(path)
		sources = append(sources, Source{
			Language: langByFilename(info.Name()),
			Path:     filepath.Join(path),
		})
		return nil
	}); err != nil {
		return nil, err
	}
	return sources, nil
}

func langByFilename(filename string) string {
	switch filename {
	case ".clang-format":
		return "yaml"
	case "CMakeLists.txt":
		return "cmake"
	case "Dockerfile":
		return "docker"
	case "frontend-embedded_nl_NL.ts":
		return "xml"
	case "nginx.conf":
		return "nginx"
	}

	if strings.HasPrefix(filename, ".env.") {
		return "ini"
	}

	switch filepath.Ext(filename) {
	case ".bat", ".cmd":
		return "bat"
	case ".cpp", ".h", ".hpp":
		return "cpp"
	case ".css":
		return "css"
	case ".Dockerfile":
		return "docker"
	case ".go":
		return "go"
	case ".html":
		return "html"
	case ".js", ".jsx":
		return "js"
	case ".md":
		return "md"
	case ".proto":
		return "proto"
	case ".qml":
		return "qml"
	case ".qrc":
		return "xml"
	case ".sh":
		return "bash"
	case ".sql":
		return "postgresql"
	case ".ts", ".tsx":
		return "typescript"
	case ".yaml", ".yml":
		return "yaml"
	}

	return "text"
}
