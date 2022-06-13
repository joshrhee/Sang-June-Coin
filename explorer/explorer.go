package explorer

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/BlockChain/blockchain"
)

const (
	templateDir string = "explorer/templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

//rw: A writer that write a data that we want to send to user
//r: a pointer of Request because we are going to use actual Request, not a copy of Request
func homeHandler(rw http.ResponseWriter, r *http.Request) {

	blockChaindata := homeData{"Home", blockchain.GetBlockChain().AllBlocks()}
	//Executing a template which name is "home"
	templates.ExecuteTemplate(rw, "home", blockChaindata)
}

func addHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchain.GetBlockChain().AddBlock(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func Start(port int) {
	//Creating a new Server Mux(Multiplexer) to prevent default Mux run multiple localhost
	handler := http.NewServeMux()

	//templates variable will load all .gohtml template files in "templates/pages/"
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	//Now templates is an object, so we can use templates.ParseGlob for partials folder's files
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))

	handler.HandleFunc("/", homeHandler)
	handler.HandleFunc("/add", addHandler)

	fmt.Printf("Listening on http://localhost %d\n", port)

	//Fatal will print error if there is an error and complete the program
	//os.Exit(1) means program exit with an error code 1
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}
