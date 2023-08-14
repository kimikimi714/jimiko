package log

import (
	"encoding/json"
	"fmt"
	"log"
)

func Info(format string, v ...any) {
	log.Println(entry{Severity: "INFO", Message: fmt.Sprintf(format, v...)})
}

func Warn(format string, v ...any) {
	log.Println(entry{Severity: "WARN", Message: fmt.Sprintf(format, v...)})
}

func Error(format string, v ...any) {
	log.Println(entry{Severity: "ERROR", Message: fmt.Sprintf(format, v...)})
}

// entry defines a log entry.
type entry struct {
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

// String renders an entry structure to the JSON format expected by Cloud Logging.
func (e entry) String() string {
	out, err := json.Marshal(e)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}
	return string(out)
}
