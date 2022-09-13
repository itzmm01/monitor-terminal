package jobs

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// Include aa
type Include struct {
	File   string
	Params map[string]string `json:"params"`
}

// ReadJobYml read job.yml
func ReadJobYml(filePath string, jobParams map[string]string, allJobRes map[string]map[string][][]string) {
	File, err := ioutil.ReadFile(filePath)
	if err != nil {
		Console.ErrorFile("YML file not found: " + filePath)
		return
	}
	var JobYmlList interface{}
	err = yaml.Unmarshal(File, &JobYmlList)
	if err != nil {
		Console.ErrorFile("YML file format error: " + filePath)
		return
	}
	if JobYmlList == nil {
		Console.ErrorFile("JobYml is nul: " + filePath)
		return
	}
	for _, jobYml := range JobYmlList.([]interface{}) {
		for key, obj := range jobYml.(map[string]interface{}) {
			if key == "job" {
				resByte, _ := json.Marshal(obj)
				var jobItem Job
				err := json.Unmarshal(resByte, &jobItem)
				if err != nil {
					Console.ErrorFile("YML file format error: job Unmarshal failed-" + jobItem.Name)
					return
				}
				jobItem.Params = jobParams
				jobRes := jobItem.RunJob()
				allJobRes[jobItem.Name] = jobRes
			} else if key == "include" {
				resByte, _ := json.Marshal(obj)
				var includeList []Include
				err := json.Unmarshal(resByte, &includeList)
				if err != nil {
					Console.ErrorFile("YML file format error: include Unmarshal failed")
					return
				}
				for _, include := range includeList {
					ReadJobYml(include.File, include.Params, allJobRes)
				}
			} else {
				Console.ErrorFile("YML file format error: unknown key " + key)
			}
		}
	}
}
