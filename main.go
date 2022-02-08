package main

import (
	"io/ioutil"
	"net/http"
	"fmt"
	"strings"
)

var saladechat []byte

func main() {
	
		
		
		//ejecutamos peticionador de caracteres en paralelo
		//ejecutamos peticionador actualizador de char en paralelo
		
		data, haynuevo := DiferenciarChat("http://ravr.webcindario.com/5_chat/consola.php?nombre=&mensaje=")
		
		if haynuevo {
			fmt.Println(data)
		}
		
		
}

/*
parsear texto
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
}*/

//true si hubo error o false sin errores
func ValidarURL(url string) bool {
	resp, err := http.Get(url) 
	if err != nil {
		return true
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return true
	}
	saladechat, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return true
	}
	return false
}


// retorno las lineas nuevas + true si hay dato, false si no hay nada
func DiferenciarChat(url string) (string, bool) {
	
	resp, err := http.Get(url) 
	if err != nil {
		return "",false
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "",false
	}
	actualizado, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "",false
	}
	
	//fmt.Println(string(actualizado))
	
	old := `46;Bry;Hola
47;Rick;Bry???
48;Demo;Prueba de Chat
49;Ricardo;Ultimo Mensaje
50;Jorge;Hola`
	
	
		lineas := strings.Split(string(actualizado) , "\n")
		lineasold := strings.Split(old , "\n")
		//lineasold := strings.Split(string(saladechat) , "\n")
		
		imprimir := false
		imprimirnada := false
		i := 0
		j := 0
		retorno := ""
		
		if 0 != len(string(actualizado)) {
		
			for i=0; i<len(lineas); i++ {
				
				campos := strings.Split(lineas[i] , ";")
				
				if imprimir==false {
				for j=0; j<len(lineasold) && imprimir==false; j++ {
					
					camposold := strings.Split(lineasold[j] , ";")
					
					if campos[0] == camposold[0] {
						imprimir = true
						if j == 0 {
							imprimirnada = true
						}
					} 
				}
				if imprimirnada == true {
					break
				} else if imprimir==true {
					diferencia := len(lineasold) - j
					i = diferencia + 1
					//fmt.Println(i)
				}				
				}
				
				if imprimir==true {
					campos = strings.Split(lineas[i] , ";")
					//fmt.Println(campos[0],"-",campos[1],": ",campos[2])
					if len(retorno) == 0 {
						retorno = campos[0]+"-"+campos[1]+": "+campos[2]
					} else {
						retorno = retorno +"\n"+campos[0]+"-"+campos[1]+": "+campos[2]
					} 
					
				}
				
			}
			
			
		}
	
	if imprimirnada == true {
		return "",false
	} else if imprimir==true {
		return retorno,true
	} else {
		return "",false
	}
	
	
	
	
	
		
}


