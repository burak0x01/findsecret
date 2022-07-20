package main

import (
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

	var JsFile string
	var secrets []Secrets = readSecrets()
	var matchedSecrets []string

	println("\n[!] Working\n")

	checkSecret()

	var input string
	var output string
	flag.StringVar(&input, "i", "", "-i local.js \n-i https://domain.tld/external.js \n-i https://domain.tld")
	flag.StringVar(&output, "o", "cli", "-o cli \n-o <filename>.txt")
	flag.Parse()

	switch checkInput(input) {
	case "url":
		JsFile = getExternalJsFile(input)
		matchedSecrets = matchSecret(secrets, JsFile)
		checkOutput(input, output, matchedSecrets)
		println("[!] Completed\n")
	case "local":
		JsFile = getLocalJsFile(input)
		matchedSecrets = matchSecret(secrets, JsFile)
		checkOutput(input, output, matchedSecrets)
		println("[!] Completed\n")
	case "domain":
		JsFileList := getScripts(input)
		println("[!] Found", len(JsFileList), "JS file\n")
		for _, v := range JsFileList {
			JsFile = getExternalJsFile(v)
			matchedSecrets = matchSecret(secrets, JsFile)
			checkOutput(v, output, matchedSecrets)
		}
		println("[!] Completed\n")
	case "not valid":
		println("[!] unrecognized input\n")

	}
	// case "list":
	// 	JsFile = getExternalJsFile(input)
	// 	matchedSecrets = matchSecret(secrets, JsFile)

	// 	for _, url := range getJsList(input) {
	// 		checkOutput(url, output, matchedSecrets)
	// 	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkSecret() {
	homeDir, err := os.UserHomeDir()
	check(err)
	if _, err := os.Stat(homeDir + "/findsecret/secrets.json"); err != nil {
		downloadSecret()
	}
}

func checkInput(input string) string {

	if len(input) >= 4 {
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
				return "domain"
			}
		}
	} else {
		return "not valid"
	}

}

func checkOutput(input, output string, result []string) {

	if output == "cli" {
		writeCli(input, result)
	} else {
		writeFile(input, output, result)
	}
}

// func getJsList(input string) []string {
// 	lines := []string{}

// 	file, err := os.Open(input)
// 	check(err)
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	err2 := scanner.Err()
// 	check(err2)

// 	for scanner.Scan() {
// 		lines = append(lines, scanner.Text())
// 	}

// 	return lines
// }

func downloadSecret() {

	homeDir, err := os.UserHomeDir()
	check(err)

	err2 := os.Mkdir(homeDir+"/findsecret", 0777)
	check(err2)

	jsonFile, err3 := os.Create(homeDir + "/findsecret/secrets.json")
	check(err3)
	defer jsonFile.Close()

	resp, err4 := http.Get("https://raw.githubusercontent.com/burak0x01/findsecret/main/secrets.json")
	check(err4)
	defer resp.Body.Close()

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

func writeFile(input, path string, data []string) {

	if len(data) != 0 {
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		check(err)
		defer file.Close()

		file.WriteString(input + "\n")
		for _, line := range data {
			file.WriteString("\t" + line + "\n")
		}
		file.WriteString("\n")
	}
}

func writeCli(input string, data []string) {

	if len(data) != 0 {
		println("[+]", input)

		for _, v := range data {
			println("\t" + v)
		}
		println()
	}
}

func getScripts(url string) []string {

	var jsFileList []string

	resp, err := http.Get(url)
	check(err)

	if resp.StatusCode == 200 {
		body, err2 := io.ReadAll(resp.Body)
		check(err2)

		re := regexp.MustCompile(`https?://.*\.js`)
		match := re.FindAllString(string(body), -1)

		for _, url := range match {

			resp2, err3 := http.Get(url)
			if err3 != nil {
				continue
			}

			if resp2.StatusCode == 200 {
				jsFileList = append(jsFileList, url)
			}
		}

		return jsFileList
	} else {
		jsFileList = append(jsFileList, "check the url please !!!")

		return jsFileList
	}
}
