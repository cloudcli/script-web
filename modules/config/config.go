package config

import (
	"gopkg.in/ini.v1"
	"runtime"
	"os/exec"
	"os"
	log "github.com/Sirupsen/logrus"
	"path/filepath"
	"errors"
	"path"
	"strings"
)

var (
	// Is os windows
	IsWindows   bool

	// App version
	AppVersion = "0.0.1"

	// App Name
	AppName    = "Cloudcli"

	// User Home dir
	HomeDir     string

	// App binary path
	AppPath     string

	// App user data path: default "~/.Cloudcli/"
	AppUsrPath  string

	// App log dir
	LogDir      string

	// APp log level
	LogLevel    int

	// Custom config path
	ConfigPath  string

	// Config
	Config      *ini.File

	MysqlDatabase  string	// Mysql database name
	MysqlAddress   string	// Mysql address
	MysqlUser      string	// Mysql username
	MysqlPassword  string	// Mysql password

	RepoPath       string   // Script git repo path
)

func init() {
	IsWindows = runtime.GOOS == "windows"

	AppPath, err := getAppPath()
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
	AppPath = strings.Replace(AppPath, "\\", "/", -1)
	ConfigPath = path.Join(path.Dir(AppPath), "app.ini")

	HomeDir, err := getHomeDir()
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}

	AppUsrPath = path.Join(HomeDir, "." + AppName)
	LogDir = path.Join(AppUsrPath, "log")

	if _, err = os.Stat(AppUsrPath); os.IsNotExist(err) {
		err = os.MkdirAll(AppUsrPath, 0777)
		if err != nil {
			log.Fatal(err)
			os.Exit(0)
		}
	}

	if _, err = os.Stat(LogDir); os.IsNotExist(err) {
		err = os.MkdirAll(LogDir, 0777)
		if err != nil {
			log.Fatal(err)
			os.Exit(0)
		}
	}
}

/**
 * Intialize configuration from config file (.ini), if configFilePath not provide ,use the default configfile
 * which is (AppWor)
 */
func InitConfig(configFilePath string) {
	if configFilePath == "" {
		configFilePath = ConfigPath
	} else {
		ConfigPath = configFilePath
	}

	cfg, err := ini.InsensitiveLoad(configFilePath)
	if err != nil {
		log.Fatalf("Can't find config file in %s with err: %v", configFilePath, err)
		os.Exit(0)
	}

	// parse database section

	sec, err := cfg.GetSection("database")
	if err != nil {
		log.Fatal("Can't find database config section")
		os.Exit(0)
	}

	key, err := sec.GetKey("address")
	if err != nil {
		log.Fatalf("Parse Database Address Error: %v", err)
		os.Exit(0)
	}
	MysqlAddress = key.String()

	key, err = sec.GetKey("db")
	if err != nil {
		log.Fatalf("Parse db Error: %v", err)
		os.Exit(0)
	}
	MysqlDatabase = key.String()

	key, err = sec.GetKey("user")
	if err != nil {
		log.Fatalf("Parse user Error: %v", err)
		os.Exit(0)
	}
	MysqlUser = key.String()

	key, err = sec.GetKey("password")
	if err != nil {
		log.Fatalf("Parse password Error: %v", err)
		os.Exit(0)
	}
	MysqlPassword = key.String()

	// parse repo section
	sec, err = cfg.GetSection("repo")
	if err != nil {
		log.Fatal("Can't find repo config section")
		os.Exit(0)
	}

	key, err = sec.GetKey("path")

	if err != nil {
		log.Fatalf("Parse repo path Error: %v", err)
		os.Exit(0)
	}

	RepoPath = key.String()
	RepoPath, _ = filepath.Abs(RepoPath)
}

func getAppPath() (string, error){
	file, err := exec.LookPath(os.Args[0])

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return filepath.Abs(file)
}

func getHomeDir() (home string, err error) {
	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
		if len(home) == 0 {
			home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		}
	} else {
		home = os.Getenv("HOME")
	}

	if len(home) == 0 {
		return "", errors.New("Cannot specify home directory because it's empty")
	}

	return home, nil
}


func forcePathSeparator(path string) {
	if strings.Contains(path, "\\") {
		log.Fatal(4, "Do not use '\\' or '\\\\' in paths, instead, please use '/' in all places")
	}
}