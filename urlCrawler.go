package main

import (
    "fmt"
    "net/url"
)

func main() {
    base, err := url.Parse("http://ravr.webcindario.com/5_chat/servidor.php")
    if err != nil {
        return
    }

    // Query params
    params := url.Values{}
    params.Add("nombre", "ricky")
    params.Add("mensaje", "El niño moreno")
    base.RawQuery = params.Encode() 

    fmt.Printf("Encoded URL is %q\n", base.String())
}

/*
package main

import (
	//"html"
	"io/ioutil"
	//"log"
	"net/http"
	//"net/url"
	//"regexp"
	"fmt"
	"strings"
)

func main() {

	//crawlUrl("http://ravr.webcindario.com/5_chat/servidor.php?nombre=&mensaje=")
	
	var (
		err     error
		content []byte
		resp    *http.Response
	)

	// GET content of URL
	resp, err = http.Get("http://ravr.webcindario.com/5_chat/consola.php",url.Values{"nombre": {"Value"}, "mensaje": {"123"}}) 
	if err != nil {
		fmt.Println("ERROR GET")
		return
	}
	defer resp.Body.Close()

	// Check if request was successful
	if resp.StatusCode != 200 {
		fmt.Println("ERROR EN SITIO WEB")
		return
	}

	// Read the body of the HTTP response
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR EN SITIO WEB")
		return
	}


	lineas := strings.Split(string(content) , "\n")
	for i:=0; i<len(lineas); i++ {
		campos := strings.Split(lineas[i] , ";")
		fmt.Println(campos[1],":",campos[2])
	}
	
	
}*/


