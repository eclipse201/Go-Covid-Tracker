
package main
 
import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "strconv"
    "time"
)

//holds covid information
type GeoData struct{
    ZipCode string `json:"zipCd"`
    Counties []Counties `json:"counties"`
    Geo Geo `json:"geo"`
}

//holds covid information for the counties pulled
type Counties struct{
    CountyName string `json:"countyName"`
    DeathCount int `json:"deathCt"`
    PositiveCount int `json:"positiveCt"`
}

//holds latitude and longitude that the covid information was collected from
type Geo struct {
    LeftBottomLatLong float64 `json:leftBottomLatLong"`
    RightBottomLatLong float64 `json:"rightBottomLatLong"`
}

//collects covid information every minute
func gatherDataOnTimer(){ 
    for {
        time.Sleep(time.Minute)
        go getWeather()
    }
}

//formats the latitude and longitude information
func FormatGeo(info *GeoData) string{
    var formatString = fmt.Sprintf("%f", info.Geo.RightBottomLatLong) + " " + fmt.Sprintf("%f", info.Geo.LeftBottomLatLong)
    return formatString
}

//formats the covid information for each county in the area
func FormatCounties(info *GeoData) string{
    var formatString = ""
    for i := 0; i < len(info.Counties); i++{
       formatString = formatString + "County: " + info.Counties[i].CountyName + "\nDeath Count: " + strconv.Itoa(info.Counties[i].DeathCount)+ "\nPositive Cases: " + strconv.Itoa(info.Counties[i].PositiveCount) + "\n"
    }
    return formatString
}


//retrieves covid information from https://localcoviddata.com/ and returns a formated string of that data
func getWeather()(string) {    
    response, err := http.Get("https://localcoviddata.com/covid19/v1/locations?zipCode=28105")
    if err != nil{
        log.Fatal(err)
    }

    bytes, _ := ioutil.ReadAll(response.Body)
    var geodata GeoData
    json.Unmarshal(bytes, &geodata)

    var location = FormatGeo(&geodata)
    
    var countyInfo = FormatCounties(&geodata)

    var data = "ZipCode: " + geodata.ZipCode + "\nLocation: " + location + "\n" + countyInfo

    defer response.Body.Close()
    return data
}

//starts the server on local host port 8080, runs the code, and presents the covid information on the server
func main() {
    go gatherDataOnTimer()
    var weather = getWeather()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                fmt.Fprint(w,"Covid information\n")
                fmt.Fprint(w,weather)
    })
    
    port := ":8080"
    fmt.Println("Server is running on port" + port)

    log.Fatal(http.ListenAndServe(port, nil))
}