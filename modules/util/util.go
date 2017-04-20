package util

import (
	"os"
	"runtime"
	"errors"
	"strings"
	"crypto/md5"
	"encoding/hex"
	"reflect"
	"fmt"
	"io/ioutil"
)

func IsFile(filePath string) bool{
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

func HomeDir() (home string, err error) {
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

func CurrentUsername() string {
	curUserName := os.Getenv("USER")
	if len(curUserName) > 0 {
		return curUserName
	}

	return os.Getenv("USERNAME")
}

func PWD() string {
	dir, _ := os.Getwd()
	return strings.Replace(dir, " ", "\\ ", -1)
}

func EncodePassword(password string) (string, error) {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(password))
	cipherStr := md5Ctx.Sum(nil)
	result := strings.ToLower(hex.EncodeToString(cipherStr))
	return result, nil
}

// ToInt64 convert any numeric value to int64
func toInt64(value interface{}) (d int64, err error) {
	val := reflect.ValueOf(value)
	switch value.(type) {
	case int, int8, int16, int32, int64:
		d = val.Int()
	case uint, uint8, uint16, uint32, uint64:
		d = int64(val.Uint())
	default:
		err = fmt.Errorf("ToInt64 need numeric not `%T`", value)
	}
	return
}

func ReadFile(filePath string) ([]byte, error) {
	var b []byte
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return b, err
	}

	b, err := ioutil.ReadFile(filePath)

	if err != nil {
		return b, err
	}

	return b, nil
}