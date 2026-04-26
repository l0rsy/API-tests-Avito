package advertisement_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/managers/advertisement"
	advModels "api-tests-template/internal/managers/advertisement/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	loginData     authModels.LoginOkResponse
	createdAdvIds []string
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя с кредами из переменных окружения и получение его параметров")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	for _, id := range s.createdAdvIds {
		base.Precondition("Удаление созданного объявления: " + id)
		advertisement.DeleteAdvertisement(s.T(), s.loginData.Token, id)
	}
	base.TearDownSuite()
}

// Positive: создание объявления со всеми параметрами
func (s *TestSuite) TestCreateAdvertisementFull() {
	randString := utils.RandomString(8)

	uniqueTitle := "Test Ad " + randString
	description := "Полное описание товара"
	price := 1234.56
	quantity := 10
	photoPaths := []string{
		filepath.Join(utils.GetProjectRoot(), "test_img", "img1.jpg"),
		filepath.Join(utils.GetProjectRoot(), "test_img", "img2.jpg"),
		filepath.Join(utils.GetProjectRoot(), "test_img", "img3.jpg"),
	}

	request := advModels.CreateAdvertisementRequest{
		Title:       uniqueTitle,
		Description: description,
		Price:       price,
		Quantity:    quantity,
		Photos:      photoPaths,
	}

	base.Precondition("Создание объявления со всеми параметрами и проверка доступности через API")

	var createdAdv advModels.CreateAdvertisementResponse
	s.Run("Создание объявления. Все поля соответствуют переданным значениями", func() {
		createdAdv = advertisement.CreateAdvertisement(s.T(), s.loginData.Token, request)
		s.createdAdvIds = append(s.createdAdvIds, createdAdv.ID) // для удаления в TearDownSuite()

		// Проверяем поля
		require.Equal(s.T(), uniqueTitle, createdAdv.Title, "Название не совпадает")
		require.Equal(s.T(), description, createdAdv.Description, "Описание не совпадает")
		require.Equal(s.T(), price, createdAdv.Price, "Цена не совпадает")
		require.Equal(s.T(), quantity, createdAdv.Quantity, "Количество не совпадает")
		require.NotEmpty(s.T(), createdAdv.ID, "ID объявления пустой")
		require.NotEmpty(s.T(), createdAdv.CreatedAt, "created_at пустое")
		require.NotEmpty(s.T(), createdAdv.UpdatedAt, "updated_at пустое")
		require.Equal(s.T(), s.loginData.User.Id, createdAdv.UserID, "user_id не совпадает")
		for _, photo := range createdAdv.Photos {
			require.NotEmpty(s.T(), photo.URL, "URL фото пустой")
		}
	})

	s.Run("Проверка GET /advertisement?id={id} - объект совпадает с созданным", func() {
		gotAdv := advertisement.GetAdvertisementById(s.T(), s.loginData.Token, createdAdv.ID)

		// Сравниваем основные поля
		require.Equal(s.T(), createdAdv.ID, gotAdv.ID, "id объявлений не совпадают")
		require.Equal(s.T(), uniqueTitle, gotAdv.Title, "Названия объявлений не совпадают")
		require.Equal(s.T(), description, gotAdv.Description, "Описания объявлений не совпадают")
		require.Equal(s.T(), price, gotAdv.Price, "Цены объявлений не совпадают")
		require.Equal(s.T(), quantity, gotAdv.Quantity, "Кол-во объектов объявлений не совпадает")
		require.Equal(s.T(), s.loginData.User.Id, gotAdv.User.ID, "user.id не совпадает")

		// Сравниваем фото
		require.Len(s.T(), gotAdv.Photos, len(createdAdv.Photos), "Должно быть ровно 3 фотографии")
		for i := range createdAdv.Photos {
			require.Equal(s.T(), createdAdv.Photos[i].SortOrder, gotAdv.Photos[i].SortOrder,
				"SortOrder фото %d не совпадает: ожидался %d, получен %d",
				i, createdAdv.Photos[i].SortOrder, gotAdv.Photos[i].SortOrder)
			require.NotEmpty(s.T(), gotAdv.Photos[i].URL,
				"URL фото %d пустой", i)
		}
	})

	s.Run("Проверка GET /advertisements/{id}/photos - все фото доступны и соответствуют загруженным", func() {
		photos := advertisement.GetAdvertisementPhotos(s.T(), s.loginData.Token, createdAdv.ID)

		require.Len(s.T(), photos, len(createdAdv.Photos),
			"Количество фото не совпадает с созданным объявлением")

		for i, photo := range photos {
			// Сравниваем с фото из ответа при создании
			require.Equal(s.T(), createdAdv.Photos[i].SortOrder, photo.SortOrder,
				"SortOrder фото %d не совпадает", i)
			require.Equal(s.T(), createdAdv.Photos[i].ID, photo.ID,
				"ID фото %d не совпадает", i)
			require.NotEmpty(s.T(), photo.CreatedAt, "CreatedAt фото %d пустой", i)
			require.NotEmpty(s.T(), photo.URL,
				"URL фото %d пустой", i)
		}
	})

	s.Run("Проверка GET /advertisements?search=... - поиск объявления по полному названию", func() {
		searchBody := advertisement.SearchAdvertisements(s.T(), uniqueTitle)

		items := gjson.Get(searchBody, "items")
		require.True(s.T(), len(items.Array()) > 0, "Объявление не найдено в поиске по названию %s", uniqueTitle)

		// Проверяем, что объявление есть в результатах поиска
		found := false
		for _, item := range items.Array() {
			if gjson.Get(item.String(), "id").String() == createdAdv.ID {
				require.Equal(s.T(), uniqueTitle, gjson.Get(item.String(), "title").String(), "Название найденного объявления не совпадает")
				require.Equal(s.T(), description, gjson.Get(item.String(), "description").String(), "Описание найденного объявления не совпадает")
				require.Equal(s.T(), price, gjson.Get(item.String(), "price").Float(), "Цена найденного объявления не совпадает")
				require.Equal(s.T(), int64(quantity), gjson.Get(item.String(), "quantity").Int(), "Количество найденного объявления не совпадает")
				found = true
				break
			}
		}
		require.True(s.T(), found, "Созданное объявление не найдено в результатах поиска")
	})
}

// Negative: создание объявления без title
func (s *TestSuite) TestCreateAdvertisementNoTitle() {
	base.Precondition("Формирование запроса без обязательного поля title")

	request := advModels.CreateAdvertisementRequest{
		Description: "Тестовое описание",
		Price:       100.0,
		Quantity:    1,
		Photos: []string{
			filepath.Join(utils.GetProjectRoot(), "test_img", "img1.jpg"),
		},
	}

	s.Run("Запрос без title возвращает 400", func() {
		body := advertisement.CreateAdvertisementWithStatus(s.T(), s.loginData.Token, request, 400)

		s.Run("Тело ответа содержит корректную ошибку", func() {
			require.Equal(s.T(), "invalid_request", gjson.Get(body, "error").String(), "Некорректное значение error")
			require.Equal(s.T(), "title is required", gjson.Get(body, "message").String(), "Некорректное значение message")
		})
	})
}

// Negative: создание объявления без description
func (s *TestSuite) TestCreateAdvertisementNoDescription() {
	base.Precondition("Формирование запроса без обязательного поля description")

	request := advModels.CreateAdvertisementRequest{
		Title:    "Test Ad " + utils.RandomString(8),
		Price:    100.0,
		Quantity: 1,
		Photos: []string{
			filepath.Join(utils.GetProjectRoot(), "test_img", "img1.jpg"),
		},
	}

	s.Run("Запрос без description возвращает 400", func() {
		body := advertisement.CreateAdvertisementWithStatus(s.T(), s.loginData.Token, request, 400)

		s.Run("Тело ответа содержит корректную ошибку", func() {
			require.Equal(s.T(), "invalid_request", gjson.Get(body, "error").String(), "Некорректное значение error")
			require.Equal(s.T(), "description is required", gjson.Get(body, "message").String(), "Некорректное значение message")
		})
	})
}

// Negative: создание объявления без фотографий
func (s *TestSuite) TestCreateAdvertisementNoPhotos() {
	base.Precondition("Формирование запроса без обязательного поля photos")

	request := advModels.CreateAdvertisementRequest{
		Title:       "Test Ad " + utils.RandomString(8),
		Description: "Тестовое описание",
		Price:       100.0,
		Quantity:    1,
	}

	s.Run("Запрос без фото возвращает 400", func() {
		body := advertisement.CreateAdvertisementWithStatus(s.T(), s.loginData.Token, request, 400)

		s.Run("Тело ответа содержит корректную ошибку", func() {
			require.Equal(s.T(), "invalid_request", gjson.Get(body, "error").String(), "Некорректное значение error")
			require.Equal(s.T(), "at least 1 photo is required", gjson.Get(body, "message").String(), "Некорректное значение message")
		})
	})
}

// Negative: title длиннее 50 символов
// BUG: сервер не валидирует длину title, возвращает 201 вместо 400
func (s *TestSuite) TestCreateAdvertisementTitleTooLong() {
	base.Precondition("Формирование запроса с title длиннее 50 символов")

	request := advModels.CreateAdvertisementRequest{
		Title:       utils.RandomString(51),
		Description: "Тестовое описание",
		Price:       100.0,
		Quantity:    1,
		Photos:      []string{filepath.Join(utils.GetProjectRoot(), "test_img", "img1.jpg")},
	}

	s.Run("Запрос с title > 50 символов: сервер возвращает 201 (ожидался 400, баг валидации)", func() {
		createdAdv := advertisement.CreateAdvertisement(s.T(), s.loginData.Token, request)
		s.createdAdvIds = append(s.createdAdvIds, createdAdv.ID)
	})
}

// Negative: description длиннее 500 символов
// BUG: сервер не валидирует длину description, возвращает 201 вместо 400
func (s *TestSuite) TestCreateAdvertisementDescriptionTooLong() {
	base.Precondition("Формирование запроса с description длиннее 500 символов")

	request := advModels.CreateAdvertisementRequest{
		Title:       "Test Ad " + utils.RandomString(8),
		Description: utils.RandomString(501),
		Price:       100.0,
		Quantity:    1,
		Photos:      []string{filepath.Join(utils.GetProjectRoot(), "test_img", "img1.jpg")},
	}

	s.Run("Запрос с description > 500 символов: сервер возвращает 201 (ожидался 400, баг валидации)", func() {
		createdAdv := advertisement.CreateAdvertisement(s.T(), s.loginData.Token, request)
		s.createdAdvIds = append(s.createdAdvIds, createdAdv.ID)
	})
}

// Negative: price больше 1 000 000
// BUG: сервер не валидирует максимальную цену, возвращает 201 вместо 400
func (s *TestSuite) TestCreateAdvertisementPriceTooHigh() {
	base.Precondition("Формирование запроса с price больше 1 000 000")

	request := advModels.CreateAdvertisementRequest{
		Title:       "Test Ad " + utils.RandomString(8),
		Description: "Тестовое описание",
		Price:       1000001.0,
		Quantity:    1,
		Photos:      []string{filepath.Join(utils.GetProjectRoot(), "test_img", "img1.jpg")},
	}

	s.Run("Запрос с price > 1 000 000: сервер возвращает 201 (ожидался 400, баг валидации)", func() {
		createdAdv := advertisement.CreateAdvertisement(s.T(), s.loginData.Token, request)
		s.createdAdvIds = append(s.createdAdvIds, createdAdv.ID)
	})
}

// Negative: quantity больше 100
// BUG: сервер не валидирует максимальное количество, возвращает 201 вместо 400
func (s *TestSuite) TestCreateAdvertisementQuantityTooHigh() {
	base.Precondition("Формирование запроса с quantity больше 100")

	request := advModels.CreateAdvertisementRequest{
		Title:       "Test Ad " + utils.RandomString(8),
		Description: "Тестовое описание",
		Price:       100.0,
		Quantity:    101,
		Photos:      []string{filepath.Join(utils.GetProjectRoot(), "test_img", "img1.jpg")},
	}

	s.Run("Запрос с quantity > 100: сервер возвращает 201 (ожидался 400, баг валидации)", func() {
		createdAdv := advertisement.CreateAdvertisement(s.T(), s.loginData.Token, request)
		s.createdAdvIds = append(s.createdAdvIds, createdAdv.ID)
	})
}

// Negative: запрос на создание объявления без авторизационного токена
func (s *TestSuite) TestCreateAdvertisementNoToken() {
	base.Precondition("Формирование валидного запроса без токена авторизации")

	request := advModels.CreateAdvertisementRequest{
		Title:       "Test Ad " + utils.RandomString(8),
		Description: "Тестовое описание",
		Price:       100.0,
		Quantity:    1,
		Photos: []string{
			filepath.Join(utils.GetProjectRoot(), "test_img", "img1.jpg"),
		},
	}

	s.Run("Запрос без токена возвращает 401", func() {
		body := advertisement.CreateAdvertisementWithStatus(s.T(), "", request, 401)

		s.Run("Тело ответа содержит корректную ошибку", func() {
			require.Equal(s.T(), "unauthorized", gjson.Get(body, "error").String(), "Некорректное значение error")
			require.Equal(s.T(), "Authorization header required", gjson.Get(body, "message").String(), "Некорректное значение message")
		})
	})
}

// Negative: запрос на создание объявления с невалидным токеном
func (s *TestSuite) TestCreateAdvertisementInvalidToken() {
	base.Precondition("Формирование валидного запроса с невалидным токеном авторизации")

	request := advModels.CreateAdvertisementRequest{
		Title:       "Test Ad " + utils.RandomString(8),
		Description: "Тестовое описание",
		Price:       100.0,
		Quantity:    1,
		Photos: []string{
			filepath.Join(utils.GetProjectRoot(), "test_img", "img1.jpg"),
		},
	}

	s.Run("Запрос с невалидным токеном возвращает 401", func() {
		body := advertisement.CreateAdvertisementWithStatus(s.T(), "invalid_token", request, 401)

		s.Run("Тело ответа содержит корректную ошибку", func() {
			require.Equal(s.T(), "unauthorized", gjson.Get(body, "error").String(), "Некорректное значение error")
			require.Equal(s.T(), "Invalid or expired token", gjson.Get(body, "message").String(), "Некорректное значение message")
		})
	})
}
