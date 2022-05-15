package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

const (
	FILE_EXECUTE_PERMISSION   = 0664
	FOLDER_EXECUTE_PERMISSION = 0775
)

func main() {

	tsPath := "src"

	// this slice will be used as a value for the 'files' field in package.json
	tsFileList := []string{"esm", "index.js", "index.d.ts"}

	files, err := ioutil.ReadDir(tsPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {

		if f.IsDir() {
			tsFileList = append(tsFileList, f.Name())
		}
	}

	// populate 'files' field value with tsFileList in the packagae.json
	ore, err := ioutil.ReadFile("package.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m := map[string]interface{}{}

	err = json.Unmarshal(ore, &m)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m["files"] = tsFileList

	newPackage, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pack, err := os.OpenFile("package.json", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, FILE_EXECUTE_PERMISSION)
	if err != nil {
		fmt.Println("Failed to open package.json", err)
		os.Exit(1)
	}

	defer pack.Close()

	_, err = pack.Write(newPackage)
	if err != nil {
		fmt.Println("Failed to write to package.json", err)
		os.Exit(1)
	}

	// setting up .npmrc file with authToken
	npmrc, err := os.OpenFile(".npmrc", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, FILE_EXECUTE_PERMISSION)
	if err != nil {
		fmt.Println("Failed to open npmrc", err)
		os.Exit(1)
	}

	defer npmrc.Close()

	if len(os.Getenv("NPM_TOKEN")) == 0 {
		fmt.Println("No NPM_TOKEN env found")
		os.Exit(1)
	}

	if _, err = npmrc.WriteString("//registry.npmjs.org/:_authToken=" + os.Getenv("NPM_TOKEN")); err != nil {
		fmt.Println("Failed to open npmrc", err)
		os.Exit(1)
	}

	// get latest published version
	getVersions := exec.Command("npm", "show", "m3o", "--time", "--json")
	getVersions.Dir = tsPath

	outp, err := getVersions.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to get versions of NPM package", string(outp))
		os.Exit(1)
	}

	type npmVers struct {
		Versions []string `json:"versions"`
	}

	npmOutput := &npmVers{}

	if len(outp) > 0 {
		err = json.Unmarshal(outp, npmOutput)
		if err != nil {
			fmt.Println("Failed to unmarshal versions", string(outp))
			os.Exit(1)
		}
	}

	fmt.Println("the last published version: ", npmOutput.Versions[len(npmOutput.Versions)-1])
}
