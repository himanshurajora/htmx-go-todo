package main

import (
	"html/template"
	"log"
	"net/http"
)

type Task struct {
	Name        string
	IsCompleted bool
}

func main() {
	tasks := map[string][]Task{
		"Tasks": {},
	}

	rootPathHandler := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, tasks)
	}

	addTaskHandler := func(w http.ResponseWriter, r *http.Request) {
		name := r.PostFormValue("task_name")
		tmpl := template.Must(template.ParseFiles("index.html"))
		tasks["Tasks"] = append(tasks["Tasks"], Task{Name: name, IsCompleted: false})
		tmpl.ExecuteTemplate(w, "task-list-item", Task{Name: name, IsCompleted: false})
	}
	http.HandleFunc("/", rootPathHandler)
	http.HandleFunc("/add", addTaskHandler)

	log.Fatal(http.ListenAndServe(":5000", nil))
}
