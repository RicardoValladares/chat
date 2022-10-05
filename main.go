package main

import (
	"github.com/eiannone/keyboard"	
	"fmt"
	"strings"
	"net/http"
	"net/url"
	"io/ioutil"
	"os/user"
	"time"
	"sync"
)

var (
	urlchat = "http://ravr.webcindario.com/5_chat/consola.php"
	usuario string
	mensaje string
	saladechat []byte
)


func main() {
	user, errorfatal := user.Current()
	usuario = user.Username //Usuario del Sistema Operativo
	mensaje = ""			//Mensaje escrito desde Consola
	if ConexionValida() && errorfatal == nil {
			
		/* Ejecutamos Go Rutina, proceso en paralelo */
		var wg sync.WaitGroup
		wg.Add(1)
		go_comando := make(chan string)
		go Routine(go_comando, &wg)
		
		/* Verificamos podemos obtener peticiones desde consola */
		if errorteclado := keyboard.Open(); errorteclado != nil {
			fmt.Println("No se logro obtener peticionador desde consola")
			return
		}
		defer func() {
			keyboard.Close()
		}()

		/* Mostramos el peticionador de mensajes */
		fmt.Printf("%s>", usuario)

		for {
			caracter, tecla, errorteclado := keyboard.GetKey() //obtenemos un pulso a la vez
			/* Error del peticionador desde consola */
			if errorteclado != nil {
				fmt.Println("Se obtuvo error del peticionador de consola")
				go_comando <- "Detener"
				wg.Wait()
				break
			}
			/* Si es un caracter el pulso desde teclado, lo agregamos al mensaje */
			if tecla == 0 {
				mensaje = mensaje + string(caracter)
				fmt.Printf("%c",caracter)
			/* Si es un espacio en blanco la tecla pulsada, agregamos espacio al mensaje */
			} else if tecla == keyboard.KeySpace {
				mensaje = mensaje + " "
				fmt.Print(" ")
			/* Si la tecla Esc es pulsada, nos salimos de la aplicacion */
			} else if tecla == keyboard.KeyEsc {
				go_comando <- "Detener"
				wg.Wait()
				break
			/* Si la peticion es borrar caracter, lo borramos del mensaje */
			} else if tecla == keyboard.KeyBackspace || tecla == keyboard.KeyBackspace2 {
				if len(mensaje) != 0 {
					mensaje = mensaje[:len(mensaje)-1]
					fmt.Print("\b \b")
				}
			/* Si la tecla Enter es pulsada y hay mensaje escrito, se envia el mensaje */
			} else if tecla == keyboard.KeyEnter && len(mensaje) > 0 {
				go_comando <- "Pausar" //Pausamos el actualizador de mensajes
				/* Solicitamos el envio de mensaje hasta que se logre enviar */
				for {
					if Enviar(usuario, mensaje) {
						mensaje = ""
						break;
					} else {
						time.Sleep(1 * time.Second)
					}
				}
				go_comando <- "Iniciar" //Retomamos el actualizador de mensajes
			}
		}

	} else {
		fmt.Println("No se logro conectar a:",urlchat)
	}
}



/* Actualizacion de sala de chat en paralelo */
func Routine(go_comando <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	var status = "Iniciar"
	for {
		select {
			case cmd := <-go_comando:
				switch cmd {
					case "Detener": return //cancela la actualizacion de mensajeria
					case "Pausar": status = "Pausar" //pausa las actualizacion de mensajeria
					default: status = "Iniciar" //ejecuta nuevamente las actualizaciones de mensajeria
				}
			break
			default:
				if status == "Iniciar" {
					texto, haynuevo := Mensajeria() //verfica hay mensajes nuevos
					if haynuevo {
						/* borramos la linea en la que esta escribiendo */
						for i := 0; i < (len(usuario) + len(">") + len(mensaje)); i++ {
							fmt.Print("\b \b")
						}
						fmt.Print("\r")
						/* mostrmos los nuevos mensajes */
						fmt.Println(texto)
						/* mostramos nuevamente el peticionador de mensajes con el mensaje que estaba escribiendo */
						fmt.Printf("%s> %s", usuario, mensaje)
					}
					time.Sleep(1000 * time.Millisecond) //actualizamos la mensajeria a cada segundo	
				}
			break
		}
	}
}



/* Valida si la conexion funciona correctamente */
func ConexionValida() bool {
	respuesta, conexionerror := http.Get(urlchat) 
	if conexionerror != nil {
		return false
	}
	defer respuesta.Body.Close()
	if respuesta.StatusCode != 200 {
		return false
	}
	_, conexionerror = ioutil.ReadAll(respuesta.Body)
	if conexionerror != nil {
		return false
	}
	return true //retorna que hay conexion al servidor de chat
}



/* Envia mensajes al chat y valida se haya enviado */
func Enviar(usuario, mensaje string) bool {
	linkparseado, conexionerror := url.Parse(urlchat)
	if conexionerror != nil {
		return false
	}
	parametros := url.Values{}
	parametros.Add("nombre", usuario)
	parametros.Add("mensaje", mensaje)
	linkparseado.RawQuery = parametros.Encode()
	respuesta, conexionerror := http.Get(linkparseado.String()) 
	if conexionerror != nil {
		return false
	}
	defer respuesta.Body.Close()
	if respuesta.StatusCode != 200 {
		return false
	}
	actualizado, conexionerror := ioutil.ReadAll(respuesta.Body)
	if conexionerror != nil {
		return false
	}

	lineas_nuevas := strings.Split(string(actualizado) , "\n")
	lineas_viejas := strings.Split(string(saladechat) , "\n")
	imprimir := false
	exactamente_iguales := false
	i := 0
	j := 0
	retorno := ""
	if len(string(actualizado)) > 0 {
		for i = 0; i < len(lineas_nuevas) /*&& imprimir==false*/; i++ {
			campos_nuevos := strings.Split(lineas_nuevas[i] , ";")

			if imprimir == false {
				/* validamos si ya tenemos los mensajes actualizados */
				for j = 0; j < len(lineas_viejas) && imprimir == false; j++ {
					campos_viejos := strings.Split(lineas_viejas[j] , ";")
					if campos_nuevos[0] == campos_viejos[0] {
						imprimir = true
						if j == 0 {
							/* la mensajeria en la nube y la local son exactamente iguales */
							exactamente_iguales = true 
						}
					} 
				}
				/* si ambas mensajeria son exactamente iguales terminamos de comparar */
				if exactamente_iguales == true {
					break
				/* si son similares, determinamos las diferencias para imprimirlas */
				} else if imprimir == true {
					diferencia := len(lineas_viejas) - j
					i = diferencia + 1
				/* si no hay nada igual, imprimimos toda la mensajeria de la nube */
				} else if j == len(lineas_viejas) {
					imprimir=true
				}				
			}
			
			/* generamos el string a imprimir */
			if imprimir == true {
				campos_nuevos = strings.Split(lineas_nuevas[i] , ";")
				if len(retorno) == 0 {
					retorno = campos_nuevos[1]+": "+campos_nuevos[2]
				} else {
					retorno = retorno + "\n" + campos_nuevos[1] + ": " + campos_nuevos[2]
				} 
			}

		}
	}
	
	/* si la mensajeria local no es exactamente igual a la existente en la nube, entonces imprimiremos los nuevos mensajes */
	if exactamente_iguales == false && imprimir == true {
		/* borramos la linea en la que estaba el mensaje que enviamos */
		for i = 0; i < (len(usuario) + len("> ") + len(mensaje)); i++ {
			fmt.Print("\b \b")
		}
		fmt.Print("\r")
		/* mostramos los mensajes nuevos */
		fmt.Println(retorno)
		/* mostramos nuevamente el peticionador de mensajes */
		fmt.Printf("%s> ", usuario)
		/* actualizamos la sala de chat local */
		saladechat = actualizado
		return true //retornamos que el mensaje se envio exitosamente
	} else {
		return false
	}

}



/* Retorna Mensajes nuevos y si existen mensajes */
func Mensajeria() (string, bool) {
	linkparseado, conexionerror := url.Parse(urlchat)
	if conexionerror != nil {
		return "",false
	}
	parametros := url.Values{}
	parametros.Add("nombre", "")
	parametros.Add("mensaje", "")
	linkparseado.RawQuery = parametros.Encode()
	respuesta, conexionerror := http.Get(linkparseado.String()) 
	if conexionerror != nil {
		return "",false
	}
	defer respuesta.Body.Close()
	if respuesta.StatusCode != 200 {
		return "",false
	}
	actualizado, conexionerror := ioutil.ReadAll(respuesta.Body)
	if conexionerror != nil {
		return "",false
	}

	lineas_nuevas := strings.Split(string(actualizado) , "\n")
	lineas_viejas := strings.Split(string(saladechat) , "\n")
	imprimir := false
	exactamente_iguales := false
	i := 0
	j := 0
	retorno := ""
	if len(string(actualizado)) > 0 {
		for i = 0; i < len(lineas_nuevas) /*&& imprimir==false*/; i++ {
			campos_nuevos := strings.Split(lineas_nuevas[i] , ";")

			if imprimir == false {
				/* validamos si ya tenemos los mensajes actualizados */
				for j = 0; j < len(lineas_viejas) && imprimir == false; j++ {
					campos_viejos := strings.Split(lineas_viejas[j] , ";")
					if campos_nuevos[0] == campos_viejos[0] {
						imprimir = true
						if j == 0 {
							/* la mensajeria en la nube y la local son exactamente iguales */
							exactamente_iguales = true 
						}
					} 
				}
				/* si ambas mensajeria son exactamente iguales terminamos de comparar */
				if exactamente_iguales == true {
					break
				/* si son similares, determinamos las diferencias para imprimirlas */
				} else if imprimir == true {
					diferencia := len(lineas_viejas) - j
					i = diferencia + 1
				/* si no hay nada igual, imprimimos toda la mensajeria de la nube */
				} else if j == len(lineas_viejas) {
					imprimir=true
				}				
			}
			
			/* generamos el string a imprimir */
			if imprimir == true {
				campos_nuevos = strings.Split(lineas_nuevas[i] , ";")
				if len(retorno) == 0 {
					retorno = campos_nuevos[1]+": "+campos_nuevos[2]
				} else {
					retorno = retorno + "\n" + campos_nuevos[1] + ": " + campos_nuevos[2]
				} 
			}

		}
	}
	
	/* si la mensajeria local no es exactamente igual a la existente en la nube, entonces imprimiremos los nuevos mensajes */
	if exactamente_iguales == false && imprimir == true {
		saladechat = actualizado
		return retorno, true //retornamos que los mensajes nuevos y que si existen mensajes nuevos
	} else {
		return "", false //retornamos que no hay mensajes nuevos
	}
}

