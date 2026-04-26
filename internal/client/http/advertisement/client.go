package advertisement

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisement/models"
)

// HttpPostAdvertisement отправляет POST запрос на создание объявления
func HttpPostAdvertisement(t *testing.T, token string, req models.CreateAdvertisementRequest) *http.Response {
	request := apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath)

	request = request.
		MultipartFormData("title", req.Title).
		MultipartFormData("description", req.Description).
		MultipartFormData("price", strconv.FormatFloat(req.Price, 'f', -1, 64)).
		MultipartFormData("quantity", strconv.Itoa(req.Quantity))

	for _, filePath := range req.Photos {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Fatalf("Файл не найден: %s", filePath)
		}
		request = request.MultipartFile("photos", filePath)
	}

	return request.
		Expect(t).
		Status(http.StatusCreated).
		End().Response
}

// HttpPostAdvertisementWithStatus отправляет POST запрос на создание объявления и проверяет ожидаемый статус-код
func HttpPostAdvertisementWithStatus(t *testing.T, token string, req models.CreateAdvertisementRequest, expectedStatus int) *http.Response {
	request := apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath)

	request = request.
		MultipartFormData("title", req.Title).
		MultipartFormData("description", req.Description).
		MultipartFormData("price", strconv.FormatFloat(req.Price, 'f', -1, 64)).
		MultipartFormData("quantity", strconv.Itoa(req.Quantity))

	for _, filePath := range req.Photos {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Fatalf("Файл не найден: %s", filePath)
		}
		request = request.MultipartFile("photos", filePath)
	}

	return request.
		Expect(t).
		Status(expectedStatus).
		End().Response
}

// HttpGetAdvertisement выполняет GET запрос на получение объявления по ID
func HttpGetAdvertisement(t *testing.T, token string, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

// HttpDeleteAdvertisement выполняет DELETE запрос на удаление объявления по ID
func HttpDeleteAdvertisement(t *testing.T, token string, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementPhotos выполняет GET запрос на получение списка фотографий объявления
func HttpGetAdvertisementPhotos(t *testing.T, token string, id string) *http.Response {
	photosPath := fmt.Sprintf(path.AdvertisementsPhotosPath, id)
	return apiRunner.GetRunner().Auth(token).Create().Get(photosPath).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementsSearch выполняет GET запрос на ПОИСК объявлений по названию
func HttpGetAdvertisementsSearch(t *testing.T, search string) *http.Response {
	return apiRunner.GetRunner().Create().Get(path.AdvertisementsSearchPath).
		Query("search", search).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}
