
package main

import (
    "fmt"
    "strings"
    "github.com/heltonmarx/goami/ami"
    "encoding/json"
    "net/http"
	"html/template"
)

func connect() (*ami.Socket, string, bool) {
    uuid, _ := ami.GetUUID()
    socket, err := ami.NewSocket("127.0.0.1:5038")
    if err != nil {
        fmt.Printf("socket error: %v\n", err)
    }
    if _, err := ami.Connect(socket); err != nil {
        fmt.Printf("Connect error (%v)\n", err)
    }

    ret, err := ami.Login(socket, "admin", "ipsafe", "Off", uuid)
    if err != nil || ret == false {
        fmt.Printf("login error (%v)\n", err)
    }
    return socket, uuid, ret
}

type Peer struct {
    Username        string `json:"username"` // Default-Username
    Description     string `json:"description"` // Description
    ObjectType      string `json:"objecttype"` // ChanObjectType
    IP              string `json:"ip"` // Address-IP
    Callerid        string `json:"callerid"` // Callerid
    Callgroup       string `json:"callgroup"` // Named Callgroup
    Context         string `json:"context"` // Context
    Pickupgroup     string `json:"pickupgroup"` // Pickupgroup
}

type Peers struct {
    Item []Peer `json:"items"`
}

func get_peers(socket *ami.Socket, uuid string) ([]byte) {
    peers := Peers{}
    list, _ := ami.SIPpeers(socket, uuid)
    if len(list) > 0 {
        for _, m := range list {
            peer := Peer{}
            message, _ := ami.SIPshowpeer(socket, uuid, m["ObjectName"])
            for k, v := range message {
                //fmt.Printf("%s : %q\n", k, v)
                if (strings.Contains(k, "Default-Username")) { peer.Username = v }
                if (strings.Contains(k, "Description")) { peer.Description = v }
                if (strings.Contains(k, "ChanObjectType")) { peer.ObjectType = v }
                if (strings.Contains(k, "Address-IP")) { peer.IP = v }
                if (strings.Contains(k, "Callerid")) { peer.Callerid = v }
                if (strings.Contains(k, "Callgroup")) { peer.Callgroup = v }
                if (strings.Contains(k, "Context")) { peer.Context = v }
                if (strings.Contains(k, "Pickupgroup")) { peer.Pickupgroup = v }
            }
            peers.Item = append(peers.Item, peer)
        }
    }
    b, _ := json.Marshal(peers)
    return b
}



type Queue struct {
    Name            string `json:"name"`
    Number          string `json:"number"`
    ActiveCalls     string `json:"active_calls"`
}

type Queues struct {
    Item []Queue `json:"items"`
}

func get_queues(socket *ami.Socket, uuid string) ([]byte) {
    queues := Queues{}
    list, _ := ami.Queues(socket, uuid)
    for _, v := range list {
        queue := Queue{}
        queue.Number = v["Queue"]
        queue.Name = v["Name"]
        queue.ActiveCalls = v["Callers"]
        queues.Item = append(queues.Item, queue)
    }
    b, _ := json.Marshal(queues)
    return b
}

func home_handler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("templates/index.html")
    t.Execute(w, r)
}

func main() {
    http.HandleFunc("/list/peers", func(w http.ResponseWriter, r *http.Request) {
        socket, uuid, _ := connect()
        w.Write(get_peers(socket, uuid))
    })

    http.HandleFunc("/list/queues", func(w http.ResponseWriter, r *http.Request) {
        socket, uuid, _ := connect()
        w.Write(get_queues(socket, uuid))
    })

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        home_handler(w, r)
    })

    fmt.Println("Listening...")
    http.ListenAndServe(":9002", nil)

}
