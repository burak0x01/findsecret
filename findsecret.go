package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"os"
	"regexp"
)

type Secrets struct {
	Name     string   `json:"name"`
	Patterns []string `json:"patterns"`
}

func main() {
	// flags
	var input string
	var output string
	flag.StringVar(&input, "i", "", "-i something.js \n-i https://domain.com/something.js \n-i domains.txt")
	flag.StringVar(&output, "o", "cli", "-o cli \n-o <filename>.txt")
	flag.Parse()

	// check secrets.json file exist
	homeDir, err := os.UserHomeDir()
	check(err)
	if _, err := os.Stat(homeDir + "/findsecret/secrets.json"); err != nil {
		downloadSecret()
	}

	switch checkInput(input) {
	case "url":
		checkOutput(input, output, matchSecret(readSecrets(), getExternalJsFile(input)))
	case "local":
		checkOutput(input, output, matchSecret(readSecrets(), getLocalJsFile(input)))
	case "list":
		for _, url := range getJsList(input) {
			checkOutput(url, output, matchSecret(readSecrets(), getExternalJsFile(url)))
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkInput(input string) string {

	isJS, err := regexp.MatchString(`^.*\.js(\?.*=?.*)?$`, input)
	check(err)

	isURL, err := regexp.MatchString(`https?://.*`, input)
	check(err)

	isTxt, err := regexp.MatchString(`^.*\.txt$`, input)
	check(err)

	if isJS {
		if isURL {
			return "url"
		} else {
			return "local"
		}
	} else {
		if isTxt {
			return "list"
		} else {
			return "Please check input"
		}
	}
}

func checkOutput(input, output string, result []string) {

	if output == "cli" {
		writeCli(input, result)
	} else {
		writeFile(input, output, result)
	}
}

func getJsList(input string) []string {
	lines := []string{}

	file, err := os.Open(input)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	err2 := scanner.Err()
	check(err2)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

// Download secrets.json
func downloadSecret() {

	// get user home directory
	homeDir, err := os.UserHomeDir()
	check(err)

	// make findsecret directory
	err2 := os.Mkdir(homeDir+"/findsecret", 0777)
	check(err2)

	// create secrets.json
	jsonFile, err3 := os.Create(homeDir + "/findsecret/secrets.json")
	check(err3)
	defer jsonFile.Close()

	// get external data
	resp, err4 := http.Get("https://raw.githubusercontent.com/burak0x01/findsecret/main/secrets.json")
	check(err4)
	defer resp.Body.Close()

	// write response data to the file
	_, err5 := io.Copy(jsonFile, resp.Body)
	check(err5)
}

func readSecrets() []Secrets {

	var secrets []Secrets

	homeDir, err := os.UserHomeDir()
	check(err)

	jsonFile, err := os.Open(homeDir + "/findsecret/secrets.json")
	check(err)
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &secrets)

	return secrets
}

// match secrets and return array
func matchSecret(secrets []Secrets, data string) []string {

	matchedArray := []string{}

	// each secret object
	for _, secret := range secrets {

		// each pattern
		for _, pattern := range secret.Patterns {
			re := regexp.MustCompile(pattern)
			match := re.FindAllString(data, -1)

			// save multiple matches
			for _, matchVal := range match {
				if matchVal != "" {
					s := secret.Name + " => " + matchVal
					matchedArray = append(matchedArray, s)
				}
			}
		}
	}

	return matchedArray
}

func getExternalJsFile(url string) string {

	resp, err := http.Get(url)

	if err == nil {
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			body, err := io.ReadAll(resp.Body)
			check(err)
			return string(body)
		} else {
			return "Not Found"
		}
	} else {
		return "Not Found"
	}
}

func getLocalJsFile(input string) string {

	data, err := os.ReadFile(input)
	check(err)

	return string(data)
}

// her data için yeniden dosya yaratılıp içine yazılıyor
func writeFile(input, path string, data []string) {

	if len(data) != 0 {
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		check(err)
		defer file.Close()

		file.WriteString(input)
		for _, line := range data {
			file.WriteString("\t" + line + "\n")
		}
	}
}

func writeCli(input string, data []string) {

	if len(data) != 0 {
		println(input)

		for _, v := range data {
			println("\t" + v)
		}
	}
}
