package main

import (
	"io"
	"log"
	"net/http"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
)

func main() {
	// make persistent connection
	tx := &smpp.Transceiver{
		Addr: "127.0.0.1:2785",
		//User:   "test_smpp_user",
		//Passwd: "123456",
		User:   "SMPPCMI01",
		Passwd: "pJ2BKFtTj3JeSJss",
	}

	conn := tx.Bind()
	// check initial connection status

	var status smpp.ConnStatus
	if status = <-conn; status.Error() != nil {
		log.Fatalln("Unable to connect, aborting:", status.Error())
	}
	log.Println("Connection completed, status:", status.Status().String())

	// example of connection checker goroutine
	go func() {
		for c := range conn {
			log.Println("SMPP connection status:", c.Status())
		}
	}()
	// example of sender handler func
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sm, err := tx.Submit(&smpp.ShortMessage{
			Src:           "Google",
			Dst:           "12345678911",
			SourceAddrTON: 5,
			SourceAddrNPI: 0,
			DestAddrTON:   1,
			DestAddrNPI:   1,
			Text:          pdutext.Raw("甲骨文甲骨文，是中国的一种古老文字，又称“契文”、“甲骨卜辞”、“殷墟文字”或“龟甲兽骨文”。是我们能见到的最早的成熟汉字，主要是指中国商朝晚期王室用于占卜记事而在龟甲或兽骨上契刻的文字。甲骨文甲骨文，是中国的一种古老文字，又称“契文”、“甲骨卜辞”、"),

			// Text:     pdutext.Raw("甲骨文甲骨文，是中国的一种古老文字，又称“契文”、“甲骨卜辞”、“殷墟文字”或“龟甲兽骨文”。"),
			Register: pdufield.FinalDeliveryReceipt,
			ESMClass: 1,
		})

		if err == smpp.ErrNotConnected {
			http.Error(w, "Oops.", http.StatusServiceUnavailable)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		io.WriteString(w, sm.RespID())
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
