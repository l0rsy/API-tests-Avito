package advertisement

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisement"
	"api-tests-template/internal/managers/advertisement/models"
)

// CreateAdvertisement создает объявление
func CreateAdvertisement(t *testing.T, token string, req models.CreateAdvertisementRequest) models.CreateAdvertisementResponse {
	require.NotEmpty(t, token, "Отсутствует токен авторизации")
	require.NotEmpty(t, req.Title, "Отсутствует название объявления")
	require.NotEmpty(t, req.Description, "Отсутствует описание объявления")
	require.GreaterOrEqual(t, req.Price, 0.0, "Цена не может быть отрицательной")
	require.Greater(t, req.Quantity, 0, "Количество должно быть больше 0")
	require.NotEmpty(t, req.Photos, "Должна быть хотя бы одна фотография")

	response := advertisement.HttpPostAdvertisement(t, token, req)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа при создании объявления")

	var adv models.CreateAdvertisementResponse
	err = json.Unmarshal(body, &adv)
	require.NoError(t, err, "Не удалось распарсить JSON ответ при создании объявления")

	return adv
}

// GetAdvertisementById получает объявление по ID
func GetAdvertisementById(t *testing.T, token string, id string) models.GetAdvertisementResponse {
	require.NotEmpty(t, token, "Отсутствует токен авторизации")
	require.NotEmpty(t, id, "Отсутствует ID объявления")

	response := advertisement.HttpGetAdvertisement(t, token, id)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа при получении объявления")

	var adv models.GetAdvertisementResponse
	err = json.Unmarshal(body, &adv)
	require.NoError(t, err, "Не удалось распарсить JSON ответ при получении объявления")

	return adv
}

// DeleteAdvertisement удаляет объявление по ID
func DeleteAdvertisement(t *testing.T, token string, id string) {
	require.NotEmpty(t, token, "Отсутствует токен авторизации")
	require.NotEmpty(t, id, "Отсутствует ID объявления")

	response := advertisement.HttpDeleteAdvertisement(t, token, id)

	require.Equalf(t, http.StatusNoContent, response.StatusCode,
		"Ожидался статус 204 при удалении объявления %s", id)
}

// GetAdvertisementPhotos получает список фотографий объявления по ID
func GetAdvertisementPhotos(t *testing.T, token string, id string) []models.Photo {
	require.NotEmpty(t, token, "Отсутствует токен авторизации")
	require.NotEmpty(t, id, "Отсутствует ID объявления")

	response := advertisement.HttpGetAdvertisementPhotos(t, token, id)

	require.Equalf(t, http.StatusOK, response.StatusCode,
		"Ожидался статус 200 при получении фото объявления %s", id)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа при получении фото")

	var photos []models.Photo
	err = json.Unmarshal(body, &photos)
	require.NoError(t, err, "Не удалось распарсить JSON ответ при получении фото")

	return photos
}

// SearchAdvertisements выполняет поиск объявлений по названию
func SearchAdvertisements(t *testing.T, search string) string {
	require.NotEmpty(t, search, "Отсутствует строка поиска")

	response := advertisement.HttpGetAdvertisementsSearch(t, search)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа при поиске объявлений")

	return string(body)
}

// CreateAdvertisementWithStatus создает объявление и проверяет ожидаемый статус-код
func CreateAdvertisementWithStatus(t *testing.T, token string, req models.CreateAdvertisementRequest, expectedStatus int) string {
	response := advertisement.HttpPostAdvertisementWithStatus(t, token, req, expectedStatus)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа")

	return string(body)
}
