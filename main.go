package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DB struct {
	Session *mgo.Session
}

type Hero struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type HeroResponse struct {
	Name string `json:"name"`
}

type Heroes []Hero

type Todo struct {
	Id          bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	TodoMessage string        `json:"todoMessage,omitempty" bson:"todoMessage"`
	CreatedAt   time.Time     `json:"createdAt,omitempty" bson:"createdAt"`
}

type Todos []Todo

const (
	col  string = "todos"
	dev  string = "./"
	node string = "node_modules"
)

func main() {
	// Creates a new serve mux
	mux := http.NewServeMux()

	// Create room for static files serving
	mux.Handle("/node_modules/", http.StripPrefix("/node_modules", http.FileServer(http.Dir("./node_modules"))))
	mux.Handle("/html/", http.StripPrefix("/html", http.FileServer(http.Dir("./html"))))
	mux.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./js"))))
	mux.Handle("/ts/", http.StripPrefix("/ts", http.FileServer(http.Dir("./ts"))))
	mux.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./css"))))

	// Do your api stuff**
	mux.Handle("/api/register", util.Adapt(api.RegisterHandler(mux),
		api.GetMongoConnection(),
		api.CheckEmptyUserForm(),
		api.EncodeUserJson(),
		api.ExpectBody(),
		api.ExpectPOST(),
	))
	mux.HandleFunc("/api/login", api.Login)
	mux.HandleFunc("/api/authenticate", api.Authenticate)

	// Any other request, we should render our SPA's only html file,
	// Allowing angular to do the routing on anything else other then the api
	// and the files it needs for itself to work.
	// Order here is critical. This html should contain the base tag like
	// <base href="/"> *href here should match the HandleFunc path below
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	mux.HandleFunc("/myhandle/heroes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			js, err := json.Marshal(getHeroesFile())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		case "POST":
			decoder := json.NewDecoder(r.Body)
			var t HeroResponse
			err := decoder.Decode(&t)
			checkError("Problem decoding POST request body into JSON", err)
			defer r.Body.Close()
			//generate a unique id, define a hero, append it to the whole heroes struct and commit it to disk
			log.Println(t)
			w.WriteHeader(http.StatusOK)
		case "PUT":
			// Update an existing record.
		case "DELETE":
			// Remove the record.
		default:
			http.Error(w, "Error", http.StatusInternalServerError)
		}
	})

	//Static Files
	//Start Server
	http.ListenAndServe(":1870", nil)
}

func getHeroesFile() Heroes {
	var h Heroes
	heroesFile := "heroes.json"
	if _, err := os.Stat(heroesFile); err == nil {
		f, _ := os.Open(heroesFile)
		decoder := json.NewDecoder(f)
		err := decoder.Decode(&h)
		checkError("Problem decoding JSON file", err)
	}
	return h
}

func setHeroesFile(h Heroes) {
	heroesFile := "heroes.json"
	if _, err := os.Stat(heroesFile); err == nil {
		os.Remove(heroesFile)
	}
	file, err := os.Create(heroesFile)
	checkError("Cannot create file", err)
	defer file.Close()
	writer := bufio.NewWriter(file)
	checkError("error opening the writer", err)
	js, err := json.Marshal(h)
	_, err = writer.Write(js)
	checkError("Cannot write to file", err)
	defer writer.Flush()
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
