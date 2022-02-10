package main

import (
	"io/ioutil"
	"net/http"
	"fmt"
	"strings"
	"net/url"
	"time"
	"github.com/eiannone/keyboard"	
	"os/user"
	"sync"
)

var (
	saladechat []byte
	urlchat = "http://ravr.webcindario.com/5_chat/consola.php"
	usuario string
	mensaje string
)





func routine(command <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	var status = "Play"
	for {
		select {
		case cmd := <-command:
			//fmt.Println(cmd)
			switch cmd {
			case "Stop":
				return
			case "Pause":
				status = "Pause"
			default:
				status = "Play"
			}
		default:
			if status == "Play" {
				
				
			
					texto, haynuevo := Mensajeria(urlchat,"","")
					if haynuevo {
						for i:=0; i<(len(usuario) + len(">") + len(mensaje)); i++ {
							fmt.Print("\b \b")
						}
						fmt.Print("\r")
						fmt.Println(texto)
						fmt.Printf("%s>%s", usuario, mensaje)
					}
					
				
					
					time.Sleep(1 * time.Second)
					
				
				
				
				
			}
		}
	}
}



func main() {
		
		
		//urlchat = "http://ravr.webcindario.com/5_chat/consola.php"
		
		user, error1 := user.Current()
		usuario = user.Username
		
		usuario = strings.Replace(usuario,"\n","\\n",-1)
		
		mensaje = ""
		
		
		if ConexionValida(urlchat) && error1==nil {
			
			
			
			
			
			var wg sync.WaitGroup
			wg.Add(1)
			command := make(chan string)
			go routine(command, &wg)
			
			

	
			
			fmt.Printf("%s>", usuario)
	
	
	
			if err := keyboard.Open(); err != nil {
				panic(err)
			}
			defer func() {
				_ = keyboard.Close()
			}()

			//fmt.Println("Press ESC to quit")
			for {
				char, key, err := keyboard.GetKey()
				if err != nil {
					panic(err)
				}
				if key == 0 {
					mensaje = mensaje + string(char)
					fmt.Printf("%c",char)
				} else if key == keyboard.KeySpace {
					mensaje = mensaje + " "
					fmt.Print(" ")
				} else if key == keyboard.KeyEsc {
					break
				} else if key == keyboard.KeyBackspace || key == keyboard.KeyBackspace2 {
					if len(mensaje) != 0 {
						mensaje = mensaje[:len(mensaje)-1]
						fmt.Print("\b \b")
					}
					
				} else if key == keyboard.KeyEnter && len(mensaje)>0 {
					//salir = true
					command <- "Pause"
					
					/*command <- "Stop"
					wg.Wait()*/
					
					Enviar(urlchat, usuario, mensaje)
					
					command <- "Play"
				}
			}
			
			//KeyEnter
			
			//time.Sleep(10 * time.Second)
			//ejecutamos peticionador de caracteres en paralelo
			
		} else {
			fmt.Println("Conexion Invalida")
		}
		
		
		
		
		
		
		
}




func ConexionValida(link string) bool {
	respuest, errorlink := http.Get(link) 
	if errorlink != nil {
		return false
	}
	defer respuest.Body.Close()
	if respuest.StatusCode != 200 {
		return false
	}
	saladechat, errorlink = ioutil.ReadAll(respuest.Body)
	if errorlink != nil {
		return false
	}
	return true
}


// retorno las lineas nuevas + true si hay dato, false si no hay nada
func Mensajeria(link, nombre, mensaje string) (string, bool) {
	
	linkparseado, error1 := url.Parse(link)
    if error1 != nil {
        return "",false
    }
    parametros := url.Values{}
    parametros.Add("nombre", nombre)
    parametros.Add("mensaje", mensaje)
    linkparseado.RawQuery = parametros.Encode()
    
	
	respuesta, error2 := http.Get(linkparseado.String()) 
	if error2 != nil {
		return "",false
	}
	defer respuesta.Body.Close()
	if respuesta.StatusCode != 200 {
		return "",false
	}
	actualizado, error3 := ioutil.ReadAll(respuesta.Body)
	if error3 != nil {
		return "",false
	}
	

	
	
		lineas_nuevas := strings.Split(string(actualizado) , "\n")
		//lineas_viejas := strings.Split(old , "\n")
		lineas_viejas := strings.Split(string(saladechat) , "\n")
		
		imprimir := false
		imprimirnada := false
		
		i := 0
		j := 0
		
		retorno := ""
		
		if len(string(actualizado)) > 0 {
			for i=0; i<len(lineas_nuevas); i++ {
				campos_nuevos := strings.Split(lineas_nuevas[i] , ";")
				if imprimir==false {
					for j=0; j<len(lineas_viejas) && imprimir==false; j++ {
						campos_viejos := strings.Split(lineas_viejas[j] , ";")
						if campos_nuevos[0] == campos_viejos[0] {
							imprimir = true
							if j == 0 {
								imprimirnada = true
							}
						} 
					}
					if imprimirnada == true {
						break
					} else if imprimir == true {
						diferencia := len(lineas_viejas) - j
						i = diferencia + 1
					} else if j == len(lineas_viejas) {
						imprimir=true
					}				
				}
				if imprimir==true {
					campos_nuevos = strings.Split(lineas_nuevas[i] , ";")
					if len(retorno) == 0 {
						retorno = campos_nuevos[1]+": "+campos_nuevos[2]
					} else {
						retorno = retorno +"\n"+campos_nuevos[1]+": "+campos_nuevos[2]
					} 
				}
			}
		}
		
		if imprimirnada == true {
			return "",false
		} else if imprimir==true {
			saladechat = actualizado
			return retorno,true
		} else {
			return "",false
		}
}





func Enviar(link, nombre, mensaje string) bool {
	
	linkparseado, error1 := url.Parse(link)
    if error1 != nil {
        return false
    }
    parametros := url.Values{}
    parametros.Add("nombre", nombre)
    parametros.Add("mensaje", mensaje)
    linkparseado.RawQuery = parametros.Encode()
    
	
	respuesta, error2 := http.Get(linkparseado.String()) 
	if error2 != nil {
		return false
	}
	defer respuesta.Body.Close()
	if respuesta.StatusCode != 200 {
		return false
	}
	_, error3 := ioutil.ReadAll(respuesta.Body)
	if error3 != nil {
		return false
	}
	return true
	
}






