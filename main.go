package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func apiHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		if getUrl := req.URL.Path; getUrl == "/" {
			t := time.Now()
			fmt.Fprintf(w, "%s", t.Format("03h04"))
		} else if getUrl := req.URL.Path; getUrl == "/entries" {
			data, err := ioutil.ReadFile("./save.data")
			if err != nil {
				fmt.Println("File reading error", err)
				fmt.Fprintf(w, "Une erreur est survenue : Impossible de lire les données")
				return
			}
			stringSlice := strings.Split(string(data), ":")
			for index, elem := range stringSlice {
				if index%2 != 0 {
					fmt.Fprintf(w, "%s\n", elem)
				}
			}
		} else {
			fmt.Fprintf(w, "Erreur 404 : la page n'existe pas")
		}
	case http.MethodPost:
		if err := req.ParseForm(); err != nil {
			fmt.Println("Something went bad")
			fmt.Fprintln(w, "Something went bad")
		}

		for key, value := range req.PostForm {
			fmt.Println(key, ":", value)
		}

		fmt.Fprintf(w, "Information received:\nAuthor : %v\nEntry : %v", req.PostForm["author"][0], req.PostForm["entry"][0])

		saveFile, err := os.OpenFile("./save.data", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		defer saveFile.Close()

		w := bufio.NewWriter(saveFile)
		if err != nil {
			fmt.Println("Erreur")
		}

		fmt.Fprintf(w, "%v:%v:\n", req.PostForm["author"][0], req.PostForm["entry"][0])
		w.Flush()
	}
}

func main() {
	http.HandleFunc("/", apiHandler)
	http.HandleFunc("/add", apiHandler)
	http.HandleFunc("/entries", apiHandler)
	http.ListenAndServe(":4567", nil)
}
