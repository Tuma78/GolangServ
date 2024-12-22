package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"github.com/Tuma78/GolangServ/calc"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}


func CalcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !isValidExpression(request.Expression) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		response := Response{Error: "Expression is not valid"}
		json.NewEncoder(w).Encode(response)
		return
	}

	result, err := calс.Calc(request.Expression)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := Response{Error: "Internal server error"}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Response{Result: fmt.Sprintf("%f", result)}
	json.NewEncoder(w).Encode(response)
}
func isValidExpression(expression string) bool {
	for _, char := range expression {
		if !isValidChar(char) {
			return false
		}
	}
	return true
}

func isValidChar(char rune) bool {
	return (char >= '0' && char <= '9') || char == '+' || char == '-' || char == '*' || char == '/' || char == '(' || char == ')' || char == '.' || char == ' '
}
