package handler

import (
	"fmt"
	"github.com/Rhymen/go-whatsapp"
	"github.com/Rhymen/go-whatsapp/binary/proto"
	h "github.com/aldinokemal/go-whatsapp-web/helpers"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// handle History

// historyHandler for acquiring chat history
type historyHandler struct {
	c        *whatsapp.Conn
	messages []string
}

func (h *historyHandler) ShouldCallSynchronously() bool {
	return true
}

// handles and accumulates history's text messages.
// To handle images/documents/videos add corresponding handle functions
func (h *historyHandler) HandleTextMessage(message whatsapp.TextMessage) {
	authorID := "-"
	screenName := "-"
	if message.Info.FromMe {
		authorID = h.c.Info.Wid
		screenName = ""
	} else {
		if message.Info.Source.Participant != nil {
			authorID = *message.Info.Source.Participant
		} else {
			authorID = message.Info.RemoteJid
		}
		if message.Info.Source.PushName != nil {
			screenName = *message.Info.Source.PushName
		}
	}

	date := time.Unix(int64(message.Info.Timestamp), 0)
	h.messages = append(h.messages, fmt.Sprintf("%s	%s (%s): %s", date,
		authorID, screenName, message.Text))

}

func (h *historyHandler) HandleError(err error) {
	log.Printf("Error occured while retrieving chat history: %s", err)
}

func GetHistory(jid string, wac *whatsapp.Conn) []string {
	// create out history handler
	handler := &historyHandler{c: wac}

	// load chat history and pass messages to the history handler to accumulate
	wac.LoadFullChatHistory(jid, 300, time.Millisecond*300, handler)
	return handler.messages
}

func GetAnyHistory(wac *whatsapp.Conn, chats map[string]struct{}) {
	// show list of chats
	var chatSlice []string
	for chat := range chats {
		chatSlice = append(chatSlice, chat)
		fmt.Printf("%d	%s\n", len(chatSlice), chat)
	}

	fmt.Println("Select chat number to get history for:")
	var index = 0
	for index < 1 || index > len(chatSlice) {
		fmt.Scanf("%d", &index)
	}

	// get history for the selected chat
	fmt.Println("Gathering chat history...")
	messages := GetHistory(chatSlice[index-1], wac)
	for _, message := range messages {
		fmt.Println(message)
	}
}

type waHandler struct {
	c     *whatsapp.Conn
	chats map[string]struct{}
}

func (h *waHandler) ShouldCallSynchronously() bool {
	return true
}

func (h *waHandler) HandleRawMessage(message *proto.WebMessageInfo) {
	// gather chats jid info from initial messages
	if message != nil && message.Key.RemoteJid != nil {
		h.chats[*message.Key.RemoteJid] = struct{}{}
	}
}

func (h *waHandler) HandleError(err error) {

	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v", e.Err)
		log.Println("Waiting 30sec...")
		<-time.After(30 * time.Second)
		log.Println("Reconnecting...")
		err := h.c.Restore()
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
	} else {
		log.Printf("error occoured: %v\n", err)
	}
}

func ReadHistory(c *gin.Context) {
	//create new WhatsApp connection
	wac, err := whatsapp.NewConn(5 * time.Second)
	if err != nil {
		log.Fatalf("error creating connection: %v\n", err)
	}

	//Add handler
	handler := &waHandler{wac, make(map[string]struct{})}
	wac.AddHandler(handler)

	//login or restore
	err = Login(wac)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
		fmt.Println(os.TempDir() + "whatsappSession.gob")
		err = os.Remove(os.TempDir() + "whatsappSession.gob")
		if err != nil {
			fmt.Println(err.Error())
		}
		ReadHistory(c)
		return
	}

	// wait while chat jids are acquired through incoming initial messages
	fmt.Println("Waiting for chats info...")
	<-time.After(5 * time.Second)

	// get history synchronously
	GetAnyHistory(wac, handler.chats)
	fmt.Println("Done. Press Ctrl+C for exit.")

	chanel := make(chan os.Signal, 1)
	signal.Notify(chanel, os.Interrupt, syscall.SIGTERM)
	<-chanel

	//Disconnect safely
	fmt.Println("Shutting down now.")
	session, err := wac.Disconnect()
	if err != nil {
		log.Fatalf("error disconnecting: %v\n", err)
	}
	if err := h.WriteSession(session, "120381"); err != nil {
		log.Fatalf("error saving session: %v", err)
	}
}
