package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type Payload struct {
	Me       string `json:"me,omitempty"`
	One2One  string `json:"one2one,omitempty"`
	Meeting  string `json:"meeting,omitempty"`
	Person   string `json:"person,omitempty"`
	Solo     *bool  `json:"solo,omitempty"`
	Calendar string `json:"calendar,omitempty"`
	Offset   string `json:"offset,omitempty"`
	Range    string `json:"range,omitempty"`
}

func main() {
	port := flag.String("port", "4096", "Port to listen on")
	help := flag.Bool("help", false, "Display help information")

	flag.Parse()

	if *help {
		printUsage()
		return
	}

	http.HandleFunc("/", usageHandler) // New handler registration
	http.HandleFunc("/calendar", corsMiddleware(calendarHandler))
	log.Printf("Service started, listening on port %s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func printUsage() {
	printUsageTo(os.Stdout)
}

func printUsageTo(w io.Writer) {
	fmt.Fprint(w, usage)
}

func usageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	printUsageTo(w)
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://app.tana.inc")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func calendarHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		log.Printf("Processed request: %s %s, returned status %s", r.Method, r.URL, w.Header().Get("X-HTTP-Status-Code"))
	}()

	if r.Method != http.MethodPost {
		w.Header().Set("X-HTTP-Status-Code", fmt.Sprintf("%d", http.StatusMethodNotAllowed))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("X-HTTP-Status-Code", fmt.Sprintf("%d", http.StatusInternalServerError))
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var payload Payload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		w.Header().Set("X-HTTP-Status-Code", fmt.Sprintf("%d", http.StatusBadRequest))
		http.Error(w, fmt.Sprintf("Error unmarshalling JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Run the calendar_auth.scpt1 script
	output, err := runCalendarAuthScript()
	if err != nil {
		w.Header().Set("X-HTTP-Status-Code", fmt.Sprintf("%d", http.StatusInternalServerError))
		http.Error(w, fmt.Sprintf("Failed to run calendar_auth.scpt1 script.\nError %v", output), http.StatusInternalServerError)
		return
	}

	output, err = runCalendarSwiftScript(payload)
	if err != nil {
		w.Header().Set("X-HTTP-Status-Code", fmt.Sprintf("%d", http.StatusInternalServerError))
		http.Error(w, fmt.Sprintf("Error running getcalendar.swift script.\nError %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-HTTP-Status-Code", fmt.Sprintf("%d", http.StatusOK))
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, output)
}

func runCalendarAuthScript() (string, error) {
	cmd := exec.Command("osascript", "./scripts/calendar_auth.scpt")
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func runCalendarSwiftScript(payload Payload) (string, error) {
	cmd := exec.Command("./scripts/getcalendar.swift", "-noheader")
	var args []string

	if payload.Calendar != "" {
		args = append(args, "-calendar", payload.Calendar)
	}
	if payload.Me != "" {
		args = append(args, "-me", payload.Me)
	}
	if payload.One2One != "" {
		args = append(args, "-one2one", payload.One2One)
	}
	if payload.Meeting != "" {
		args = append(args, "-meeting", payload.Meeting)
	}
	if payload.Person != "" {
		args = append(args, "-person", payload.Person)
	}
	if payload.Offset != "" {
		args = append(args, "-offset", payload.Offset)
	}
	if payload.Range != "" {
		args = append(args, "-range", payload.Range)
	}
	if payload.Solo != nil && *payload.Solo {
		args = append(args, "-solo")
	}

	cmd.Args = append(cmd.Args, args...)

	output, err := cmd.CombinedOutput()
	return string(output), err
}

const usage = `Usage: go run main.go [-port <port_number>] [-help]

Options:
  -port <port_number>  Port number to listen on. Default is 4096.
  -help                Display help information.

API Endpoint:
  POST /calendar

JSON Payload:
{
  "me": "<string>",
  "one2one": "<string>",
  "meeting": "<string>",
  "person": "<string>",
  "solo": <boolean>,
  "calendar": "<string>",
  "offset": "<string>",
  "range": "<string>"
}`
