package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		encrypted, err := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(encrypted))
	})

	log.Println("서버 시작 : 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
