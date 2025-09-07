package googleads

import "client-runaway-zenoti/internal/config"

func getKeyData() []byte {
	return []byte(config.Confs.GAjson)
}
