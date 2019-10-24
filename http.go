package main

import ( 	"fmt"
			"net/http"
			"encoding/json"
			"os"
			"bufio"
		)

const fileURL = "dbLang.txt"

type Language struct{
	Language string 
	Hello string
}

//GET /hello
func getHello(param string, writer http.ResponseWriter){
	/*if param == "en" {
		fmt.Fprintf(writer, "Hello\n")
	} else if param == "fr" {
		fmt.Fprintf(writer, "Bonjour\n")
	} else {
		fmt.Fprintf(writer, "Hello\n")
	}*/
	check := false
	langArray := readFile()
	for i :=0; i < len(langArray); i++{
		if param == langArray[i].Language{
			fmt.Fprintf(writer, langArray[i].Hello)
			check = true
		}
	}
	if check == false{
		fmt.Fprintf(writer, "hello")
	}
}

//Read File 
func readFile() []Language {
	var languageArray []Language
	var langageRead Language
	file, err := os.Open(fileURL)
	if err != nil{
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan(){
		langageRead.Language = scanner.Text()
		scanner.Scan()
		langageRead.Hello = scanner.Text()
		languageArray = append(languageArray, langageRead)
	}
	file.Close()
	return languageArray
}
//Write File
func writeFile(array []Language){
	file, err := os.Create(fileURL)
	if err != nil{
		panic(err)
	}
	for i :=0; i<len(array); i++{
		fmt.Fprintln(file, array[i].Language, array[i].Hello)

	}
	file.Close()
}
//Add new language
func addLanguage(array []Language, new Language) []Language{
	var check = false
	for i :=0; i < len(array); i++{
		if new.Language == array[i].Language{
			check = true
		}
	}
	if check == false{
		array = append(array, new)
	}
	return array
}
//Delete language
func deleteLang(array []Language, param string) []Language{
	var newArray []Language
	for i :=0; i < len(array); i++{
		if array[i].Language != param{
			newArray = append(newArray, array[i])
		}
	}
	return newArray
}

//API fucntion
func API(writer http.ResponseWriter, request *http.Request){
	switch request.Method {
	case "GET":
		getHello(request.URL.Query().Get("lang"), writer)

	case "POST":
		var newLanguage Language
		err := json.NewDecoder(request.Body).Decode(&newLanguage)
		if err != nil{
			panic(err)
		}
		array := readFile()
		array = addLanguage(array, newLanguage)
		writeFile(array)

	case "DELETE":
		array := readFile()
		array = deleteLang(array, request.URL.Query().Get("lang"))
		writeFile(array)

	default:
		getHello(request.URL.Query().Get("lang"), writer)
	}
}


//server
func main(){
	http.HandleFunc("/", API)
	http.ListenAndServe(":8080", nil)
}