package main

import(
	"fmt"
	"net/http"
	"log"
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/gorilla/mux"
)

//Definición de la tareas
type task struct {
	ID int `json:ID`
	Name string `json:Name`
	Content string `json:Content`
}

type allTasks []task 

var tasks = allTasks {
	{
		ID: 1,
		Name: "Task One",
		Content: "Some Content",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func getTaskById(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for _, task := range tasks {
		if task.ID == taskId {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["id"])
	
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for i, task := range tasks {
		if task.ID == taskId {
			tasks = append(tasks[:i], tasks[i + 1:]...) 
			//Elimina el elemento del indice que coincida con el ID y conservando los anteriores y posteriores al indice
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "The task with ID %v has been removed successfully", taskId)
		}
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a valid task")
	}

	json.Unmarshal(reqBody, &newTask)

	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func updateTask(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["id"])
	var updateTask task

	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please enter valid data")
		return
	}
	json.Unmarshal(reqBody, &updateTask)

	for i, task := range tasks {
		if task.ID == taskId {
			// tasks = append(tasks[:i], tasks[i + 1:]...) 
			// updateTask.ID = taskId
			// tasks = append(tasks, updateTask)
			tasks[i].Name = updateTask.Name
			tasks[i].Content = updateTask.Content

			fmt.Fprintf(w, "The task with ID %v has been updated successfully", taskId)
		}
	}
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my API")
}

func main()  {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/task/{id}", getTaskById).Methods("GET")
	router.HandleFunc("/task/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/task/{id}", updateTask).Methods("PUT")
	log.Fatal(http.ListenAndServe(":3000", router))
}