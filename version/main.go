package version

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
)

const BASE_URL = "https://version.yyjlincoln.app"

func GetVersion(name string, channel string) (*VersionSpec, error) {
	resp, err := http.Get(BASE_URL + "/" + name + "/" + channel)
	if err != nil {
		return nil, fmt.Errorf("could not get version: %v", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %v", err)
	}
	versionSpec := &VersionSpec{}
	err = json.Unmarshal(data, versionSpec)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %v", err)
	}
	return versionSpec, nil
}

func PrintVersionInformation(v VersionSpec) {
	fmt.Printf("Version: ")
	Successf("%s\n", v.Version)
	fmt.Printf("Channel: ")
	Successf("%s\n", v.Channel)
	fmt.Printf("Released At: ")
	Successf("%s\n", v.ReleasedAt)
	fmt.Printf("Description: ")
	Successf("%s\n", v.Description)
	fmt.Printf("URL: ")
	Successf("%s\n", v.URL)
}

func Errorf(format string, v ...interface{}) {
	if len(format) == 0 {
		return
	}
	hasTrailingNewLine := false

	if format[len(format)-1] == '\n' {
		format = format[:len(format)-1]
		hasTrailingNewLine = true
	}

	coloredErrorString := color.RedString(fmt.Sprintf(format, v...))

	if hasTrailingNewLine {
		fmt.Printf("%v\n", coloredErrorString)
	} else {
		fmt.Printf("%v", coloredErrorString)
	}
}

func Successf(format string, v ...interface{}) {
	if len(format) == 0 {
		return
	}
	hasTrailingNewLine := false

	if format[len(format)-1] == '\n' {
		format = format[:len(format)-1]
		hasTrailingNewLine = true
	}

	coloredErrorString := color.GreenString(fmt.Sprintf(format, v...))

	if hasTrailingNewLine {
		fmt.Printf("%v\n", coloredErrorString)
	} else {
		fmt.Printf("%v", coloredErrorString)
	}
}

func Warnf(format string, v ...interface{}) {
	if len(format) == 0 {
		return
	}
	hasTrailingNewLine := false

	if format[len(format)-1] == '\n' {
		format = format[:len(format)-1]
		hasTrailingNewLine = true
	}

	coloredErrorString := color.YellowString(fmt.Sprintf(format, v...))

	if hasTrailingNewLine {
		fmt.Printf("%v\n", coloredErrorString)
	} else {
		fmt.Printf("%v", coloredErrorString)
	}
}

func Infof(format string, v ...interface{}) {
	if len(format) == 0 {
		return
	}
	hasTrailingNewLine := false

	if format[len(format)-1] == '\n' {
		format = format[:len(format)-1]
		hasTrailingNewLine = true
	}

	coloredErrorString := color.BlueString(fmt.Sprintf(format, v...))

	if hasTrailingNewLine {
		fmt.Printf("%v\n", coloredErrorString)
	} else {
		fmt.Printf("%v", coloredErrorString)
	}
}
