package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/willqiang/bookstore_items-api/domain/items"
	"github.com/willqiang/bookstore_items-api/domain/queries"
	"github.com/willqiang/bookstore_items-api/services"
	"github.com/willqiang/bookstore_items-api/utils/http_utils"
	"github.com/willqiang/bookstore_utils-go/rest_errors"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	ItemsController itemsControllerInterface = &itemsController{}
)

type itemsControllerInterface interface {
	Create(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	Search(http.ResponseWriter, *http.Request)
}

type itemsController struct{}

func (c *itemsController) Create(w http.ResponseWriter, r *http.Request) {
	/*	if err := oauth.AuthenticateRequest(r); err != nil {
			//http_utils.ResponseError(w, err)
			return
		}

		sellerId := oauth.GetCallerId(r)
		if sellerId == 0 {
			requestErr := rest_errors.NewUnauthorizedError("unable to retrieve user infomation from given access token")
			http_utils.ResponseError(w, *requestErr)
			return
		}
	*/
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		requestErr := rest_errors.NewBadRequestError("invalid request body")
		http_utils.ResponseError(w, *requestErr)
		return
	}
	defer r.Body.Close()

	var itemRequest items.Item
	if err := json.Unmarshal(requestBody, &itemRequest); err != nil {
		jsonErr := rest_errors.NewBadRequestError("invalid item json body")
		http_utils.ResponseError(w, *jsonErr)
		return
	}

	itemRequest.Seller = 1 //sellerId
	result, createErr := services.ItemsService.Create(itemRequest)
	if createErr != nil {
		http_utils.ResponseError(w, *createErr)
		return
	}

	http_utils.ResponseJson(w, http.StatusCreated, result)
}

func (c *itemsController) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemId := strings.TrimSpace(vars["id"])

	item, err := services.ItemsService.Get(itemId)
	if err != nil {
		http_utils.ResponseError(w, *err)
		return
	}
	http_utils.ResponseJson(w, http.StatusOK, item)
}

func (c *itemsController) Search(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		reqErr := rest_errors.NewBadRequestError("invalid request body")
		http_utils.ResponseError(w, *reqErr)
		return
	}
	defer r.Body.Close()

	var query queries.EsQuery
	if err := json.Unmarshal(bytes, &query); err != nil {
		jsonErr := rest_errors.NewBadRequestError("invalid json body")
		http_utils.ResponseError(w, *jsonErr)
		return
	}

	items, queryErr := services.ItemsService.Search(query)
	if queryErr != nil {
		http_utils.ResponseError(w, *queryErr)
		return
	}

	http_utils.ResponseJson(w, http.StatusOK, items)
}
