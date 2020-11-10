package rest

import "encoding/json"

func HeaderFormat(header string) (result map[string][]string, err error) {
	if header != "" {
		err = json.Unmarshal([]byte(header), &result)
	}
	return
}
