package main

import (
	"fmt"
)

var m map[string]string

func InsertToDo() {
	var heading, content string
	fmt.Print("Enter the heading for the todo: ")
	fmt.Scan(&heading)
	fmt.Print("Enter the content: ")
	fmt.Scan(&content)
	if heading == "" || content == "" {
		fmt.Println("The heading and content should not be empty")
	}
	m[heading] = content
	fmt.Println("Todo added successfully")
}

func DeleteToDo() {
	var heading string
	fmt.Print("Enter the heading to be deleted in the todo: ")
	fmt.Scan(&heading)
	_, found := m[heading]
	if found {
		delete(m, heading)
		fmt.Println("Todo deleted successfully")
	}
	fmt.Println("Todo not found")
}

func ReadToDos() {
	for heading, content := range m {
		fmt.Printf("\n%s: %s\n", heading, content)
	}
}

func UpdateTodo() {
	var heading, content string
	fmt.Print("Enter the heading to be updated in the todo: ")
	fmt.Scan(&heading)
	_, found := m[heading]
	if found{
		fmt.Print("Enter the content: ")
		fmt.Scan(&content)
		m[heading] = content
		fmt.Println("Todo updated successfully")
	} else{
		fmt.Println("Todo not found")
	}
}

func main() {
	m = make(map[string]string)
	var option int
	for {
		fmt.Print("\n1. Insert Todo\n2. Delete Todo\n3. Read Todos\n4. Update Todo\nEnter your option: ")
		fmt.Scan(&option)
		switch option {
		case 1:
			InsertToDo()
		case 2:
			DeleteToDo()
		case 3:
			ReadToDos()
		case 4:
			UpdateTodo()
		default:
			fmt.Println("Invalid option. Try again")
		}
	}
}
