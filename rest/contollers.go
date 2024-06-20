package rest

import (
	"net/http"

	"github.com/chapar-rest/mock-server/api"
)

func (ro *Rest) AddPet(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (ro *Rest) UpdatePet(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (ro *Rest) FindPetsByStatus(w http.ResponseWriter, r *http.Request, params api.FindPetsByStatusParams) {
	//TODO implement me
	panic("implement me")
}

func (ro *Rest) FindPetsByTags(w http.ResponseWriter, r *http.Request, params api.FindPetsByTagsParams) {
	//TODO implement me
	panic("implement me")
}

func (ro *Rest) DeletePet(w http.ResponseWriter, r *http.Request, petId int64, params api.DeletePetParams) {
	//TODO implement me
	panic("implement me")
}

func (ro *Rest) GetPetById(w http.ResponseWriter, r *http.Request, petId int64) {
	pet := api.Pet{
		Id:   &petId,
		Name: "dog",
	}

	writeJSON(w, http.StatusOK, pet)
}

func (ro *Rest) UpdatePetWithForm(w http.ResponseWriter, r *http.Request, petId int64, params api.UpdatePetWithFormParams) {
	//TODO implement me
	panic("implement me")
}

func (ro *Rest) UploadFile(w http.ResponseWriter, r *http.Request, petId int64, params api.UploadFileParams) {
	//TODO implement me
	panic("implement me")
}

func (ro *Rest) GetInventory(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (ro *Rest) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (ro *Rest) DeleteOrder(w http.ResponseWriter, r *http.Request, orderId int64) {
	//TODO implement me
	panic("implement me")
}

func (ro *Rest) GetOrderById(w http.ResponseWriter, r *http.Request, orderId int64) {
	//TODO implement me
	panic("implement me")
}

func (ro *Rest) CreateUser(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (ro *Rest) CreateUsersWithListInput(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (ro *Rest) LoginUser(w http.ResponseWriter, r *http.Request, params api.LoginUserParams) {
	//TODO implement me
	panic("implement me")
}

func (ro *Rest) LogoutUser(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (ro *Rest) DeleteUser(w http.ResponseWriter, r *http.Request, username string) {
	//TODO implement me
	panic("implement me")
}

func (ro *Rest) GetUserByName(w http.ResponseWriter, r *http.Request, username string) {
	//TODO implement me
	panic("implement me")
}

func (ro *Rest) UpdateUser(w http.ResponseWriter, r *http.Request, username string) {
	//TODO implement me
	panic("implement me")
}
