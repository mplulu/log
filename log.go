package log

import (
	"database/sql"
	"fmt"
	"os"
	"runtime"
	"time"
)

var rootDirectory string
var httpRootUrl string
var shouldLogToFile bool

type DBItf interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var dataCenter DBItf

func RegisterDataCenter(registeredDataCenter DBItf) {
	dataCenter = registeredDataCenter
}

// Z2lhLmRhbmduZ3V5ZW5ob2FuZ0BnbWFpbC5jb20=
// YmdnbGd3dWl4emF3Y3Rvdw==

func EnableLogToFile() {
	shouldLogToFile = true
}

func SetRootDirectory(aRootDirectory string) {
	rootDirectory = aRootDirectory
}

func SetHttpRootUrl(aHttpRootUrl string) {
	httpRootUrl = aHttpRootUrl
}

func Log(format string, a ...interface{}) string {
	logStr := fmt.Sprintf(format, a...)
	logStr = fmt.Sprintf("Info <%s>: %s", time.Now().Format(time.ANSIC), logStr)
	fmt.Println(logStr)

	return logStr
}

func GetStack() string {
	trace := make([]byte, 8192)
	count := runtime.Stack(trace, false)
	content := fmt.Sprintf("Dump (%d bytes):\n %s \n", count, trace[:count])
	return content
}

func LogSerious(format string, a ...interface{}) string {
	logStr := fmt.Sprintf(format, a...)
	logStrSerious := fmt.Sprintf("SERIOUS <%s>: %s", time.Now().Format(time.ANSIC), logStr)
	fmt.Println(logStrSerious)

	// fmt.Println("SENDADMINMAIL")
	trace := make([]byte, 8192)
	count := runtime.Stack(trace, false)
	content := fmt.Sprintf("%s \n Dump (%d bytes):\n %s \n", logStrSerious, count, trace[:count])
	go sendMail(logStr, content)

	return logStrSerious
}

func LogSeriousWithStack(format string, a ...interface{}) string {
	title := fmt.Sprintf(format, a...)
	trace := make([]byte, 8192)
	count := runtime.Stack(trace, false)
	content := fmt.Sprintf("Dump (%d bytes):\n %s \n", count, trace[:count])
	logStr := fmt.Sprintf("SERIOUS <%s>: %s\n%s", time.Now().Format(time.ANSIC), title, content)
	fmt.Println(logStr)

	go sendMail(title, content)

	return logStr
}

func CreateFileAndLog(fileName string, format string, a ...interface{}) (content string, filePath string) {
	if !shouldLogToFile {
		return "", ""
	}
	logStr := fmt.Sprintf(format, a...)
	logStr = fmt.Sprintf("Info <%s>: %s\n", time.Now().Format(time.ANSIC), logStr)

	// generate file path
	filePath = fmt.Sprintf("%s/conf/log/%s", rootDirectory, fileName)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		LogSerious("log to file %s fail %s,log content %s", fileName, err.Error(), logStr)
		return "", ""
	}

	defer f.Close()

	if _, err = f.WriteString(logStr); err != nil {
		LogSerious("log to file %s fail %s,log content %s", fileName, err.Error(), logStr)
	}
	return logStr, filePath
}

func LogFile(fileName string, format string, a ...interface{}) string {
	if !shouldLogToFile {
		return ""
	}
	logStr := fmt.Sprintf(format, a...)
	logStr = fmt.Sprintf("Info <%s>: %s\n", time.Now().Format(time.ANSIC), logStr)

	// generate file path
	filePath := fmt.Sprintf("%s/conf/log/%s", rootDirectory, fileName)
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		LogSerious("log to file %s fail %s,log content %s", fileName, err.Error(), logStr)
	}

	defer f.Close()

	if _, err = f.WriteString(logStr); err != nil {
		LogSerious("log to file %s fail %s,log content %s", fileName, err.Error(), logStr)
	}
	return logStr
}

func LogSeriousFile(fileName string, format string, a ...interface{}) string {
	if !shouldLogToFile {
		return ""
	}
	logStr := fmt.Sprintf(format, a...)
	logStr = fmt.Sprintf("SERIOUS <%s>: %s\n", time.Now().Format(time.ANSIC), logStr)

	// generate file path
	filePath := fmt.Sprintf("%s/conf/log/%s", rootDirectory, fileName)
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		LogSerious("log to file %s fail %s,log content %s", fileName, err.Error(), logStr)
	}

	defer f.Close()

	if _, err = f.WriteString(logStr); err != nil {
		LogSerious("log to file %s fail %s,log content %s", fileName, err.Error(), logStr)
	}
	return logStr
}

func SendMailWithCurrentStack(message string) string {
	trace := make([]byte, 8192)
	count := runtime.Stack(trace, false)
	content := fmt.Sprintf("%s \n Dump (%d bytes):\n %s \n", message, count, trace[:count])
	sendMail(message, content)
	Log(content)
	return content
}

func sendMail(title, content string) {
	logToDb(title, fmt.Sprintf("%s", content))
}

func DumpRows(rows *sql.Rows) {
	cols, err := rows.Columns()
	if err != nil {
		fmt.Println("Failed to get columns", err)
		return
	}

	// Result is your slice string.
	rawResult := make([][]byte, len(cols))
	result := make([]string, len(cols))

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i, _ := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			fmt.Println("Failed to scan row", err)
			return
		}

		for i, raw := range rawResult {
			if raw == nil {
				result[i] = "\\N"
			} else {
				result[i] = string(raw)
			}
		}

		fmt.Printf("%#v\n", result)
	}
}

func logToDb(title string, content string) {
	if dataCenter != nil {
		_, err := dataCenter.Exec("INSERT INTO log_serious (title, content) VALUES ($1, $2)", title, content)
		if err != nil {
			fmt.Println("err log serious", title, content, err)
		}
	}

}
