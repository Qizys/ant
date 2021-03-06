package router

import (
	"encoding/json"
	"github.com/anatasluo/ant/backend/engine"
	"github.com/anatasluo/ant/backend/setting"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func getSetting(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	WriteResponse(w, clientConfig.GetWebSetting())
}

func getStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	runningEngine.TorrentEngine.WriteStatus(w)
}

func getRunningQueue(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var tmp engine.TorrentLogsAndID
	runningEngine.TorrentDB.GetLogs(&tmp)
	WriteResponse(w, tmp)
}

func applySetting(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	isApplied := false
	var newSettings setting.WebSetting
	err := decoder.Decode(&newSettings)
	if err != nil {
		logger.WithFields(log.Fields{"Error": err}).Error("Failed to get new settings")
	}else{
		if runningEngine.EngineRunningInfo.HasRestarted == false {
			runningEngine.EngineRunningInfo.HasRestarted = true
			clientConfig.UpdateConfig(newSettings)
			logger.WithFields(log.Fields{"Settings": newSettings}).Info("Setting update")
			isApplied = true
			runningEngine.Restart()
			runningEngine.EngineRunningInfo.HasRestarted = false
		}
	}
	WriteResponse(w, JsonFormat{
		"IsApplied":isApplied,
	})
}

func handleSetting(router *httprouter.Router)  {
	router.GET("/settings/config", getSetting)
	router.GET("/settings/status", getStatus)
	router.GET("/settings/queue", getRunningQueue)
	router.POST("/settings/apply", applySetting)
}