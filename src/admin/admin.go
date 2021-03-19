package admin

import (
	"fmt"
)

type Admin struct {
	ID       uint64 `json: id`
	Username string `json: username`
	Email    string `json: email`
	// TODO(Jovan): Hash
	Password string `json: password`
}

// func GetAll(db *database.DBConn) func(http.ResponseWriter, http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		admins := []Admin{}
// 		if err := db.GetAllAdmins(&admins); err != nil {
// 			log.Print(err)
// 		}
// 		json.NewEncoder(w).Encode(admins)
// 	}
// }

func (a Admin) ToString() string {
	return fmt.Sprintf("Admin {ID: %s, Username: %s, Email: %s, Password: %s}",
		a.ID, a.Username, a.Email, a.Password)
}
