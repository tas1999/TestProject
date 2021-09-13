package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/tas1999/TestProject/docs"
)

var (
	// DSN это соединение с базой
	// вы можете изменить этот на тот который вам нужен
	// docker run -p 3306:3306 -v $(PWD):/docker-entrypoint-initdb.d -e MYSQL_ROOT_PASSWORD=1234 -e MYSQL_DATABASE=golang -d mysql
	// DSN = "root@tcp(localhost:3306)/golang2017?charset=utf8"
	//DSN = "User=postgres Password=example Server=localhost Port=5432 Database=testgo sslmode=disable"
	DSN = "user=postgres password=example dbname=testgo port=5432 host=localhost sslmode=disable"
)

type Player struct {
	Id    int
	Name  string
	Email string
	Age   int
}

// @title Swagger TestProject API
// @version 1.0
// @description This is a sample server Petstore TestProject server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8082
// @BasePath /
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationurl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information

// @x-extension-openapi {"example": "value on a json format"}
func main() {

	db, err := sql.Open("postgres", DSN)
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		panic(err)
	}

	handler, err := NewDbExplorer(db)
	if err != nil {
		panic(err)
	}
	serveMux := http.NewServeMux()
	serveMux.Handle("/", handler)
	serveMux.Handle("/swagger/", Swagger())
	fmt.Println("starting server at :8082")
	handler = panicMiddleware(serveMux)

	http.ListenAndServe(":8082", handler)
}

func Swagger() http.Handler {
	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8082/swagger/doc.json"),
	))
	return r
}
func NewDbExplorer(db *sql.DB) (http.Handler, error) {
	dbEx := DbExplorer{Db: db}
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/players", dbEx.List)
	return serveMux, nil
}
func panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("panicMiddleware", r.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recovered", err)
				http.Error(w, "Internal server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type DbExplorer struct {
	Db *sql.DB
}

// List players
// @Summary Show a players
// @Description get list players
// @ID get-list-players
// @Accept  json
// @Produce  json
// @Success 200 {object} []Player
// @Router /players [get]
func (ex *DbExplorer) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := ex.Db.Query("Select * from players")
	__err_panic(err)
	var players []Player
	for rows.Next() {
		var pl Player
		err = rows.Scan(&pl.Id, &pl.Name, &pl.Email, &pl.Age)
		players = append(players, pl)
	}

	__err_panic(err)
	err = json.NewEncoder(w).Encode(players)
	__err_panic(err)
}

// Add player
// @Summary Add player
// @Description add player
// @ID add-player
// @Accept json
// @Produce json
// @Param player body Player true "Add player"
// @Success 200 {object} int
// @Router /player [post]
func (ex *DbExplorer) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var player Player
	err := json.NewDecoder(r.Body).Decode(&player)
	__err_panic(err)
	ret := ex.Db.QueryRow("INSERT INTO Players (name,email,age) VALUES ('?', '?', ?) return id", player.Name, player.Email, player.Age)
	__err_panic(err)
	var id int64
	ret.Scan(&id)

	__err_panic(err)
	err = json.NewEncoder(w).Encode(id)
	__err_panic(err)
}
func __err_panic(err error) {
	if err != nil {
		panic(err)
	}
}
