package logrus

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/docopt/docopt.go"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus/hooks/source_file"
	"github.com/spf13/viper"
)

var (
	arguments map[string]interface{}
	appName   = ""
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
  level: debug
  format: text # text|json
  	timestamp: 2006-01-02 15:04:05
  	full_time: true
  hookfile: true
`)
)

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

	if viper.GetString("log.format") == "text" {
		customFormatter := new(TextFormatter)
		if viper.GetString("log.format.timestamp") != "" {
			customFormatter.TimestampFormat = viper.GetString("log.format.timestamp")
		}
		customFormatter.FullTimestamp = viper.GetBool("log.format.full_time")
		SetFormatter(customFormatter)
	} else {
		customFormatter := new(JSONFormatter)
		SetFormatter(customFormatter)
	}

	if viper.GetBool("log.hookfile") {
		if level, err := ParseLevel(viper.GetString("log.level")); err == nil {
			SetLevel(level)
			AddHook(&source_file.CodeLineHook{LogLevel: level})
			if level == DebugLevel {
				Debugln("logrus & viper init done")
			}
		} else {
			SetLevel(DebugLevel)
			AddHook(&source_file.CodeLineHook{LogLevel: DebugLevel})
		}
	} else {
		if level, err := ParseLevel(viper.GetString("log.level")); err == nil {
			SetLevel(level)
		} else {
			SetLevel(DebugLevel)
		}
	}
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

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
