package log

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"strconv"

	"github.com/xanderxampp-be/franco/log/entity"

	"github.com/sirupsen/logrus"
)

const (
	httpRequest      = "REQUEST"
	httpResponse     = "RESPONSE"
	folder           = "logs"
	timeformat       = "2006-01-02T15:04:05-0700"
	nameformat       = "log-2006-01-02.log"
	nameformatTrxLog = "trxlog-2006-01-02.log"
)

var (
	currentFileName string
	currentFile     *os.File
	logText         *logrus.Logger
	logJSON         *logrus.Logger
	serviceName     string
	debug           bool
	err             error
)

func Init() {
	serviceName = os.Getenv("MICRO_NAME")
	setText()
	setJSON()
	setFolder()

	debug = false

	debugStr := os.Getenv("DEBUG")
	debug, err = strconv.ParseBool(debugStr)

	if err != nil {
		fmt.Println(err)
	}

	if debug {
		logText.SetLevel(logrus.DebugLevel)
		logJSON.SetLevel(logrus.DebugLevel)
	} else {
		logText.SetLevel(logrus.InfoLevel)
		logJSON.SetLevel(logrus.InfoLevel)
	}
}

func setFolder() {
	dir, _ := os.Getwd()
	folderlogs := dir + "/" + folder

	if _, err := os.Stat(folderlogs); os.IsNotExist(err) {
		err := os.Mkdir(folderlogs, 0777)
		// TODO: handle error
		fmt.Println(err)
	}
	/*
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err := os.Mkdir(path, mode)
			// TODO: handle error
		}
	*/
}

func setJSON() {
	logJSON = logrus.New()
	formatter := new(logrus.JSONFormatter)
	formatter.DisableTimestamp = true
	logJSON.SetFormatter(formatter)
}

func setText() {
	logText = logrus.New()
	formatter := new(logrus.TextFormatter)
	formatter.DisableTimestamp = true
	formatter.DisableQuote = true
	logText.SetFormatter(formatter)
}

func setLogFile(mode int) string {
	currentTime := time.Now()
	timestamp := currentTime.Format(timeformat)

	fileFormat := nameformat

	if mode == 1 {
		fileFormat = nameformatTrxLog
	}

	filename := folder + "/" + currentTime.Format(fileFormat)
	if filename == currentFileName {
		// not changing date, therefore keep using the same logfile
		return timestamp
	}

	// changing date in which leads to different file name
	newLogFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	} else {
		// change the current file name to the new file name
		currentFileName = filename
		logText.SetOutput(newLogFile)
		logJSON.SetOutput(newLogFile)

		// close the old file
		if currentFile != nil {
			currentFile.Close()
			currentFile = newLogFile
		}
	}

	return timestamp
}

func LogReq(trxType string, body *interface{}, header *http.Header) {
	// setJSON()
	timestamp := setLogFile(0)
	mapRequest := Minify(body)
	mapRequest = getCleanMap(mapRequest)

	logJSON.WithFields(logrus.Fields{
		"service":         serviceName,
		"http_type":       httpRequest,
		"request_header":  header,
		"request_body":    mapRequest,
		"trx_type":        trxType,
		"device_id":       header.Get("Device-Id"),
		"device_type":     header.Get("Device-Type"),
		"device_name":     header.Get("Device-Name"),
		"device_version":  header.Get("Device-Version"),
		"device_sequence": header.Get("Device-Sequence"),
		"version":         header.Get("Version"),
		"type_user":       header.Get("Type-User"),
		"timestamp":       timestamp,
	}).Info("REQUEST")

}

func LogRespNonfin(param *entity.Responselog) {
	timestamp := setLogFile(0)
	mapResponse := Minify(param.ResponseBody)
	logJSON.WithFields(logrus.Fields{
		"username":        param.Username,
		"third_party":     param.ThirdParty,
		"service":         serviceName,
		"http_type":       httpResponse,
		"response_header": param.ResponseHeader,
		"response_body":   mapResponse,
		"trx_type":        param.TrxType,
		"device_id":       param.DeviceId,
		"device_type":     param.DeviceType,
		"device_name":     param.DeviceName,
		"device_version":  param.DeviceVersion,
		"response_code":   param.ResponseCode,
		"trace":           param.Trace,
		"timestamp":       timestamp,
		"elapsed":         param.Elapsed,
		"fast_menu":       param.FastMenuFlag,
		"type_user":       param.TypeUser,
	}).Info(httpResponse)
}

func LogRespFin(param *entity.Responselog) {
	timestamp := setLogFile(0)
	mapResponse := Minify(param.ResponseBody)
	logJSON.WithFields(logrus.Fields{
		"username":           param.Username,
		"account_debet":      param.AccountDebet,
		"amount":             param.Amount,
		"amount_float":       param.AmountFloat,
		"fee":                param.Fee,
		"transaction_refnum": param.TransactionRefnum,
		"third_party":        param.ThirdParty,
		"service":            serviceName,
		"http_type":          httpResponse,
		"response_header":    param.ResponseHeader,
		"response_body":      mapResponse,
		"trx_type":           param.TrxType,
		"device_id":          param.DeviceId,
		"device_type":        param.DeviceType,
		"device_name":        param.DeviceName,
		"device_version":     param.DeviceVersion,
		"response_code":      param.ResponseCode,
		"trace":              param.Trace,
		"timestamp":          timestamp,
		"elapsed":            param.Elapsed,
		"fast_menu":          param.FastMenuFlag,
		"type_user":          param.TypeUser,
	}).Info(httpResponse)
}

func LogTrxLog(param *entity.TrxLog) {
	_ = setLogFile(1)
	logJSON.WithFields(logrus.Fields{
		"id":                param.Id,
		"username":          param.Username,
		"account":           param.Account,
		"reference_num":     param.ReferenceNum,
		"logged":            param.Logged,
		"trx_type":          param.TrxType,
		"trx_status":        param.TrxStatus,
		"trx_object":        param.TrxObject,
		"ip_address_source": param.IpAddressSource,
		"agent":             param.Agent,
		"trx_date":          param.TrxDate,
		"type_user":         param.TypeUser,
	}).Info()
}

func LogDebug(msg string) {
	timestamp := setLogFile(0)
	logText.Debug(fmt.Sprintf("%s [%s] %s", timestamp, "", msg))
}

func LogDebugs(refnum, msg string) {
	timestamp := setLogFile(0)
	logText.Debug(fmt.Sprintf("%s [%s] %s", timestamp, refnum, msg))
}

func LogDebugJSON(data map[string]interface{}) {
	timestamp := setLogFile(0)
	data["timestamp"] = timestamp
	logJSON.WithFields(data).Debug()
}

func LogInfoJSON(data map[string]interface{}) {
	timestamp := setLogFile(0)
	data["timestamp"] = timestamp
	logJSON.WithFields(data).Info()
}

func getCleanMap(m map[string]interface{}) map[string]interface{} {
	delete(m, "password")
	delete(m, "mother_maiden_name")
	delete(m, "cellphone_number")
	delete(m, "cif")
	delete(m, "address")
	delete(m, "born_date")
	return m
}

func Minify(r interface{}) map[string]interface{} {
	js, _ := json.Marshal(r)
	var m map[string]interface{}
	_ = json.Unmarshal(js, &m)

	minifyThreshold := 100

	minifyThresholdRaw := os.Getenv("LOG_MINIFY_TRESHOLD")
	if threshold, err := strconv.Atoi(minifyThresholdRaw); err == nil {
		minifyThreshold = threshold
	}

	for k, v := range m {
		if k == "response_data" || k == "responseData" {
			_, ok := v.(map[string]interface{})
			if !ok {
				m[k] = map[string]interface{}{}
			}
		}

		s := fmt.Sprintf("%v", v)
		if len(s) > minifyThreshold {
			m[k] = "panjang"

			_, ok := v.(string)
			if !ok || k == "response_data" || k == "responseData" {
				m[k] = map[string]interface{}{}
			}
		}
	}

	// marshal m to string json
	jsm, _ := json.Marshal(m)
	strJsm := string(jsm)

	for _, key := range maskedKey {
		regexPattern := `"` + key + `":"(.*?)"`
		re := regexp.MustCompile(regexPattern)
		strJsm = re.ReplaceAllString(strJsm, `"`+key+`":"`+maskedStr+`"`)

		regexPattern = `\\"` + key + `\\":\\"(.*?)\\"`
		re = regexp.MustCompile(regexPattern)
		strJsm = re.ReplaceAllString(strJsm, `\"`+key+`\":\"`+maskedStr+`\"`)
	}

	// unmarshal maskedInput to map
	_ = json.Unmarshal([]byte(strJsm), &m)

	return m
}
