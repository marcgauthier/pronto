package pronto

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

/*
	Hold the file name currently open to receive new log entries
*/

var log_file string

/*
	Hold the container for the os file object currently open
*/

var log_f *os.File

/*
	path of where the files are saved, last character must be /
*/

var logpath string

/*
	channel to send line to the file
*/
var channelFile chan string

func Start(path string, LogToConsole bool, chanelSize int) error {

	if string(path[len(path)-1:]) != "/" {
		path = path + "/"
	}
	logpath = path

	channelFile = make(chan string, chanelSize)

	go func() {
		for line := range channelFile {

			if LogToConsole {
				fmt.Println(line)
			}

			err := rotate_file_check()

			if err != nil {
				// no log has to use fmt.print.
				log.Println("unable to open new log file")
			} else {
				log_f.Write([]byte(line))
			}
		}
	}()

	return nil

}

/*
	Check if the day as changed since we open the current log file
	if it as then close the current file and create a new file.
	if there is an error creating the file keep trying
*/

func rotate_file_check() error {

	var err error

	/*
		Use today date to create a log file name.
	*/

	f := logpath + time.Now().UTC().Format("2006-01-02.log")

	if f != log_file {

		// close the current file before we open a new one.
		if log_f != nil {
			log_f.Close()
		}

		log_f, err = os.OpenFile(
			f,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)
		if err != nil {
			return err

		}

	}
	return nil
}

/*
	Write a log message in the file and on the console.
*/

func log_msg(s, logtype, filename string, line int) {

	filename = filepath.Base(filename)
	message := time.Now().UTC().Format("[2006-01-02 15:04:05] ["+logtype+"] [") + filename + " " + strconv.Itoa(line) + "] " + s + "\r"

	channelFile <- message

}

/*
	Write a log message in the file
*/
func Info(s string) {
	_, filename, line, _ := runtime.Caller(1)
	log_msg(s, "INFO ", filename, line)
}

/*
	Write a log message in the file
*/
func Warn(s string) {
	_, filename, line, _ := runtime.Caller(1)
	log_msg(s, "WARN ", filename, line)
}

/*
	Write a log message in the file
*/
func Error(s string) {
	_, filename, line, _ := runtime.Caller(1)
	log_msg(s, "ERROR", filename, line)
}
