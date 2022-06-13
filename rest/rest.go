package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/BlockChain/blockchain"
	"github.com/BlockChain/utils"
	"github.com/gorilla/mux"
)

// const port string = ":4000"
var port string

type url string

//Implemented TextMarshaler interface which is already built in
//If url type is instantiated, this TextMarshaler interface would be called
//so this MarshalText() will be called
func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type urlDescription struct {
	//Using struct field tag to show URL to url
	//because we cannot export lowercase starting variable
	//URL string `json: "url"`   means when my struct is json type,
	//URL would be shown url"

	//omitempty means if that field is empty, not show that field
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type addBlockBody struct {
	Message string
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

// //Implemented Stringer interface which is already built in
// //If a function's name is String(), fmt package will call this Stringer function
//
// func (urlDescription URLDescription) String() string {
// 	return "Hello I am the URL Description"
// }

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/blocks"),
			Method:      "GET",
			Description: "See All Blocks",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{height}"),
			Method:      "GET",
			Description: "See A Block",
		},
	}

	// //Even though we changed the data to jsonFormat in the below,
	// //this is not still actual json, so browser will understand this type as text
	// //This part is telling browser that whatever we are writing in rw
	// //will be "application/json" type for "Content-Type"
	// rw.Header().Add("Content-Type", "application/json")

	// //json.Marshal returns data object into encoded []byte type
	// byteData, err := json.Marshal(data)
	// utils.ErrHandler(err)
	// //If we didn't do rw.Header().Add("Content-Type", "application/json")
	// //Still jsonFormatData's "Content-Type" will be still text
	// jsonFormatData := string(byteData)
	// fmt.Fprint(rw, jsonFormatData)

	//This part is same with right above part
	//NewEncoder writes encoded data to rw
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rw.Header().Add("Content-Type", "application/json ")
		json.NewEncoder(rw).Encode(blockchain.GetBlockChain().AllBlocks())
	case "POST":
		var addBlockBody addBlockBody
		//We put decode version of request's content into our empty variable addBlockBody
		utils.ErrHandler(json.NewDecoder(r.Body).Decode(&addBlockBody))
		blockchain.GetBlockChain().AddBlock(addBlockBody.Message)
		//rw sends http response and 201 status code
		rw.WriteHeader(http.StatusCreated)

	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	//mux.Vars returns current Request as a Map like map[height:2]
	//if API request is like "http://localhost:4000/blocks/2"
	vars := mux.Vars(r)

	//strconv package converts other type to string or string to other type
	//Atoi() converts string to integer
	height, err := strconv.Atoi(vars["height"])
	utils.ErrHandler(err)

	block, err := blockchain.GetBlockChain().GetBlock(height)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
		return
	}
	encoder.Encode(block)
}

//MiddleWare is a function that is called right before the final destination
//Final destination means next router's api endpoint in this case
func jsonContentTypeMiddleWare(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})

}

func Start(aPort int) {
	// //Creating a new Server Mux(Multiplexer) to prevent default Mux run multiple localhost
	// handler := http.NewServeMux()

	//Using Gorila mux's routre instead of Go's builtin ServeMux()
	//Gorila mux has many useful things that ServeMux doesn't have
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", aPort)

	router.Use(jsonContentTypeMiddleWare)

	//.Method() fix our HandleFunc run only the parameter's methods
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")

	//Gorial mux's router find "height: number" in url
	router.HandleFunc("/blocks/{height:[0-9]+}", block).Methods("GET")

	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
