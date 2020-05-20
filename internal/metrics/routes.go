package metrics

import (
	"encoding/json"
	"net/http"

	"github.com/lbryio/lbrytv/internal/errors"
	"github.com/lbryio/lbrytv/internal/monitor"
	"github.com/lbryio/lbrytv/internal/responses"

	"github.com/spf13/cast"
)

var Logger = monitor.NewModuleLogger("metrics")

func TrackUIMetric(w http.ResponseWriter, req *http.Request) {
	responses.AddJSONContentType(w)
	resp := make(map[string]string)
	code := http.StatusOK

	metricName := req.FormValue("name")
	resp["name"] = metricName

	switch metricName {
	case "buffer":
		UIBufferCount.Inc()
	case "time_to_start":
		UITimeToStart.Observe(cast.ToFloat64(req.FormValue("value")))
	default:
		Logger.Log().Errorf("invalid UI metric name: %s", metricName)
		code = http.StatusBadRequest
		resp["error"] = "Invalid metric name"
	}

	if errMsg, ok := resp["error"]; ok {
		monitor.ErrorToSentry(errors.Err(errMsg))
	}

	w.WriteHeader(code)
	respByte, _ := json.Marshal(&resp)
	w.Write(respByte)
}
