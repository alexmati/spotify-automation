package handler

import "net/http"

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := Templates.ExecuteTemplate(w, "welcome.html", nil); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		//CHANGE TO LOGIN HANDLER
	}
}

//func CallbackHandler(w http.ResponseWriter, r *http.Request) {
//
//}
