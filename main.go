package main

import (
	"fmt"
	"html"
	"io"
	"math"
	"net/http"
	"os"
)

//Variables Globales
var Kenobi = [2]int{-500, -200}
var Skywalker = [2]int{100, -100}
var Sato = [2]int{500, 100}

func main() {

	//Working Status
	fmt.Printf("Server is Up")

	//Routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World, %q", html.EscapeString(r.URL.Path))
		io.WriteString(w, "\n Is working this site")
	})

	http.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Contact: %q", "Escobar Sebastian")
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

	//Start Servidor
	port := os.Getenv("PORT")

	http.ListenAndServe(":"+port, nil)

}

//Funcion GetLocation: Localizar mediante distancias la ubicacion (Trilateracion)
func GetLocation(distance ...float32) (x, y float32) {

	//Prueba para fines practicos
	P1 := [2]int{0, 0}
	P2 := [2]int{3, 0}
	P3 := [2]int{2, -4}

	//P1 := Kenobi
	//P2 := Skywalker
	//P3 := Sato

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
