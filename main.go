package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type Mount struct {
	Type string  `json:"type"`
	Tons uint16   `json:"tons"`
	MCr  float32 `json:"mcr"`
}

type Range struct {
	Type 	string  `json:"type"`
	TLMod	uint8   `json:"tlmod"`
	CostMod	float32 `json:"costmod"`
	TonsMod float32 `json:"tonsmod"`
}

type Component struct {
   Type    string  `json:"type"`
   Stage   string  `json:"stage"`
   Mount   string  `json:"mount"`
   Range   string  `json:"range"`
   TL      uint8   `json:"tl"`
   Tons    uint16  `json:"tons"`
   MCr     float32 `json:"mcr"`
}

type Sensor struct {
	Type 	string 	`json:"type"`
	TL  	uint8	`json:"tl"`
	MCr		float32	`json:"mcr"`
}

//
//  global Articles array
//  the poor man's database
//
var MountMap map[string]Mount
var RangeMap map[string]Range
var SensorMap map[string]Sensor

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Homepage!")
	fmt.Println("Endpoint hit: homePage")
}

/*
func createNewArticle(w http.ResponseWriter, r *http.Request) {
   // get the body of our POST request
   // unmarshal this into a new Article struct
   // append this to our Articles array.    
   reqBody, _ := ioutil.ReadAll(r.Body)
   var article Article 
   json.Unmarshal(reqBody, &article)
   // update our global Articles array to include
   // our new Article
   Articles = append(Articles, article)

   json.NewEncoder(w).Encode(article)
}
*/

func createSensor(w http.ResponseWriter, r *http.Request) {
   vars := mux.Vars(r)
   typ := vars["type"]
   
   var comp_mcr float32 = 0
   var comp_tl  uint8 = 0

   var mnt_mcr float32 = 0.0
   var mnt_vol uint16 = 0

   reqBody, _ := ioutil.ReadAll(r.Body)

   var component Component
       component.Stage = "std"

   json.Unmarshal(reqBody, &component)

   component.Type = typ

   sensor_object, ncheck := SensorMap[typ]
   if ncheck {
      comp_mcr = sensor_object.MCr
	  comp_tl  = sensor_object.TL
   }

   mount_object, mocheck := MountMap[component.Mount]
   if mocheck {
      mnt_vol = mount_object.Tons
      mnt_mcr = mount_object.MCr
   }

   range_object, rcheck := RangeMap[component.Range]
   if rcheck {
      mnt_vol = uint16(float32(mnt_vol) * range_object.TonsMod)
      mnt_mcr *= range_object.CostMod
   }

   component.Tons = mnt_vol 
   component.MCr  = float32(mnt_mcr) + comp_mcr
   component.TL   = comp_tl

   json.NewEncoder(w).Encode(component)
}

func handleRequests() {

	//http.HandleFunc("/", homePage)
	//http.HandleFunc("/articles", returnAllArticles)

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/sensor/{type}", createSensor).Methods("POST")

	log.Fatal(http.ListenAndServe(":1317", myRouter))
}

func main() {

	fmt.Println("Rest API 2 - Mux Routers")

	//
	// populate our trivial database
	//
	MountMap = map[string]Mount {

		"T1": {Type:"T1",	Tons:1,		MCr:0.2},
		"T2": {Type:"T2",	Tons:1,		MCr:0.6},
		"T3": {Type:"T3",	Tons:1,		MCr:1.0},
		"T4": {Type:"T4",	Tons:1,		MCr:1.5},
		"B1": {Type:"B1",	Tons:3,		MCr:5},
		"B2": {Type:"B2",	Tons:6,		MCr:7},
	}

   	RangeMap = map[string]Range {
		"Vd": {Type:"Vd",	TLMod:0,	CostMod:1.0,	TonsMod:1.0},
		"Or": {Type:"Or",   TLMod:1,    CostMod:2.0,    TonsMod:2.0},
		"Fo": {Type:"Fo",   TLMod:2,    CostMod:3.0,    TonsMod:3.0},
 	}

	SensorMap = map[string]Sensor {
		"N": {Type:"N",	TL:10, MCr: 1.0},
	}

	handleRequests()
}
