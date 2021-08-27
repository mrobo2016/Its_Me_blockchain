package main

/*
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"net/http"
)

var test map[string]int
var conf = readConfig()

type Config struct {
	Authorization string `json:"authorization"`
	KasAccount string `json:"kas_account"`
	MasterContract string `json:"master_contract"`
}
 
func init() {
    test = make(map[string]int)
    test["foo"] = 0
    test["bar"] = 1
}

func main() {
    fmt.Println(test) // prints map[foo:0 bar:1]
	fmt.Println(string(conf.Authorization))

	router := gin.Default()

	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
	router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})
	router.Run(":8081")
}

func readConfig() Config {
	var conf Config
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &conf); err != nil {
		panic(err)
	}
	return conf
}
*/ 