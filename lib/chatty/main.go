package chatty

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"time"
)

var App = mux.NewRouter()

func init() {
	App.HandleFunc("/", home)
	App.HandleFunc("/send", send)
	App.HandleFunc("/await", await)
}

var R = New()

var baseTpl = template.Must(template.ParseFiles("templates/base.tpl"))

var homeTpl = template.Must(baseTpl.ParseFiles("templates/home.tpl"))

func home(w http.ResponseWriter, r *http.Request) {
	homeTpl.Execute(w, nil)
}

func send(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Bad Request Method", 500)
		return
	}
	r.ParseForm()
	R.Put(Message{Name: r.Form.Get("name"), Message: r.Form.Get("message"), Time: time.Now()})
}

func await(w http.ResponseWriter, r *http.Request) {
	m := R.Await()
	s, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	w.Header().Add("content-type", "application/json")
	fmt.Fprintln(w, string(s))
}

type Message struct {
	Name    string
	Time    time.Time
	Message string
}

func (m Message) String() string {
	return fmt.Sprintf("%s: %s: %s", m.Time, m.Name, m.Message)
}

type Messages []Message
type Await chan<- Messages
type Put chan<- Message

type Reactor struct {
	await, recv chan<- Await
	put         chan<- Message
}

func awaitFor(c chan<- Await) Messages {
	x := make(chan Messages, 1)
	c <- x
	return <-x
}

func (r Reactor) Await() Messages {
	return awaitFor(r.await)
}

func (r Reactor) Recv() Messages {
	return awaitFor(r.recv)
}

func (r Reactor) Put(m Message) {
	r.put <- m
}

func New() Reactor {
	a, r, p := make(chan Await), make(chan Await), make(chan Message)
	go func() {
		aqs, ms := []Await{}, Messages{}
		for {
			select {
			case aq := <-a:
				aqs = append(aqs, aq)
			case rx := <-r:
				rx <- append(Messages{}, ms...)
			case m := <-p:
				ms = append(ms, m)
				for i := range aqs {
					aqs[i] <- append(Messages{}, ms...)
				}
				if len(aqs) > 0 {
					aqs = []Await{}
					ms = Messages{}
				}
			}
		}
	}()
	return Reactor{a, p}
}
