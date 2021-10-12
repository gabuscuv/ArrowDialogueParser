package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {

	appConfig := loadConfig("./config.json")
	var JsonPath string

	if len(os.Args) > 1 {
		JsonPath = os.Args[1]
	} else {
		JsonPath = appConfig.MainCon.DefaultJson
	}

	parsedCSV := ParseJSON(JsonPath)

	file, err := os.Create(appConfig.MainCon.OutputPath + appConfig.MainCon.OutputFile)
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.WriteAll(parsedCSV)
}

func ParseJSON(JsonPath string) [][]string {
	outputcsv := [][]string{
		{"ArrowID", "DialogueWaveId", "Voice", "Location", "Text"},
	}

	var tmp_DialogueOwner string

	var data map[string]interface{}
	raw, err := ioutil.ReadFile(JsonPath)
	checkError("Failed trying to read the json file: "+JsonPath, err)
	err = json.Unmarshal(raw, &data)
	checkError("Failed trying to unmarshall the json file: "+JsonPath, err)
	for _, node := range data["resources"].(map[string]interface{})["nodes"].(map[string]interface{}) {
		if node.(map[string]interface{})["type"] == "dialog" {
			tmp_DialogueOwner = resolveDialogueOwnerNames(data, node)
			for _, line := range node.(map[string]interface{})["data"].(map[string]interface{})["lines"].([]interface{}) {
				// TODO: Find the Node Place.
				outputcsv = append(
					outputcsv, []string{
						fmt.Sprintf("%v", node.(map[string]interface{})["name"]),
						strings.Split(fmt.Sprintf("%v", node.(map[string]interface{})["notes"]), "\n")[0],
						tmp_DialogueOwner,
						"",
						fmt.Sprintf("%v", line)})
			}
		}
	}

	return outputcsv
}

func resolveDialogueOwnerNames(data map[string]interface{}, node interface{}) string {

	characterId := node.(map[string]interface{})["data"].(map[string]interface{})["character"].(float64)

	if characterId == -1 {
		return "Anonymous"
	}

	return fmt.Sprintf("%v", data["resources"].(map[string]interface{})["characters"].(map[string]interface{})[fmt.Sprintf("%v", characterId)].(map[string]interface{})["name"])

}

/*func resolvePlaceName(data map[string]interface{}, node interface{}) string {

	return fmt.Sprintf("%v", data["resources"].(map[string]interface{})["characters"].(map[string]interface{})[fmt.Sprintf("%v", node)].(map[string]interface{})["name"])

}
*/
func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func loadConfig(file string) appConfig {
	var config appConfig

	resp, err := ioutil.ReadFile(file)
	if err == nil {
		err := json.Unmarshal(resp, &config)
		if err == nil {
			return config
		}
	}
	return appConfig{}
}

type appConfig struct {
	MainCon MainConfig `json:"MainConfig"`
}
type MainConfig struct {
	DefaultJson string `json:"defaultJSON"`
	OutputPath  string `json:"outputPath"`
	OutputFile  string `json:"outputFile"`
}
