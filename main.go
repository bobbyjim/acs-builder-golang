package main

import (
	"encoding/json"
	"fmt"
	"strconv"
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
	TLMod	int8    `json:"tlmod"`
	CostMod	float32 `json:"costmod"`
	TonsMod float32 `json:"tonsmod"`
}

type Component struct {
   Type    string  `json:"type"`
   Name    string  `json:"name"`
   Label   string  `json:"label"`
   Stage   string  `json:"stage"`
   Mount   string  `json:"mount"`
   Range   string  `json:"range"`
   TL      uint8   `json:"tl"`
   Tons    uint16  `json:"tons"`
   MCr     float32 `json:"mcr"`
}

type Sensor struct {
	Type 	string 	`json:"type"`
	Name	string  `json:"name"`
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
 
   var rng_tl  int8   = 0

   reqBody, _ := ioutil.ReadAll(r.Body)

   var component Component
       component.Stage = "std"

   json.Unmarshal(reqBody, &component)

   component.Type = typ
   component.Name = "unknown"
   component.Label = ""

   sensor_object, ncheck := SensorMap[typ]
   if ncheck {
      comp_mcr = sensor_object.MCr
	  comp_tl  = sensor_object.TL
	  component.Name = sensor_object.Name
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
	  rng_tl  = range_object.TLMod
   }

   component.Tons = mnt_vol 
   component.MCr  = float32(mnt_mcr) + comp_mcr
   component.TL   = uint8(int8(comp_tl) + rng_tl)
   component.Label = range_object.Type + " " + mount_object.Type + " " + sensor_object.Name + "-" + strconv.Itoa(int(component.TL))

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
		"BR": {Type:"BR",   TLMod:-3,   CostMod:0.25,   TonsMod:0.25},
        "FR": {Type:"FR",   TLMod:-2,   CostMod:0.33,   TonsMod:0.33},
		"SR": {Type:"SR",   TLMod:-1,   CostMod:0.5,    TonsMod:0.5},
		"AR": {Type:"AR",   TLMod:0,    CostMod:1.0,    TonsMod:1.0},
		"LR": {Type:"LR",   TLMod:1,    CostMod:3.0,    TonsMod:2.0},
		"DS": {Type:"DS",   TLMod:2,    CostMod:5.0,    TonsMod:3.0},

		"Vl": {Type:"Vl",   TLMod:-2,   CostMod:0.33,    TonsMod:0.33},
		"D" : {Type:"D",    TLMod:-1,   CostMod:0.5,    TonsMod:0.5},
		"Vd": {Type:"Vd",	TLMod:0,	CostMod:1.0,	TonsMod:1.0},
		"Or": {Type:"Or",   TLMod:1,    CostMod:3.0,    TonsMod:2.0},
		"Fo": {Type:"Fo",   TLMod:2,    CostMod:5.0,    TonsMod:3.0},
		"G":  {Type:"G",    TLMod:3,    CostMod:6.0,    TonsMod:4.0},
 	}

	SensorMap = map[string]Sensor {
		"C": {Type:"C", Name:"Communicator", 		TL:8,  MCr: 1.0},
		"H": {Type:"H", Name:"HoloVisor",    		TL:18, MCr: 1.0},
		"T": {Type:"T", Name:"Scope",        		TL:9,  MCr: 1.0},
		"V": {Type:"V", Name:"Visor",        		TL:14, MCr: 1.0},
		"W": {Type:"W", Name:"CommPlus",     		TL:15, MCr: 1.0},

		"E": {Type:"E", Name:"EMS",				 	TL:12, MCr: 1.0},
		"G": {Type:"G", Name:"Grav Sensor",			TL:13, MCr: 1.0},
		"N": {Type:"N",	Name:"Neutrino Detector", 	TL:10, MCr: 1.0},
		"R": {Type:"R", Name:"Radar",				TL:9,  MCr: 1.0},
		"S": {Type:"S", Name:"Scanner",				TL:19, MCr: 1.0},

		"A": {Type:"A", Name:"Activity Sensor",		TL:11, MCr: 0.1},
	}

	handleRequests()
}
