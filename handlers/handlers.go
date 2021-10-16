package handlers



func UserRegistration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	    // create an empty user of type models.User
		var user model.User

		// decode the json request to user
		err := json.NewDecoder(r.Body).Decode(&user)
	
		if err != nil {
			log.Printf("")
		}
	
		// call insert user function and pass the user
		insertID := database.UserCreated(user)
	
		// format a response object
		// res := response{
		//     ID:      insertID,
		//     Message: "Comment created successfully",
		// }
	
		// send the response
		json.NewEncoder(w).Encode(insertID)
}
