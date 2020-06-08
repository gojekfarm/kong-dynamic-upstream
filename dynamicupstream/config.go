package dynamicupstream

import (
	"encoding/json"
	"fmt"
)

func RouteConfig(route, upstreamConfig string) (string, string) {
	configMap := make(map[string]interface{})

	json.Unmarshal([]byte(upstreamConfig), &configMap)
	routeConfig := configMap[route].(map[string]interface{})
	return fmt.Sprintf("%v", routeConfig["expression"]), fmt.Sprintf("%v", routeConfig["port"])
}
