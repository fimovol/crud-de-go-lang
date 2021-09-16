package main

import (
	"database/sql"
	//"log"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexion() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contrasenia := ""
	Nombre := "Sistema"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+
		"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	fmt.Println("servidor corriendo...")
	http.ListenAndServe(":8080", nil)
}

type Empleados struct {
	Id     int
	Nombre string
	Correo string
}

func Inicio(w http.ResponseWriter, r *http.Request) {

	conexionestablesida := conexion()
	registros, err := conexionestablesida.
		Query("SELECT * FROM `empleados`;")

	if err != nil {
		panic(err.Error())
	}
	empleado := Empleados{}
	arregloEmpleado := []Empleados{}

	for registros.Next() {
		var id int
		var nombre, correo string
		err = registros.Scan(&id, &nombre, &correo)

		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		arregloEmpleado = append(arregloEmpleado, empleado)
	}
	//fmt.Fprintf(w, "hola javier")
	plantillas.ExecuteTemplate(w, "inicio", arregloEmpleado)
}
func Crear(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "hola javier")
	plantillas.ExecuteTemplate(w, "crear", nil)
}
func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionestablesida := conexion()
		insertarregistros, err := conexionestablesida.
			Prepare("INSERT INTO `empleados`(`id`, `nombre`, `correo`) VALUES (NULL, ?, ?);")
		if err != nil {
			panic(err.Error())
		}
		insertarregistros.Exec(nombre, correo)

		http.Redirect(w, r, "/", 301)
	}
}
func Borrar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")

	conexionestablesida := conexion()
	insertarregistros, err := conexionestablesida.
		Prepare("DELETE FROM `empleados` WHERE `empleados`.`id` = ?")
	if err != nil {
		panic(err.Error())
	}
	insertarregistros.Exec(idEmpleado)

	http.Redirect(w, r, "/", 301)
}
func Editar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")

	conexionestablesida := conexion()
	registro, err := conexionestablesida.
		Query("SELECT * FROM `empleados` where id=?;", idEmpleado)

	empleado := Empleados{}
	for registro.Next() {
		var id int
		var nombre, correo string
		err = registro.Scan(&id, &nombre, &correo)

		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
	}
	plantillas.ExecuteTemplate(w, "editar", empleado)
}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionestablesida := conexion()
		modificarregistros, err := conexionestablesida.
			Prepare("UPDATE `empleados` SET `nombre` = ?, `correo` = ? WHERE `empleados`.`id` = ?;")
		if err != nil {
			panic(err.Error())
		}
		modificarregistros.Exec(nombre, correo, id)

		http.Redirect(w, r, "/", 301)
	}
}
