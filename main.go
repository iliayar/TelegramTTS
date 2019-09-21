package main

import (
	"log"
	"os/exec"
	"os"
	"time"
	"errors"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)


func storeFile(name, file string) string {

	cmd := exec.Command("curl", "-F", "name="+name+".ogg",
	 "-F", "file=@" + file, "https://uguu.se/api.php?d=upload-tool")
	b, err := cmd.Output()

	if err != nil {
		log.Panic(err)
	}

	return string(b)

}

func genText(id, text string) error {
	
	if text == "" {
		return errors.New("text string is empty")
	}

	dir, _ := os.Getwd()
	
	BALCON := "C:\\Program Files (x86)\\balcon\\balcon.exe"
	VOICE := "IVONA 2 Maxim OEM"
	FILE :=  dir + "\\data\\" + id + ".ogg"

	cmd := exec.Command(BALCON, "-n", VOICE, "-t", text, "-w", FILE)

	err := cmd.Run()
	if err != nil {
		log.Panic(err)
		return err
	}

	return nil
}

func main() {

	b, err := tb.NewBot(tb.Settings{
		Token: getToken(),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("[Tg] Authorized")

	b.Handle(tb.OnQuery, func(q *tb.Query) {
		if strings.ContainsAny(q.Text,"\"\\/'") {
			log.Panic("Escape cahracters met")
			return
		}

		results := make(tb.Results, 1)

		if strings.Index(q.Text, "...") == -1 {

			results[0] = &tb.ArticleResult{
				Text: "Add \"...\" in the end of the message",
				Title: "Tip",
			}

		} else {

			err = genText(q.From.Username,q.Text)

			if err != nil {
				log.Panic(err)
				return
			}

			dir, _ := os.Getwd()
			url := storeFile(q.From.Username,dir + "\\data\\" + q.From.Username + ".ogg")

			log.Println("[URL] "+url)

			
			results[0] = &tb.VoiceResult{
				URL: url,
				Title: "TTS",
			}
			
		}

		results[0].SetResultID("0")

		err := b.Answer(q, &tb.QueryResponse{
			Results: results,
			CacheTime: 60,
		})

		if err != nil {
			log.Panic(err)
		}
	})

	b.Start()

}