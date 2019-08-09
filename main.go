package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"net/url"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"strings"
	"time"
)


func resourceUrl(event *v1.Event) string {
	return os.Getenv("OPENSHIFT_CONSOLE_URL") + "/project/" + event.InvolvedObject.Namespace + "/browse/" + strings.ToLower(event.InvolvedObject.Kind) + "s/" + event.InvolvedObject.Name
}

func monitoringUrl(event *v1.Event) string {
	return os.Getenv("OPENSHIFT_CONSOLE_URL") + "project/" + event.InvolvedObject.Namespace + "/monitoring"
}
func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}
func notifyTelegram(event *v1.Event) {
	telegramApiUrl := getEnv("TELEGRAM_API_URL","https://api.telegram.org")
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	channelName := os.Getenv("TELEGRAM_CHANNEL")

	message := event.InvolvedObject.Namespace+" "+monitoringUrl(event)+"\n" +
		event.InvolvedObject.Name)+" "+ resourceUrl(event)+"\n"+
		event.Message)+"\n"+
		"Reason: " + event.Reason + " Kind: " + event.InvolvedObject.Kind
	params := url.Values{}
	params.Add("chat_id",channelName)
	params.Add("text", message)

	client := http.Client{}

	req, err := http.NewRequest("POST", telegramApiUrl+"/bot"+botToken+"/sendMessage", bytes.NewBufferString(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	_, err = client.Do(req)
	if err != nil {
		fmt.Println("Unable to reach the server.")
	}
}

func watchEvents(clientset *kubernetes.Clientset) {
	startTime := time.Now()
	log.Printf("Watching events after %v", startTime)

	watcher, err := clientset.CoreV1().Events("").Watch(v1.ListOptions{FieldSelector: "type=Warning"})
	if err != nil {
		panic(err.Error())
	}

	for watchEvent := range watcher.ResultChan() {
		event := watchEvent.Object.(*v1.Event)
		if event.FirstTimestamp.Time.After(startTime) {
			notifyTelegram(event)
		}
	}
}

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	go func() {
		for {
			watchEvents(clientset)
			time.Sleep(5 * time.Second)
		}
	}()

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
