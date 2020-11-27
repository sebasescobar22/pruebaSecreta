package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
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
	name     string
	distance float64
	message  []string
}
type PuntoPosicion struct {
	x float64
	y float64
}

type JsonPost struct {
	satellites [3]Satellite
}

type JsonAnswer struct {
	position PuntoPosicion
	message  string
}

//}

func main() {

	//Working Status
	fmt.Printf("Server is Up")

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
			t, err := template.ParseFiles("form.html")
			check(err)
			t.Execute(w, nil)
		}

		if r.Method == "POST" {
			r.ParseForm()

			var auxErr error
			// Extraigo la informacion del formulario en una estructura
			newInfo := Info{}
			newInfo.sat1_distancia, auxErr = strconv.ParseFloat(r.Form.Get("sat1_distancia"), 32)
			newInfo.sat1_mensaje = r.Form.Get("sat1_mensaje")
			newInfo.sat2_distancia, auxErr = strconv.ParseFloat(r.Form.Get("sat2_distancia"), 32)
			newInfo.sat2_mensaje = r.Form.Get("sat2_mensaje")
			newInfo.sat3_distancia, auxErr = strconv.ParseFloat(r.Form.Get("sat3_distancia"), 32)
			newInfo.sat3_mensaje = r.Form.Get("sat3_mensaje")
			check(auxErr)

			arrayStringMsg1 := separarString(newInfo.sat1_mensaje)
			arrayStringMsg2 := separarString(newInfo.sat2_mensaje)
			arrayStringMsg3 := separarString(newInfo.sat3_mensaje)

			//Prueba con formato JSON
			sat1 := Satellite{"Kenobi", newInfo.sat1_distancia, arrayStringMsg1}
			sat2 := Satellite{"Skywalker", newInfo.sat2_distancia, arrayStringMsg2}
			sat3 := Satellite{"Sato", newInfo.sat3_distancia, arrayStringMsg3}

			auxSat := [3]Satellite{sat1, sat2, sat3}
			satelites := JsonPost{auxSat}

			jsonSatelites, auxErr := json.Marshal(satelites)

			_, err := http.Post("/top-secretsjson/", "application/json", bytes.NewBuffer(jsonSatelites))

			// Fin Prueba Json (No se llego a terminar)

			//Continuacion de la Solucion

			secretX, secretY := GetLocation(float32(newInfo.sat1_distancia), float32(newInfo.sat2_distancia), float32(newInfo.sat3_distancia))

			newDesencriptado := InfoSec{}
			newDesencriptado.PosX = secretX
			newDesencriptado.PosY = secretY
			newDesencriptado.Position[0], newDesencriptado.Position[1] = secretX, secretY

			newDesencriptado.MsgDsncrptd = GetMessage(arrayStringMsg1, arrayStringMsg2, arrayStringMsg3)

			//Muestra resultados
			y, err := template.ParseFiles("datos.html")
			check(err)

			y.Execute(w, newDesencriptado)
		}

	})

	http.HandleFunc("/top-secretjson/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			decoder := json.NewDecoder(r.Body)
			var t JsonPost
			err := decoder.Decode(&t)
			if err != nil {
				panic(err)
			}
			log.Println(t.satellites)
			fmt.Printf("Json satellites working")

		}
	})

	//Servidor Heroku
	port := os.Getenv("PORT")

	http.ListenAndServe(":"+port, nil)

	//Prueba local
	//http.ListenAndServe(":3000", nil)

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
		println(x, y)

		xF32 := float32(x)
		yF32 := float32(y)
		println(xF32, yF32)

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
		println("\n Cadena: ", indexM)
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

							println("\n Encontrados iguales: ", encontrados)
							println("\n Diferencia Posiciones: ", diferenciaPos)

						}
					}

				}

			}
		}

		for pos, valorMsg := range message {

			found := false
			posiblePos := pos

			if valorMsg != "" {

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
	println(msj)
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

			if auxSep == 2 {
				// Found 2 times the separator (' '), append a slice
				fields = append(fields, string(tosplit[last:i]))
				last = i + 1

				auxSep = 0
			}
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
	println(secuencia)
	return secuencia
}
