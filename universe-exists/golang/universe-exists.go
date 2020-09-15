// fetch https://lobby.ogame.gameforge.com/api/servers
// break data according to server
// output to separate files in format Name.Language.Number.json - TODO

package main

import "fmt"
import "net/http"
import "io/ioutil"
import "encoding/json"

// Define data structure 
// structure members have to start with capital letters to be exported
type Parameters struct {
	Name		string
	Language	string
}
type Server struct {
  Language    		string	`json:"language"`
  Number      		int		`json:"number"`
  Name       		string	`json:"name"`
  PlayerCount       int		`json:"playerCount"`
  PlayersOnline     int		`json:"playersOnline"`
  Opened     		string	`json:"opened"`
  StartDate       	string	`json:"startDate"`
  EndDate			string	`json:"endDate"`
  ServerClosed		int		`json:"serverClosed"`
  Prefered			int		`json:"prefered"`
  SignupClosed		int		`json:"signupClosed"`
  Settings struct {
	Aks							int		`json:"aks"`
	FleetSpeed					int		`json:"fleetSpeed"`
	WreckField					int		`json:"wreckField"`
	ServerLabel					string	`json:"serverLabel"`
	EconomySpeed				int		`json:"economySpeed"`
	PlanetFields				int		`json:"planetFields"`
	UniverseSize				int		`json:"universeSize"`
	ServerCategory				string	`json:"serverCategory"`
	EspionageProbeRaids			int		`json:"espionageProbeRaids"`
	PremiumValidationGift		int		`json:"premiumValidationGift"`
	DebrisFieldFactorShips		int		`json:"debrisFieldFactorShips"`
	ResearchDurationDivisor		int		`json:"researchDurationDivisor"`
	DebrisFieldFactorDefence 	int		`json:"debrisFieldFactorDefence"`
  }`json:"settings"`
 
}

// report structure
type Report struct {
	StatusCode		int		`json:"statusCode"`
	StatusMessage	string	`json:"statusMessage"`
}


func get_content(serverLanguage string, serverName string) (Server, Report){
	// json data
	url := "https://lobby.ogame.gameforge.com/api/servers"

	res, err := http.Get(url)
	if err != nil {panic(err.Error())}

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {panic(err.Error())}

	var report Report
    var server []Server
	var resulting_server Server
    // unmarshal
    json.Unmarshal(body, &server)
	if err != nil {panic(err.Error())}
	
    // print values of the object
	
	for i := range server {
		if (server[i].Language == serverLanguage && server[i].Name == serverName) {
			resulting_server = server[i]
			report.StatusCode = 200
			report.StatusMessage = "Universe found"
		}
	}
	if (resulting_server == Server{}) {
			report.StatusCode = 404
			report.StatusMessage = "Universe not found"
	}
	return resulting_server, report
}

func main() {

	http.HandleFunc("/", UniverseExists)
	http.ListenAndServe(":4000", nil)
}

func UniverseExists (w http.ResponseWriter, r *http.Request) {
	var params Parameters
	switch r.Method {
	case "POST":	
		json.NewDecoder(r.Body).Decode(&params)
		result, report := get_content(params.Language, params.Name)
		if (report.StatusCode == 404) {
			jsonReport, err := json.Marshal(report)
			if err != nil {panic(err.Error())}
			fmt.Fprintf(w, "%s", string(jsonReport))
		}
		
		if (report.StatusCode == 200) {
			jsonResult, err := json.Marshal(result)
			if err != nil {panic(err.Error())}
			fmt.Fprintf(w, "%s", jsonResult)
		}
	default:
		w.WriteHeader(http.StatusForbidden)
	}	
}
