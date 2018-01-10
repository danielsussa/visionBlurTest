package main

import (
	"os"
	"log"
	"encoding/csv"
	"bufio"
	"io"
	"fmt"
	"github.com/go-resty/resty"
	"encoding/json"
	"io/ioutil"
	"strings"
	"strconv"
)

type Docs struct {
	Doc string
	Foto string
	Foto1 string
	Foto2 string
	Foto3 string
	Foto4 string
	Foto5 string
	Foto6 string
}


func main(){
	//Load CSV
	docs := loadSource()

	selectedIndex := make(map[int]bool)

	for _,sel := range strings.Split(os.Getenv("selected"),","){
		if sel == ""{
			continue
		}
		i,_ := strconv.Atoi(sel)
		selectedIndex[i] = true
	}


	selectedDoc := os.Getenv("doc")

	for i,doc := range docs {

		if len(selectedIndex) > 0 && selectedIndex[i] == false {
			continue
		}

		if selectedDoc != "" && doc.Doc != selectedDoc {
			continue
		}

		if i == 0 {
			continue
		}
		if doc.Foto != ""{
			run(doc.Foto,i)
		}
		if doc.Foto1 != ""{
			run(doc.Foto1,i)
		}
		if doc.Foto2 != ""{
			run(doc.Foto2,i)
		}
		if doc.Foto3 != ""{
			run(doc.Foto3,i)
		}
		if doc.Foto4 != ""{
			run(doc.Foto4,i)
		}
		if doc.Foto5 != ""{
			run(doc.Foto5,i)
		}
		if doc.Foto6 != ""{
			run(doc.Foto6,i)
		}
	}
}

func run(path string,index int){
	fmt.Println("Running analizes:",index)
	resp := analizeBlur(path)

	if resp.Result == 0 {
		return
	}

	//folder := ""
	//
	//if resp.Result < 300 {
	//	folder = "low"
	//}
	//
	//if resp.Result >= 300 {
	//	folder = "high"
	//}

	img := downloadImage(path)
	ioutil.WriteFile(fmt.Sprintf("src/blur_test/%d_image_%d.jpg",resp.Result,index), []byte(img), 0777)
}


type response struct {
	Result int
}

func downloadImage(path string)(r []byte){
	query := map[string]interface{}{
	}

	resp, err := resty.R().SetHeader("Content-Type", "application/json").
		SetBody(query).
		Post("http://localhost:8085/download/"+path)

	if err != nil {
		panic(err)
	}

	return resp.Body()
}

func analizeBlur(path string)(r response){

	query := map[string]interface{}{
		"scale":map[string]interface{}{"kind":"size","size":2000},
		"cropped":true,
	}

	resp, err := resty.R().SetHeader("Content-Type", "application/json").
		SetBody(query).
		Post("http://localhost:8085/blur/"+path)

	if err != nil {
		panic(err)
	}

	json.Unmarshal(resp.Body(),&r)
	return
}

func loadSource()(docs []Docs){
	csvFile, err := os.Open("./src/import/export_final.csv")

	if err != nil {
		log.Fatal("Error opening file:", err)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		doc := Docs{
			Doc:line[0],
			Foto:line[1],
			Foto1:line[2],
			Foto2:line[3],
			Foto3:line[4],
			Foto4:line[5],
			Foto5:line[6],
			Foto6:line[7],

		}
		docs = append(docs, doc)
	}
	return
}