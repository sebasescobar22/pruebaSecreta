package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"

	"os"
	"strconv"
)

//Variables y Estructuras para facilitar
var Kenobi = [2]int{-500, -200}
var Skywalker = [2]int{100, -100}
var Sato = [2]int{500, 100}

//Info de Formulario Entrante
type Info struct {
	sat1_distancia float64
	sat1_mensaje   string
	sat2_distancia float64
	sat2_mensaje   string
	sat3_distancia float64
	sat3_mensaje   string
}

//Info de Salida "Secreta"
type InfoSec struct {
	PosX        float32
	PosY        float32
	Position    [2]float32
	MsgDsncrptd string
}

//EstructurasPara formatear JSON {
type Satellite struct {
	Name     string
	Distance float64
	Message  []string
}
type PuntoPosicion struct {
	X float64
	Y float64
}

type JsonPost struct {
	Satellites [3]Satellite
}

type JsonAnswer struct {
	Position PuntoPosicion
	Message  string
}

//}

func main() {

	//Working Status
	fmt.Println("Server is Up")

	//Modo de trabajo
	modo := []string{"on", "off"}
	local := modo[1]

	//Routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
		io.WriteString(w, "\n This site is working")
	})

	http.HandleFunc("/localizar", func(w http.ResponseWriter, r *http.Request) {

		// Prueba con distancias
		locationX, locationY := GetLocation(3.0, 3.0, 3.0)
		//xInt, yInt := int(locationX), int(locationY)

		if locationX == 0 && locationY == 0 {
			fmt.Fprintf(w, "No se ha podido Localizar")
		} else {
			fmt.Fprintf(w, "Localizacion aproximada: %.2f , %.2f", locationX, locationY)
		}

	})

	http.HandleFunc("/desencriptado", func(w http.ResponseWriter, r *http.Request) {
		slice1 := []string{"", "Este", "es", ""}
		slice2 := []string{"Este", "", "", "Mensaje"}
		slice3 := []string{"", "es", "un", ""}

		fmt.Fprintf(w, "Mensaje Desencriptado Posible: %+v\n", GetMessage(slice1, slice2, slice3))

	})

	http.HandleFunc("/top-secret/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			t, _ := template.ParseFiles("form.html")

			t.Execute(w, nil)
		}

		if r.Method == "POST" {
			r.ParseForm()

			// Extraigo la informacion del formulario en una estructura
			newInfo := Info{}
			newInfo.sat1_distancia, _ = strconv.ParseFloat(r.Form.Get("sat1_distancia"), 32)
			newInfo.sat1_mensaje = r.Form.Get("sat1_mensaje")
			newInfo.sat2_distancia, _ = strconv.ParseFloat(r.Form.Get("sat2_distancia"), 32)
			newInfo.sat2_mensaje = r.Form.Get("sat2_mensaje")
			newInfo.sat3_distancia, _ = strconv.ParseFloat(r.Form.Get("sat3_distancia"), 32)
			newInfo.sat3_mensaje = r.Form.Get("sat3_mensaje")

			arrayStringMsg1 := separarString(newInfo.sat1_mensaje)
			arrayStringMsg2 := separarString(newInfo.sat2_mensaje)
			arrayStringMsg3 := separarString(newInfo.sat3_mensaje)

			secretX, secretY := GetLocation(float32(newInfo.sat1_distancia), float32(newInfo.sat2_distancia), float32(newInfo.sat3_distancia))

			newDesencriptado := InfoSec{}
			newDesencriptado.PosX = secretX
			newDesencriptado.PosY = secretY
			newDesencriptado.Position[0], newDesencriptado.Position[1] = secretX, secretY

			newDesencriptado.MsgDsncrptd = GetMessage(arrayStringMsg1, arrayStringMsg2, arrayStringMsg3)

			//Muestra resultados
			y, _ := template.ParseFiles("datos.html")
			//check(err)

			y.Execute(w, newDesencriptado)

			//Aplicando Formato JSON y Req Post
			sat1 := Satellite{"Kenobi", newInfo.sat1_distancia, arrayStringMsg1}
			sat2 := Satellite{"Skywalker", newInfo.sat2_distancia, arrayStringMsg2}
			sat3 := Satellite{"Sato", newInfo.sat3_distancia, arrayStringMsg3}

			auxSat := [3]Satellite{sat1, sat2, sat3}
			satelites := JsonPost{auxSat}

			jsonSatelites, _ := json.Marshal(satelites)
			w.Write(jsonSatelites)
			client := &http.Client{}

			var url string
			if local == modo[0] {
				url = "http://localhost:3000/top-secretjson/"
			} else {
				url = "https://immense-bastion-33822.herokuapp.com/top-secretjson/"

			}
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonSatelites))

			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				check(err)
				fmt.Println("Unable to reach the server.")
			} else {
				body, _ := ioutil.ReadAll(resp.Body)
				fmt.Println(string(body))
			}

			// Fin Formato Json y POST Json
		}

	})

	http.HandleFunc("/top-secretjson/", func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var t JsonPost
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}

		secretX, secretY := GetLocation(float32(t.Satellites[0].Distance), float32(t.Satellites[1].Distance), float32(t.Satellites[2].Distance))

		newDesencriptado := InfoSec{}
		newDesencriptado.PosX = secretX
		newDesencriptado.PosY = secretY
		newDesencriptado.Position[0], newDesencriptado.Position[1] = secretX, secretY

		newDesencriptado.MsgDsncrptd = GetMessage(t.Satellites[0].Message, t.Satellites[1].Message, t.Satellites[2].Message)

		newAnswer := JsonAnswer{}
		newAnswer.Position.X = float64(secretX)
		newAnswer.Position.Y = float64(secretY)
		newAnswer.Message = newDesencriptado.MsgDsncrptd

		fmt.Fprintf(w, "Response Code 200 \n")

		jsonData, auxErr := json.Marshal(newAnswer)
		check(auxErr)

		w.Write(jsonData)

	})

	if local == modo[0] {
		//Prueba local
		http.ListenAndServe(":3000", nil)

	} else {
		//Servidor Heroku
		port := os.Getenv("PORT")

		http.ListenAndServe(":"+port, nil)
	}

}

//Funcion GetLocation: Localizar mediante distancias la ubicacion (Trilateracion)
func GetLocation(distance ...float32) (x, y float32) {

	//Prueba para fines practicos
	//P1 := [2]int{0, 0}
	//P2 := [2]int{3, 0}
	//P3 := [2]int{2, -4}

	P1 := Kenobi
	P2 := Skywalker
	P3 := Sato

	if len(distance) == 3 {

		//Hay valores que no se llegan a utilizar (Diferentes Formulas)
		//cambia taza de error al aproximar
		ax := P1[0]
		ay := P1[1]
		ar := float64(distance[0])

		bx := P2[0]
		//by:=P2[1]
		br := float64(distance[1])

		cx := P3[0]
		cy := P3[1]
		//cr:=float64(distance[2])

		d := bx - ax
		i := cx - ax
		j := cy - ay

		x := (math.Pow(ar, 2) - math.Pow(br, 2) + math.Pow(float64(d), 2)) / (2 * float64(d))
		y := ((math.Pow(ar, 2) - math.Pow(br, 2) + math.Pow(float64(i), 2) + math.Pow(float64(j), 2)) / (2 * float64(j))) - ((float64(i) / float64(j)) * x)

		//Liberia Math opera en f64 -> Conversion a f32 para formato de retorno
		//println(x, y)

		xF32 := float32(x)
		yF32 := float32(y)
		//println(xF32, yF32)

		return xF32, yF32
	} else {
		return 0, 0
	}

}

//Funcion GetMessage: Recibe arrays de String, descarta los string vacios,
//compara y anexa los strings para que no se repitan, devolviendo una cadena con los valores unicos.
func GetMessage(messages ...[]string) (msg string) {

	//
	var cadenaMsg [50]string
	var msj string

	for indexM, message := range messages {
		diferenciaPos := 0
		encontrados := 0
		//println("\n Cadena: ", indexM)
		if indexM != 0 {
			for pos, valorMsg := range message {

				if valorMsg != "" {

					for i, valorYaAnexado := range cadenaMsg {

						if valorYaAnexado == valorMsg {
							posEscontrada := i
							encontrados += 1

							if posEscontrada != pos {
								diferenciaPos = posEscontrada - pos
							}

							//println("\n Encontrados iguales: ", encontrados)
							//println("\n Diferencia Posiciones: ", diferenciaPos)

						}
					}

				}

			}
		}

		for pos, valorMsg := range message {

			found := false
			posiblePos := pos

			if valorMsg != "" && valorMsg != "." {

				for _, valorYaAnexado := range cadenaMsg {
					if valorYaAnexado == valorMsg {
						found = true

						if diferenciaPos == 0 {
							cadenaMsg[posiblePos] = valorMsg
						} else {
							cadenaMsg[diferenciaPos+posiblePos] = valorMsg
						}

					}
				}

				if found == false {
					if indexM == 0 {
						cadenaMsg[posiblePos] = valorMsg

					} else {
						if diferenciaPos != 0 {
							cadenaMsg[diferenciaPos+posiblePos] = valorMsg

						} else {

							//Caso que no haya encontrado palabras similares y ya exista un valor en ese orden
							if encontrados == 0 && cadenaMsg[posiblePos] != "" {
								// Posible Solucion hacer un recursivo en el futuro
								// (Volver a recorrer esta cadena por si en el futuro ya hay similutedes)
								// Estado no implementado

							} else {

								cadenaMsg[posiblePos] = valorMsg
							}

						}
					}

				}
			}
		}

	}

	for _, decod := range cadenaMsg {
		msj = msj + " " + decod
	}
	//println(msj)
	return msj
}

// FUNCION PARA REG DE ERRORES EN EL LOG
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//Funcion para Separar texto en palabras String en []String (Sacado de un foro)
func split(tosplit string, sep rune) []string {
	var fields []string
	auxSep := 0

	last := 0
	for i, c := range tosplit {
		if c == sep {
			auxSep += 1
			// Found the separator, append a slice
			fields = append(fields, string(tosplit[last:i]))
			last = i + 1

		}
	}

	// Don't forget the last field
	fields = append(fields, string(tosplit[last:]))

	return fields
}

func separarString(str string) []string {
	var secuencia []string
	for _, field := range split(str, ' ') {
		secuencia = append(secuencia, field)

	}
	//println(secuencia)
	return secuencia
}
