package service

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/mdp/qrterminal/v3"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type WhatsMeowService struct {
	DB     *sql.DB
	Client *whatsmeow.Client
}

func (s *WhatsMeowService) StartClient() {
	deviceStore := s.DeviceStore()
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	s.Client = whatsmeow.NewClient(deviceStore, clientLog)
	s.Client.AddEventHandler(s.EventHandler)

	if s.Client.Store.ID == nil {
		qrChan, _ := s.Client.GetQRChannel(context.Background())
		err := s.Client.Connect()
		if err != nil {
			panic(err)
		}
		for qr := range qrChan {
			if qr.Event == "code" {
				image, _ := qrcode.Encode(qr.Code, qrcode.Medium, 256)
				base64qrcode := "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)
				// show qr interminal with image format
				config := qrterminal.Config{
					Level:     qrterminal.L,
					Writer:    os.Stdout,
					BlackChar: qrterminal.WHITE,
					WhiteChar: qrterminal.BLACK,
					QuietZone: 1,
				}
				qrterminal.GenerateWithConfig(qr.Code, config)
				// save qrcode to database
				sqlStmt := `UPDATE users SET qrcode=? WHERE id=1`
				s.DB.Exec(sqlStmt, base64qrcode)
			} else {
				// other event
				fmt.Println("event nya : ", qr.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err := s.Client.Connect()
		if err != nil {
			panic(err)
		}
	}
}

func (s *WhatsMeowService) DeviceStore() *store.Device {
	var deviceStore *store.Device
	var err error
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite", "file:datastore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	deviceStore, err = container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	return deviceStore
}

func (s *WhatsMeowService) EventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		// message := v.Message.GetConversation()
	case *events.Connected:
		sqlStmt := `UPDATE users SET connected=1 WHERE id=?`
		s.DB.Exec(sqlStmt, 1)
	case *events.Disconnected:
		sqlStmt := `UPDATE users SET connected=0 WHERE id=?`
		s.DB.Exec(sqlStmt, 1)
	case *events.PairSuccess:
		sqlStmt := `UPDATE users SET jid=?, qrcode=? WHERE id=?`
		s.DB.Exec(sqlStmt, v.ID, "", 1)
	}
}

func (s *WhatsMeowService) GetQR() (string, string) {
	sqlStmt := `SELECT qrcode, jid FROM users WHERE id=?`
	var qrcode string
	var jid string
	s.DB.QueryRow(sqlStmt, 1).Scan(&qrcode, &jid)
	return qrcode, jid
}

func (s *WhatsMeowService) SendTextMessage(to string, message string) error {
	if to[0:1] == "+" {
		to = to[1:]
	}
	if to[0:1] == "0" {
		to = "62" + to[1:]
	}
	_, err := s.Client.SendMessage(context.Background(), types.JID{User: to, Server: "s.whatsapp.net"}, &waProto.Message{
		Conversation: proto.String(message),
	})

	return err
}
