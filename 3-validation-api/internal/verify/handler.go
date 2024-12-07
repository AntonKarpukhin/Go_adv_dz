package verify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"sync"
	"time"
	"validation/pkg/account"
	"validation/pkg/file"
	"validation/pkg/request"
	"validation/pkg/utils"
)

type VerifierHandler struct {
	Accounts     []*account.Account // Кэш аккаунтов
	Db           *file.JsFile       // Интерфейс для работы с файлом
	mu           sync.Mutex         // Мьютекс для синхронизации доступа
	stopSaveChan chan struct{}      // Канал для остановки горутины
}

// Валидация данных
var validate = validator.New()

// NewVerifierHandler создает новый обработчик и запускает горутину с возможностью остановки
func NewVerifierHandler(router *http.ServeMux, db *file.JsFile) *VerifierHandler {
	// Читаем данные из файла при старте
	dbData, err := db.Read()
	if err != nil {
		panic(fmt.Errorf("не удалось прочитать данные из файла: %w", err))
	}

	var accountsFromFile []*account.Account
	err = json.Unmarshal(dbData, &accountsFromFile)
	if err != nil {
		panic(fmt.Errorf("не удалось декодировать JSON из файла: %w", err))
	}

	handler := &VerifierHandler{
		Accounts:     accountsFromFile, // Загружаем аккаунты в память
		Db:           db,
		stopSaveChan: make(chan struct{}), // Канал для управления горутиной
	}

	// Регистрируем обработчики маршрутов
	router.HandleFunc("/auth/verify", handler.Register())
	router.HandleFunc("/auth/verify/{hash}", handler.CheckEmail())

	// Запускаем периодическую запись данных
	go handler.periodicSaveToFile()

	return handler
}

// Register обрабатывает POST-запрос для регистрации аккаунтов
func (verifier *VerifierHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload LoginRequest

		// Декодируем JSON-запрос в структуру LoginRequest
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			ErrorsRegister(w, errors.New("некорректный JSON"))
			return
		}

		// Проверяем обязательные поля в запросе
		err = validate.Struct(payload)
		if err != nil {
			ErrorsRegister(w, errors.New("поле email отсутствует или некорректно"))
			return
		}

		// Проверяем, существует ли уже аккаунт с таким email
		verifier.mu.Lock()
		for _, dataAccount := range verifier.Accounts {
			if utils.IsEmail(dataAccount.Email, payload.Email) {
				verifier.mu.Unlock()
				ErrorsRegister(w, errors.New("такой аккаунт уже существует"))
				return
			}
		}

		// Создаем новый аккаунт и добавляем его в кэш
		newAccount := account.NewAccount(payload.Email)
		verifier.Accounts = append(verifier.Accounts, newAccount)
		verifier.mu.Unlock()

		// Отправляем письмо с паролем (эмуляция)
		CreateMail(newAccount.Password, w)

		// Данные будут сохранены в файл автоматически через periodicSaveToFile()
	}
}

// CheckEmail обрабатывает GET-запрос для проверки email по хэшу
func (verifier *VerifierHandler) CheckEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем хэш из маршрута
		hash := r.URL.Query().Get("hash")
		if hash == "" {
			ErrorsRegister(w, errors.New("хэш не предоставлен"))
			return
		}

		// Проверяем аккаунты в кэше и фильтруем их
		verifier.mu.Lock()
		var newAccounts []*account.Account
		for _, dataAccount := range verifier.Accounts {
			if utils.IsEmail(dataAccount.Password, hash) {
				newAccounts = append(newAccounts, dataAccount)
			} else {
				fmt.Println("Аккаунт удален")
			}
		}
		verifier.Accounts = newAccounts
		verifier.mu.Unlock()

		// Данные будут сохранены в файл автоматически через periodicSaveToFile()
	}
}

// periodicSaveToFile выполняет периодическую запись кэша аккаунтов в файл
func (verifier *VerifierHandler) periodicSaveToFile() {
	ticker := time.NewTicker(5 * time.Minute) // Интервал записи данных
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C: // Выполняем запись данных раз в 5 минут
			verifier.mu.Lock()
			marshal, err := json.Marshal(verifier.Accounts)
			if err != nil {
				fmt.Println("Ошибка сериализации аккаунтов:", err)
			} else {
				err = verifier.Db.Write(marshal)
				if err != nil {
					fmt.Println("Ошибка записи аккаунтов в файл:", err)
				}
			}
			verifier.mu.Unlock()
		case <-verifier.stopSaveChan: // Получен сигнал остановки
			fmt.Println("Горутина periodicSaveToFile остановлена")
			return
		}
	}
}

// Stop завершает работу горутины периодической записи
func (verifier *VerifierHandler) Stop() {
	close(verifier.stopSaveChan) // Закрываем канал для завершения горутины
}
