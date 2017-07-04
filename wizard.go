package logrus

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/docopt/docopt.go"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	arguments map[string]interface{}
	vendor    = "/vendor/"
	appName   = ""
	goPath    = fmt.Sprintf("%s/src/", os.Getenv("GOPATH"))
	index     = len(goPath)
	usage     = `${APP}.

	Usage:
		${APP} [options]

	Options:
		-h --help                    Show this screen
		--conf=FILE                  Load a named configuration file
	`

	KGWF_config = []byte(`
kgwf:
  -
    raw_import:	https://gopkg.in/
    new_import: https://github.com/
    rule:  true
  -
    raw_import: https://go.googlesource.com/
    new_import: https://github.com/golang/
  -
    raw_import: https://golang.org/x/
    new_import: https://github.com/golang/
    split_sub: true
  -
    raw_import: https://google.golang.org/grpc
    new_import: https://github.com/grpc/grpc-go
  -
    raw_import: https://rsc.io/letsencrypt
    new_import: https://github.com/penhauer-xiao/letsencrypt

log:
  hookfile: true
  level: debug
`)
)

type sourceFileHook struct {
	LogLevel Level
}

type filePos struct {
	Pkg  string `json:"pkg"`
	File string `json:"file"`
	Func string `json:"func"`
	Line int    `json:"line"`
}

func (hook *sourceFileHook) Fire(entry *Entry) (_ error) {
	if pc, fullPath, line, ok := runtime.Caller(entry.Skip); ok {
		funcName := runtime.FuncForPC(pc).Name()
		relativePath := fullPath
		if temp := vendorPath(fullPath); temp != "" {
			relativePath = temp
		}

		pkg, file := path.Split(relativePath)
		if pkg != "" {
			pkg = pkg[:len(pkg)-1]
		}
		entry.Data["pos"] = filePos{
			Pkg:  pkg,
			File: file,
			Func: path.Base(funcName),
			Line: line,
		}
	}
	return
}

func (hook *sourceFileHook) Levels() []Level {
	levels := make([]Level, hook.LogLevel+1)
	for i, _ := range levels {
		levels[i] = Level(i)
	}
	return levels
}

func GetConfig() map[string]interface{} {
	return arguments
}

func init() {
	setAPPName(strings.Replace(os.Args[0], `./`, "", -1))

	file := "wizard.kgwf.yaml"
	if appName != "glide" {
		arguments = getCmdArguments()

		viper.SetConfigName("config")
	}

	if conf := arguments["--conf"]; conf != nil {
		file = conf.(string)
		if !exist(file) {
			Fatalf("No found configuration file : %s", file)
		}

		viper.AddConfigPath(path.Dir(file))
		viper.SetConfigName(strings.Split(path.Base(file), ".")[0])
		Infof("Using %s for configuration", file)
	} else if exist(file) {
		Infof("No configuration file defined, will try to use default (%s)", file)
		viper.AddConfigPath(".")
		viper.SetConfigName(file[:strings.LastIndex(file, ".")])

		err := viper.ReadInConfig()
		if err != nil {
			Fatalf("Fatal error config file: %s", err)
		}
	} else if appName == "glide" {
		Infof("No configuration file defined, Especially for glide configuration")
		viper.SetConfigType("yaml")
		err := viper.ReadConfig(bytes.NewBuffer(KGWF_config))
		if err != nil {
			Fatalf("Fatal error config file: %s", err)
		}
	} else {
		Infof("No configuration file defined, will try to use default (config.yaml)")
		viper.AddConfigPath(fmt.Sprintf("/etc/%s/", appName))
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil {
			Fatalf("Fatal error config file: %s", err)
		}
	}

	if appName != "glide" {
		arguments = cmdArguments(viper.GetString("server.version"))
		Infofp("The configure the command line parameters such as:", arguments)
	}

	if viper.GetBool("confg.watch") {
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			Infoln("Config file changed:", e.Name)
		})
	}

	SetFieldsLogger()
	SetFormatter(new(TextFormatter))

	if viper.GetBool("log.hookfile") {
		if level, err := ParseLevel(viper.GetString("log.level")); err == nil {
			SetLevel(level)
			AddHook(&sourceFileHook{LogLevel: level})
		} else {
			SetLevel(DebugLevel)
			AddHook(&sourceFileHook{LogLevel: DebugLevel})
		}
	} else {
		if level, err := ParseLevel(viper.GetString("log.level")); err == nil {
			SetLevel(level)
		} else {
			SetLevel(DebugLevel)
		}
	}
	Debugln("logrus & viper init done")
}

func setAPPName(name string) {
	appName = name
	usage = strings.Replace(usage, `${APP}`, name, -1)
}

func getCmdArguments() map[string]interface{} {
	arguments, err := docopt.Parse(usage, nil, true, "0.0", false)
	if err != nil {
		Warning("Error while parsing arguments: ", err)
	}

	return arguments
}

func cmdArguments(ver string) map[string]interface{} {
	arguments, err := docopt.Parse(usage, nil, true, ver, false)
	if err != nil {
		Warning("Error while parsing arguments: ", err)
	}

	return arguments
}

func vendorPath(fullPath string) string {
	if i := strings.Index(fullPath, vendor); i != -1 {
		return fullPath[i+len(vendor):]
	}
	return golangPath(fullPath)
}

func golangPath(fullPath string) string {
	if i := strings.Index(fullPath, goPath); i != -1 {
		return fullPath[i+index:]
	}
	return ""
}

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// fmt.Println(fmt.Sprintf("%s", stack()))
func stack() []byte {
	buf := make([]byte, 1<<20)
	n := runtime.Stack(buf, true)
	return buf[:n]
}
