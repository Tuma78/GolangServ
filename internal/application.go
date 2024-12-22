package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"github.com/Tuma78/calc"
)

// Config содержит конфигурацию приложения.
type Config struct {
	Addr string
}

// ConfigFromEnv загружает конфигурацию из переменных окружения.
func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

// Application представляет приложение.
type Application struct {
	config *Config
}

// RunServer запускает HTTP-сервер.
func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}

// New создает новое приложение с конфигурацией.
func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

// Request представляет запрос для обработки выражения.
type Request struct {
	Expression string `json:"expression"`
}

// Response представляет ответ на запрос.
type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

// CalcHandler обрабатывает запросы для вычисления выражений.
func CalcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Декодирование JSON в структуру Request
	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверка на недопустимые символы в выражении
	if !isValidExpression(request.Expression) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		response := Response{Error: "Expression is not valid"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Вычисление результата
	result, err := calс.Calc(request.Expression)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := Response{Error: "Internal server error"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Отправка результата пользователю
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Response{Result: fmt.Sprintf("%f", result)}
	json.NewEncoder(w).Encode(response)
}

// isValidExpression проверяет, состоит ли выражение только из допустимых символов.
func isValidExpression(expression string) bool {
	// Проверка на недопустимые символы
	for _, char := range expression {
		if !isValidChar(char) {
			return false
		}
	}
	return true
}

// isValidChar проверяет, является ли символ допустимым для арифметического выражения.
func isValidChar(char rune) bool {
	// Разрешены цифры, операторы и пробелы
	return (char >= '0' && char <= '9') || char == '+' || char == '-' || char == '*' || char == '/' || char == '(' || char == ')' || char == '.' || char == ' '
}
