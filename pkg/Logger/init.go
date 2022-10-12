package Logger

import (
	"os"
	"strings"
	"time"
)

var (
	strLogDir      string = "d:/Logs/Reveroxy/"
	strFormLogFile string = "20060102"
	strFormTime    string = "[15:04:05]"
)

func WriteLog(_strContent ...string) {
	strLogFile := strLogDir + time.Now().Format(strFormLogFile) + ".log"
	fpLog, err := os.OpenFile(strLogFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		os.FileMode(0644))

	if err != nil {
		panic(err)
	}

	defer fpLog.Close()

	strContent := strings.Join(_strContent, " ")
	if !strings.HasSuffix(strContent, "\n") {
		strContent = strContent + "\n"
	}
	fpLog.WriteString(time.Now().Format(strFormTime) + " " + strContent)
}

func WriteError(_strContent ...string) {
	strLogFile := strLogDir + "error/" + time.Now().Format(strFormLogFile) + ".log"
	fpLog, err := os.OpenFile(strLogFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		os.FileMode(0644))

	if err != nil {
		panic(err)
	}

	defer fpLog.Close()

	strContent := strings.Join(_strContent, " ")
	if !strings.HasSuffix(strContent, "\n") {
		strContent = strContent + "\n"
	}
	fpLog.WriteString(time.Now().Format(strFormTime) + " " + strContent)
}
