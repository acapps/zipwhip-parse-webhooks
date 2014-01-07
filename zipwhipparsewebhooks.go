package zipwhipparsewebhooks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

/**
 * All Zipwhip Web Hooks return the same data for the message object.
 */
type Message struct {
	Body             string `json:"body"`             // Body of Message.
	BodySize         int    `json:"bodySize"`         // BodySize, number of characters
	Visible          bool   `json:"visible"`          // Is the message visible on portal: True / False
	HasAttachment    bool   `json:"hasAttachment"`    // Does the message have an attachment, can be retrieved via message/get
	FinalDestination string `json:"finalDestination"` // Number, contact, for whom the message is for.
	MessageType      string `json:"messageType"`      // MO, Mobile Originated; ZO, Zipwhip Originated; MT, Mobile Terminated
	Deleted          bool   `json:"deleted"`          // Has the message been deleted; True / False
	Id               int64  `json:"id"`               // ID of the message.
	StatusCode       int    `json:"statusCode"`       // 0 or 4 means it was successfully sent; 1 is prepping; all others probably fail scenario.
	MessageTransport int    `json:"messageTransport"` // Used for internal Zipwhip routing purposes.
	DateCreated      string `json:"dateCreated"`      // Date the message was created in the system.
	Read             bool   `json:"read"`             // Has the message been marked as read: True / False
	FinalSource      string `json:"finalSource"`      // Number, contact, for whom the messsage was sent from.
	DeviceId         int    `json:"deviceId"`         // Id for the owner of account.
}

func (m *Message) ParseJson(inputJson *[]byte) error {
	return json.Unmarshal(*inputJson, m)
}

/**
 * Sidenote, formatting looks nice Mac OS X Terminal, monospaced font. Results with your IDE's console may very.
 */
func (m Message) String() string {
	return fmt.Sprintf("\nBody:\t\t\t\t%s\nBodySize:\t\t\t%d\nVisible:\t\t\t%t\nHasAttachment:\t\t\t%t\nFinalDestination:\t\t%s\nMessageType:\t\t\t%s\nDeleted:\t\t\t%t\nId:\t\t\t\t%d\nStatusCode:\t\t\t%d\nMessageTransport:\t\t%d\nDateCreated:\t\t\t%s\nRead:\t\t\t\t%t\nFinalSource:\t\t\t%s\nDeviceId:\t\t\t%d\n", m.Body, m.BodySize, m.Visible, m.HasAttachment, m.FinalDestination, m.MessageType, m.Deleted, m.Id, m.StatusCode, m.MessageTransport, m.DateCreated, m.Read, m.FinalSource, m.DeviceId)
}

/**
 * Primary workhorse, take the request's body and turn into a Message object
 */
func bodyToMessage(body *[]byte) {
	var message Message
	err := message.ParseJson(body)
	if err != nil {
		log.Printf("An error occurred while parsing JSON:\n\t%s", message)
	}

	log.Println(message)
}

/**
 * Process the specific URI's log them, read the body, and then turn it over to the bodyToMessage function.
 */
func messageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("New event received: %s\n", r.RequestURI)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("An error occurred while reading in the Request's Body:\n\t%s", err)
	}
	go bodyToMessage(&body)
}

func main() {
	http.HandleFunc("/message/send", messageHandler)
	http.HandleFunc("/message/progress", messageHandler)
	http.HandleFunc("/message/receive", messageHandler)
	http.HandleFunc("/message/read", messageHandler)
	http.HandleFunc("/message/delete", messageHandler)

	log.Println("Server listening on port 8090.")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
