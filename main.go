package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

var id = 0

type Task struct {
	Name string
	Id   int
}

var tasks = map[string][]Task{
	"Tasks": {},
}

// function to write tasks in a file
func writeTasksToFile() {
	file, err := os.Create("tasks.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, task := range tasks["Tasks"] {
		file.WriteString(fmt.Sprintf("%s\n", task.Name))
	}
}

// function to read tasks from a file
func readTasksFromFile() {
	// empty the tasks map
	tasks["Tasks"] = []Task{}

	if _, err := os.Stat("tasks.txt"); os.IsNotExist(err) {
		return
	}

	file, err := os.Open("tasks.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var task Task
	var i = 0
	for {
		fmt.Printf("Reading task %d\n", i)
		_, err := fmt.Fscanf(file, "%s\n", &task.Name)
		task.Id = i
		i++
		if err != nil {
			break
		}
		tasks["Tasks"] = append(tasks["Tasks"], task)
		// log the task
		fmt.Printf("Task: %s\n", task.Name)
	}
}

func main() {
	fmt.Println("Server is running on port 5000")

	rootPathHandler := func(w http.ResponseWriter, r *http.Request) {
		readTasksFromFile()
		fmt.Println("Root path handler")
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, tasks)
	}

	addTaskHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Add task handler")
		name := r.PostFormValue("task_name")
		tmpl := template.Must(template.ParseFiles("index.html"))
		id := id + 1
		tasks["Tasks"] = append(tasks["Tasks"], Task{Name: name, Id: id})
		tmpl.ExecuteTemplate(w, "task-list-item", Task{Name: name, Id: id})
		writeTasksToFile()
	}

	deleteTaskHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Delete task handler")
		id, err := strconv.Atoi(r.URL.Query().Get("id"))

		if err != nil {
			fmt.Println("Error parsing id")
			panic(err)
		}

		for i, task := range tasks["Tasks"] {
			if task.Id == id {
				tasks["Tasks"] = append(tasks["Tasks"][:i], tasks["Tasks"][i+1:]...)
				break
			}
		}

		writeTasksToFile()
	}

	http.HandleFunc("/delete", deleteTaskHandler)
	http.HandleFunc("/add", addTaskHandler)
	http.HandleFunc("/", rootPathHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
