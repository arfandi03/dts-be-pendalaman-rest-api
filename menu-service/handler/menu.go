package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/arfandidts/dts-be-pendalaman-microservice/menu-service/database"
	"github.com/arfandidts/dts-be-pendalaman-microservice/utils"
	"github.com/gorilla/context"
	"gorm.io/gorm"
)

type Menu struct {
	Db *gorm.DB
}

// AddMenuHandler handle add menu
func (handler *Menu) AddMenu(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.WrapAPIError(w, r, "can't read body", http.StatusBadRequest)
		return
	}

	// mengambil data yang dikirim dari handler sebelumnya
	username := context.Get(r, "user")

	var menu database.Menu
	err = json.Unmarshal(body, &menu)
	if err != nil {
		utils.WrapAPIError(w, r, "error unmarshal : "+err.Error(), http.StatusInternalServerError)
		return
	}

	menu.Username = fmt.Sprintf("%v", username)

	err = menu.Insert(handler.Db)
	if err != nil {
		utils.WrapAPIError(w, r, "insert menu error : "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WrapAPISuccess(w, r, "success", 200)
}

func (menu *Menu) GetAllMenu(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	menuDb := database.Menu{}

	menus, err := menuDb.GetAll(menu.Db)
	if err != nil {
		utils.WrapAPIError(w, r, "failed get menu:"+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WrapAPIData(w, r, menus, 200, "success")
}
