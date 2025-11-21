package profile

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Profile struct {
	ENV       string   //当前环境
	FilePaths []string //配置文件路径
	FileType  string   //配置文件后缀
}

var (
	PATH           = "configs"
	PATH_SEPARATOR = os.PathSeparator

	FILE_BASE_NAME = "app"

	FIEL_SEPARATOR = "-"

	FILE_TYPE = "yaml"
)

func LoadProfile() Profile {
	flag.Parse()
	var pro Profile
	pro.FileType = FILE_TYPE
	pro.ENV = os.Getenv("APP_ENV")
	flag.StringVar(&pro.ENV, "env", "", "runtime environment, eg: -env remote")

	if len(pro.ENV) > 0 {
		fmt.Printf("Using ENV %s\n", pro.ENV)
	}

	filePaths := make([]string, 0, 6)
	var builder strings.Builder
	builder.WriteString(PATH)
	builder.WriteByte(byte(PATH_SEPARATOR))
	builder.WriteString(FILE_BASE_NAME)
	builder.WriteString(".")
	builder.WriteString(FILE_TYPE)

	filePaths = append(filePaths, builder.String())
	if len(pro.ENV) > 0 {
		builder.Reset()
		builder.WriteString(PATH)
		builder.WriteByte(byte(PATH_SEPARATOR))
		builder.WriteString(FILE_BASE_NAME)
		builder.WriteString(FIEL_SEPARATOR)
		builder.WriteString(pro.ENV)
		builder.WriteString(".")
		builder.WriteString(FILE_TYPE)
		filePaths = append(filePaths, builder.String())
	}
	pro.FilePaths = filePaths
	return pro

}
