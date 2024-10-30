package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alienxp03/spectral/api/service"
	"github.com/alienxp03/spectral/client"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	form = `<form method="POST">
                    Start time: <input type="text" name="start_time" value="%s" />
                    End time: <input type="text" name="end_time" value="%s" />
                    <input type="submit" value="Get Usage" />
                  </form>`
)

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("Starting web server at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		w.Header().Set("Content-Type", "text/html")

		startTime := r.FormValue("start_time")
		endTime := r.FormValue("end_time")
		usages, err := client.GetUsages(startTime, endTime)
		if err != nil {
			renderError(w, err)
			return
		}

		renderBody(w, startTime, endTime, usages)
		return
	}
}

func renderError(w http.ResponseWriter, err error) {
	fmt.Fprintf(w, form, "", "")
	fmt.Fprintf(w, fmt.Sprintf("Error: %s", err.Error()), "")
}

func renderBody(w http.ResponseWriter, startTime, endTime string, res *service.GetUsageResponse) {
	f := message.NewPrinter(language.English)
	fmt.Fprintf(w, form, startTime, endTime)
	fmt.Fprintf(w, "Total usage: %s<br/>", f.Sprintf("%f", res.Data.Total))
	fmt.Fprintf(w, "Records: %d<br/>", len(res.Data.Usages))

	fmt.Fprintf(w, "<table>")
	fmt.Fprintf(w, "<tr><th>Time</th><th>Usage</th></tr>")
	for _, usage := range res.Data.Usages {
		fmt.Fprintf(w, "<tr><td>%s</td><td>%f</td></tr>", usage.Time, usage.Usage)
	}
	fmt.Fprintf(w, "</table>")
}
