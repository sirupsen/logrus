package logrus

import (
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
		if temp := vendor(fullPath); temp != "" {
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
	arguments = getCmdArguments()
	Infofp("The configure the command line parameters such as:", arguments)

	file := ""
	viper.SetConfigName("config")
	if conf := arguments["--conf"]; conf != nil {
		file = conf.(string)
		if !exist(file) {
			Fatalf("No found configuration file : %s", file)
		}

		viper.AddConfigPath(path.Dir(file))
		viper.SetConfigName(strings.Split(path.Base(file), ".")[0])
		Infof("Using %s for configuration", file)
	} else {
		Infof("No configuration file defined, will try to use default (config.yaml)")
	}

	viper.SetConfigType("yaml")
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", appName))
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		Fatalf("Fatal error config file: %s", err)
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
	arguments, err := docopt.Parse(usage, nil, true, VERSION, false)
	if err != nil {
		Warning("Error while parsing arguments: ", err)
	}

	return arguments
}

func vendor(fullPath string) string {
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
