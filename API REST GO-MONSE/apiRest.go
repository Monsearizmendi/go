package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//La estructura se maneja con Json

type task struct {
	ID        int    `json:ID`
	Nombre    string `json:Nombre`
	Contenido string `json:Contenido`
}

type Contacto struct {
	Nombre, Dirección, CorreoElectronico string
	id                                   int
}
type allTask []task

var tareas = allTask{
	{
		ID:        1,
		Nombre:    "Monserrath Bello Arizmendi",
		Contenido: "Primer prueba en Golang",
	},
}

func main() {
	//Esta linea permite que las rutas sean estrictas
	router := mux.NewRouter().StrictSlash(true)
	//Se definen las rutas del servidor
	//_______________________________________________________
	router.HandleFunc("/", rutaDePrueba)
	router.HandleFunc("/obtenerTareas", obtenerTareas).Methods("GET")
	router.HandleFunc("/obtenerTareas", crearTareas).Methods("POST")
	router.HandleFunc("/obtenerTareas/{id}", obtenerTarea).Methods("GET")
	router.HandleFunc("/obtenerTareas/{id}", borrarTareas).Methods("DELETE")
	router.HandleFunc("/obtenerTareas/{id}", actualizarDatos).Methods("PUT")
	router.
		//_______________________________________________________
		//Se inicia el servidor en el puerto 8000
		log.Fatal(http.ListenAndServe(":8000", router))
	log.Println("Servidor corriendo")

}

// Primer prueba jaja
func rutaDePrueba(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bienvenido al Server")
}

// Esta funcion mostrara las tareas existentes
func obtenerTareas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tareas)
}

// Esta funcion creara tareas y las pondra en un Json
func crearTareas(w http.ResponseWriter, r *http.Request) {
	var nuevaTarea task
	//	Este modulo recibira los datos del cliente al servidor
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Inserta una Tarea valida")
	}
	json.Unmarshal(reqBody, &nuevaTarea)
	//Esta linea hara que el Id se incremente automaticamente
	nuevaTarea.ID = len(tareas) + 1
	//Esta linea de codigo almacenara los datos en un nuevo Json
	tareas = append(tareas, nuevaTarea)
	//Esto permitira que solo se trabaje con Archivos Json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevaTarea)
}

// Esta funcion permitira obtener un dato a traves del id colocando al final del url /id
func obtenerTarea(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//Este metodo convertira un string en un entero
	idTareas, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "ID invalido")
		return
	}

	for _, tareas := range tareas {
		if tareas.ID == idTareas {
			w.Header().Set("Content Type", "application/json")
			json.NewEncoder(w).Encode(tareas)
		}
	}
}

// Esta funcion permitira borrar tareas
func borrarTareas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idTareas, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Id Invalido")
		return
	}

	for i, t := range tareas {
		if t.ID == idTareas {
			tareas = append(tareas[:i], tareas[i+1:]...)
			fmt.Fprintf(w, "La tarea ha sido eliminada con exito")
		}
	}
}

// Esta funcion permitira  actualizar los Datos del Crud
func actualizarDatos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idTareas, err := strconv.Atoi(vars["id"])
	var tareaActualizada task
	if err != nil {
		fmt.Fprintf(w, "Id Invalido")
		return
	}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Ingresa datos validos")
	}
	json.Unmarshal(reqBody, &tareaActualizada)

	for i, t := range tareas {
		if t.ID == idTareas {
			tareas = append(tareas[:i], tareas[i+1:]...)
			tareaActualizada.ID = idTareas
			tareas = append(tareas, tareaActualizada)
			fmt.Fprintf(w, "La tarea con el id %v ha sido actualizada con exito", idTareas)
		}
	}
}

func ObtenerBaseDeDatos() (db *sql.DB, e error) {
	usuario := "root"
	pass := ""
	host := "tcp(127.0.0.1:3306)"
	nombreBaseDeDatos := "monse"
	db, err := sql.open("mysql", fmt.sprintf("%s:%s@%s/%s", usuario, pass, host, nombreBaseDeDatos))
	if err != nil {
		return nil, err
	}
	return db nil
}
