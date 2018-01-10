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
	"github.com/danielsussa/checkBlurFactor/shared"
)

//VISION

type responseVision struct {
	Keywords []keywords
	BorderCtrl borderCtrl
	ImageSpec imageSpec
}

type keywords struct {
	Word string
	Out string
	Score float64
	Position int
}

type borderCtrl struct {
	Up int
	Down int
	Left int
	Right int
}

type imageSpec struct {
	Width int
	Height int
}

//Classify

type classify struct {
	Classification string
	Formations []formation
}

type formation struct {
	Words []string
	Dist []int
}


var query map[string]interface{}

var classifications []classify

func main(){

	//Load Inteligence
	loadAI()

	//Load CSV
	docs := shared.LoadSource()

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

func loadAI(){
	query := make(map[string]interface{},0)

	file, err := ioutil.ReadFile("./step3_test_vision/classify.json")

	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	json.Unmarshal(file,&classifications)

	k := make([]map[string]interface{},0)

	for _,class := range classifications{
		for _, form := range class.Formations {
			for _, word := range form.Words {
				k = append(k, map[string]interface{}{"word":word})
			}
		}
	}
	query["keywords"] = k
}

func run(path string,index int){
	fmt.Println("Running analizes:",index)
	resp := analizeVision(path)

	//Compare AI
	for _,class := range classifications{

	}



	img := downloadImage(path)
	ioutil.WriteFile(fmt.Sprintf("src/blur_test/%d_image.jpg",index), []byte(img), 0777)
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

func analizeVision(path string)(r responseVision){

	resp, err := resty.R().SetHeader("Content-Type", "application/json").
		SetBody(query).
		Post("http://localhost:8085/vision/"+path)

	if err != nil {
		panic(err)
	}

	json.Unmarshal(resp.Body(),&r)
	return
}
