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
   Class string `json:"class"`
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
	Name	string   `json:"name"`
	TL  	uint8	   `json:"tl"`
	MCr	float32	`json:"mcr"`
   RangeClass string   `json:"rangeClass"`
   MountClass    string   `json:"mountClass"`
}

//
//  global Articles array
//  the poor man's database
//
var MountMap map[string]Mount
var RangeMap map[string]Range
var SensorMap map[string]Sensor
var WeaponMap map[string]Sensor  // for now at least
var RangeClass map[string]map[string]Range

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

func createComponent(typ string) Component {
   var component Component
   component.Stage = "Std"
   component.Type = typ
   component.Name = "unknown"
   component.Label = "unknown"
   component.MCr = 0.0
   component.TL = 0
   return component
}

///////////////////////////////////////////////////////////////////////
//
// RANGE MANAGEMENT
//
///////////////////////////////////////////////////////////////////////
func getAllRanges(w http.ResponseWriter, r *http.Request) {
   json.NewEncoder(w).Encode(RangeMap)
}

func createNewRange(w http.ResponseWriter, r *http.Request) {
   // get the body of our POST request
   // unmarshal this into a new Article struct
   // append this to our array.    
   reqBody, _ := ioutil.ReadAll(r.Body)
   var acsRange Range 
   json.Unmarshal(reqBody, &acsRange)
   // update our map
   RangeMap[ acsRange.Type ] = acsRange;

   json.NewEncoder(w).Encode(acsRange)
}

///////////////////////////////////////////////////////////////////////
//
// MOUNT MANAGEMENT
//
///////////////////////////////////////////////////////////////////////
func getAllMounts(w http.ResponseWriter, r *http.Request) {
   json.NewEncoder(w).Encode(MountMap)
}

func createNewMount(w http.ResponseWriter, r *http.Request) {
   // get the body of our POST request
   // unmarshal this into a new Article struct
   // append this to our array.    
   reqBody, _ := ioutil.ReadAll(r.Body)
   var mount Mount 
   json.Unmarshal(reqBody, &mount)
   // update our map
   MountMap[ mount.Type ] = mount;

   json.NewEncoder(w).Encode(mount)
}

///////////////////////////////////////////////////////////////////////
//
// SENSOR MANAGEMENT
//
///////////////////////////////////////////////////////////////////////
func getAllSensors(w http.ResponseWriter, r *http.Request) {
   json.NewEncoder(w).Encode(SensorMap)
}

func createNewSensor(w http.ResponseWriter, r *http.Request) {
   // get the body of our POST request
   // unmarshal this into a new Article struct
   // append this to our array.    
   reqBody, _ := ioutil.ReadAll(r.Body)
   var sensor Sensor 
   json.Unmarshal(reqBody, &sensor)
   // update our map
   SensorMap[ sensor.Type ] = sensor;

   json.NewEncoder(w).Encode(sensor)
}

///////////////////////////////////////////////////////////////////////
//
// WEAPON MANAGEMENT
//
///////////////////////////////////////////////////////////////////////
func getAllWeapons(w http.ResponseWriter, r *http.Request) {
   json.NewEncoder(w).Encode(WeaponMap)
}

func createNewWeapon(w http.ResponseWriter, r *http.Request) {
   // get the body of our POST request
   // unmarshal this into a new Article struct
   // append this to our array.    
   reqBody, _ := ioutil.ReadAll(r.Body)
   var weapon Sensor 
   json.Unmarshal(reqBody, &weapon)
   // update our map
   WeaponMap[ weapon.Type ] = weapon;

   json.NewEncoder(w).Encode(weapon)
}

///////////////////////////////////////////////////////////////////////
//
//  BUILDERS
// 
///////////////////////////////////////////////////////////////////////
func buildMount(mnt string, rng string) Mount {
   //
   //  Figure out the Mount
   //
   var mount_object Mount
	
   mount_object.MCr = 0.0
   mount_object.Tons = 0
 
   mount_object= MountMap[mnt]

   // 
   //  Modify the Mount by Range
   //
   var range_object Range
   var rcheck bool

   range_object, rcheck = RangeMap[rng]
   if rcheck {
	   mount_object.Tons *= uint16(range_object.TonsMod)
	   mount_object.MCr *= range_object.CostMod
   }

   return mount_object
}

func buildSensor(w http.ResponseWriter, r *http.Request) {
   vars := mux.Vars(r)
   typ := vars["type"]
   
   reqBody, _ := ioutil.ReadAll(r.Body)

   component := createComponent(typ)

   //
   //   Figure out the sensor type
   //
   var sensor_object Sensor
   var ncheck bool

   sensor_object, ncheck = SensorMap[typ]
   if ncheck {
      component.MCr = sensor_object.MCr
	   component.TL  = sensor_object.TL
	   component.Name = sensor_object.Name
	}

   //
   //   Is this the right place to do this?
   //
   json.Unmarshal(reqBody, &component)

   //
   //  Figure out the Mount and Range
   //
   if component.Mount == "" {
      component.Mount = "Surf"
   }
   mount_object := buildMount(component.Mount, component.Range)
   range_object := RangeMap[component.Range] // still need the TL Mod

   // 
   //  Now put the component together
   //
   component.Tons = mount_object.Tons 
   component.MCr  += float32(mount_object.MCr)
   component.TL   += uint8(range_object.TLMod)
   component.Label = component.Range + " " + mount_object.Type + " " + component.Name + "-" + strconv.Itoa(int(component.TL))

   json.NewEncoder(w).Encode(component)
}

func buildWeapon(w http.ResponseWriter, r *http.Request) {
   vars := mux.Vars(r)
   typ := vars["type"]
   
   reqBody, _ := ioutil.ReadAll(r.Body)

   component := createComponent(typ)

   //
   //   Figure out the type
   //
   var weapon_object Sensor
   var ncheck bool

   weapon_object, ncheck = WeaponMap[typ]
   if ncheck {
      component.MCr = weapon_object.MCr
	   component.TL  = weapon_object.TL
	   component.Name = weapon_object.Name
	}

   //
   //   Is this the right place to do this?
   //
   json.Unmarshal(reqBody, &component)

   //
   //  Figure out the Mount and Range
   //
   if component.Mount == "" {
      component.Mount = "T1"    // TODO
   }
   mount_object := buildMount(component.Mount, component.Range)
   range_object := RangeMap[component.Range] // still need the TL Mod

   // 
   //  Now put the component together
   //
   component.Tons = mount_object.Tons 
   component.MCr  += float32(mount_object.MCr)
   component.TL   += uint8(range_object.TLMod)
   component.Label = component.Range + " " + mount_object.Type + " " + component.Name + "-" + strconv.Itoa(int(component.TL))

   json.NewEncoder(w).Encode(component)
}

func handleRequests() {

	//http.HandleFunc("/", homePage)
	//http.HandleFunc("/articles", returnAllArticles)

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)

   myRouter.HandleFunc("/ranges",  getAllRanges).Methods("GET")
   myRouter.HandleFunc("/ranges",  createNewRange).Methods("POST")

   myRouter.HandleFunc("/mounts",  getAllMounts).Methods("GET")
   myRouter.HandleFunc("/mounts",  createNewMount).Methods("POST")

   myRouter.HandleFunc("/sensors", getAllSensors).Methods("GET")
   myRouter.HandleFunc("/sensors",  createNewSensor).Methods("POST")
	myRouter.HandleFunc("/sensors/{type}", buildSensor).Methods("POST")

   myRouter.HandleFunc("/weapons", getAllWeapons).Methods("GET")
   myRouter.HandleFunc("/weapons", createNewWeapon).Methods("POST")
   myRouter.HandleFunc("/weapons/{type}", buildWeapon).Methods("POST")

	log.Fatal(http.ListenAndServe(":1317", myRouter))
}

func main() {

	fmt.Println("Rest API 2 - Mux Routers")

	//
	// populate our trivial database
	//
	MountMap = map[string]Mount {
		"Surf": {Type:"Surf", Tons:0,  MCr:0},
		"T1": {Type:"T1",	Tons:1,		MCr:0.2},
		"T2": {Type:"T2",	Tons:1,		MCr:0.6},
		"T3": {Type:"T3",	Tons:1,		MCr:1.0},
		"T4": {Type:"T4",	Tons:1,		MCr:1.5},
		"B1": {Type:"B1",	Tons:3,		MCr:5},
		"B2": {Type:"B2",	Tons:6,		MCr:7},
	}
   
   /*
   RangeClass = map[string] map[string] Range {
      "R": {
		"Vl": {Type:"Vl",   TLMod:-2,   CostMod:0.33,   TonsMod:0.33},
		"D" : {Type:"D",    TLMod:-1,   CostMod:0.5,    TonsMod:0.5},
		"Vd": {Type:"Vd",	  TLMod:0,	  CostMod:1.0,	   TonsMod:1.0},
		"Or": {Type:"Or",   TLMod:1,    CostMod:3.0,    TonsMod:2.0},
		"Fo": {Type:"Fo",   TLMod:2,    CostMod:5.0,    TonsMod:3.0},
		"G":  {Type:"G",    TLMod:3,    CostMod:6.0,    TonsMod:4.0},
      },
      "S": {
		"BR": {Type:"BR",   TLMod:-3,   CostMod:0.25,   TonsMod:0.25},
      "FR": {Type:"FR",   TLMod:-2,   CostMod:0.33,   TonsMod:0.33},
		"SR": {Type:"SR",   TLMod:-1,   CostMod:0.5,    TonsMod:0.5},
		"AR": {Type:"AR",   TLMod:0,    CostMod:1.0,    TonsMod:1.0},
		"LR": {Type:"LR",   TLMod:1,    CostMod:3.0,    TonsMod:2.0},
		"DS": {Type:"DS",   TLMod:2,    CostMod:5.0,    TonsMod:3.0},
      },
   }
   */

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
		"C": {Type:"C", Name:"Communicator", 		RangeClass:"S",  TL:8,  MCr: 1.0, MountClass: "Surf"},
		"H": {Type:"H", Name:"HoloVisor",    		RangeClass:"S",  TL:18, MCr: 1.0, MountClass: "Surf"},
		"T": {Type:"T", Name:"Scope",        		RangeClass:"S",  TL:9,  MCr: 1.0, MountClass: "Surf"},
		"V": {Type:"V", Name:"Visor",        		RangeClass:"S",  TL:14, MCr: 1.0, MountClass: "Surf"},
		"W": {Type:"W", Name:"CommPlus",     		RangeClass:"S",  TL:15, MCr: 1.0, MountClass: "Surf"},

		"E": {Type:"E", Name:"EMS",				 	RangeClass:"S",  TL:12, MCr: 1.0, MountClass: "Surf"},
		"G": {Type:"G", Name:"Grav Sensor",			RangeClass:"S",  TL:13, MCr: 1.0, MountClass: "Surf"},
		"N": {Type:"N", Name:"Neutrino Detector", RangeClass:"S",  TL:10, MCr: 1.0, MountClass: "Surf"},
		"R": {Type:"R", Name:"Radar",				   RangeClass:"S",  TL:9,  MCr: 1.0, MountClass: "Surf"},
		"S": {Type:"S", Name:"Scanner",				RangeClass:"S",  TL:19, MCr: 1.0, MountClass: "Surf"},

		"A": {Type:"A", Name:"Activity Sensor",	RangeClass:"R",  TL:11, MCr: 0.1, MountClass: "Surf"},
	}

   WeaponMap = map[string]Sensor {
      "K": {Type:"K", Name:"Pulse Laser",       RangeClass:"R",  TL:9,  MCr: 0.3, MountClass: "Turret"},
      "L": {Type:"L", Name:"Beam Laser",        RangeClass:"R",  TL:10, MCr: 0.5, MountClass: "Turret"},

      "A": {Type:"A", Name:"Particle Accelerator", RangeClass:"S", TL:11, MCr: 2.5, MountClass: "Barbette"},
      "M": {Type:"M", Name:"Meson Gun",            RangeClass:"S", TL:13, MCr: 5.0, MountClass: "Main"},
   }

	handleRequests()
}
