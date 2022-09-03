package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var wClient *whatsmeow.Client
var dClient *discordgo.Session

func Organize(m whatsmeow.DownloadableMessage, mType string) []*discordgo.File {
    b, _ := wClient.Download(m)
    r := bytes.NewBuffer(b)
    
    if mType == "" {
        mType = "text/plain"
    }

    ext := strings.Split(mType, "/")[1]
    
    if mType == "text/plain" {
        ext = "txt"
    }

    name := "attach." + ext

    if len(b) > 7900000 {
        return []*discordgo.File{}
    }
    
    return []*discordgo.File{
        {
            Name:        name,
            ContentType: mType,
            Reader:      r,
        },
    }
}


func eventHandler(evt interface{}) {
    switch v := evt.(type) {
    case *events.Message:
        msg := v.Message
        if !v.Info.IsGroup {return}
        group, err := wClient.GetGroupInfo(v.Info.Chat)
        PanicIfErr(err)
        if group.Name != CHAT_NAME {return}

        files := []*discordgo.File{}
        content := ""

        if m := msg.Conversation; m != nil {
            content = *m
        } else if m := msg.ImageMessage; m != nil {
            files = Organize(m, m.GetMimetype())
            content = m.GetCaption()
        } else if m := msg.DocumentMessage; m != nil {
            files = Organize(m, m.GetMimetype())
            content = m.GetCaption()
        } else if m := msg.AudioMessage; m != nil {
            files = Organize(m, m.GetMimetype())
        } else if m := msg.VideoMessage; m != nil {
            files = Organize(m, m.GetMimetype())
            content = m.GetCaption()
        } else if m := msg.ContactMessage; m != nil {
            b, _ := json.Marshal(m)
            fmt.Println("contact", string(b))
        } else if m := msg.LocationMessage; m != nil {
            content =  fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%v,%v", m.DegreesLatitude, m.DegreesLongitude)
        } else if m := msg.LiveLocationMessage; m != nil {
            if m.Caption != nil {
                content = *m.Caption + " - "
            }
            content += fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%v,%v", m.DegreesLatitude, m.DegreesLongitude)
        } else if m := msg.Chat; m != nil {
            b, _ := json.Marshal(m)
            fmt.Println("chat", string(b))
        }

        if content == "" && files == nil {
            content = "Message type not supported!"
        }

        content = "||<@&" + GUILD_ROLE + ">|| " + content

        _, err = dClient.ChannelMessageSendComplex(GUILD_CHANNEL, &discordgo.MessageSend{
        	Content:         content,
        	Files:           files,

        })

        if err != nil {
            fmt.Println(err)
        }
    }
}

func init() {
    ConfigInit()
}

func init() {
    dbLog := waLog.Stdout("Database", "WARN", true)
    // Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
    container, err := sqlstore.New("sqlite3", "file:deviceStore.db?_foreign_keys=on", dbLog)
    if err != nil {
        panic(err)
    }
    // If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
    deviceStore, err := container.GetFirstDevice()
    if err != nil {
        panic(err)
    }
    clientLog := waLog.Stdout("Client", "WARN", true)
    wClient = whatsmeow.NewClient(deviceStore, clientLog)
    wClient.AddEventHandler(eventHandler)

    if wClient.Store.ID == nil {
        // No ID stored, new login
        qrChan, _ := wClient.GetQRChannel(context.Background())
        err = wClient.Connect()
        if err != nil {
            panic(err)
        }
        for evt := range qrChan {
            if evt.Event == "code" {
                qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)                
            } else {
                fmt.Println("Login event:", evt.Event)
            }
        }
    } else {
        // Already logged in, just connect
        err = wClient.Connect()
        if err != nil {
            panic(err)
        }
    }
}

func init() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
    dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages
    dClient = dg

    err = dClient.Open()
    PanicIfErr(err)
}

func main() {
    // Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    <-c
    
    dClient.Close()
    wClient.Disconnect()
}
