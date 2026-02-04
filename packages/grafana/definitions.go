package grafana

type (
	logEntry struct {
		Stream map[string]string `json:"stream"`
		Values [][2]string       `json:"values"`
	}
	logsPayload struct {
		Streams []logEntry `json:"streams"`
	}

	LogMessage struct {
		Msg          string
		Channel      string
		LocationName string
		LocationId   string
	}
)
