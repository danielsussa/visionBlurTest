package main

import (
	"os"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/go-resty/resty"

	"github.com/danielsussa/visionBlurTest/shared"
)

//VISION

type responseVision struct {
	Keywords   map[string]keywords
	Phrases    map[string]phrasesKind
	BorderCtrl borderCtrl
	ImageSpec  imageSpec
}

type keywords struct {
	Word     string
	Out      string
	Score    float64
	Pass     bool
	Position int
}

type borderCtrl struct {
	Up    int
	Down  int
	Left  int
	Right int
}

type imageSpec struct {
	Width  int
	Height int
}

type phrasesKind struct {
	Text      string
	Threshold float64
	Pass      bool
}

func main() {

	//Load CSV
	docs := shared.LoadSource()

	selectedIndex := make(map[int]bool)

	for _, sel := range strings.Split(os.Getenv("selected"), ",") {
		if sel == "" {
			continue
		}
		i, _ := strconv.Atoi(sel)
		selectedIndex[i] = true
	}

	selectedDoc := os.Getenv("doc")

	for i, doc := range docs {

		if len(selectedIndex) > 0 && selectedIndex[i] == false {
			continue
		}

		if selectedDoc != "" && doc.Doc != selectedDoc {
			continue
		}

		if i == 0 {
			continue
		}

		//Run RG Validation
		if doc.Doc == "c2069ecf-ea5d-4029-9960-6f802392c6d7" {
			runRG(doc.Foto1, doc.Foto2, i)
		}

		//if doc.Foto != ""{
		//	run(doc.Foto,i)
		//}
		//if doc.Foto1 != ""{
		//	run(doc.Foto1,i)
		//}
		//if doc.Foto2 != ""{
		//	run(doc.Foto2,i)
		//}
		//if doc.Foto3 != ""{
		//	run(doc.Foto3,i)
		//}
		//if doc.Foto4 != ""{
		//	run(doc.Foto4,i)
		//}
		//if doc.Foto5 != ""{
		//	run(doc.Foto5,i)
		//}
		//if doc.Foto6 != ""{
		//	run(doc.Foto6,i)
		//}
	}
}

func runRG(pathFrente string, pathVerso string, index int) {

	queryVerso := map[string]interface{}{
		"keywords": map[string]interface{}{
			"data":       map[string]interface{}{"threshold": 0.9},
			"nascimento": map[string]interface{}{"threshold": 0.9, "nextTo": "data"},

			"valida":     map[string]interface{}{"threshold": 0.9},
			"todo":       map[string]interface{}{"threshold": 0.9, "nextTo": "valida"},
			"territorio": map[string]interface{}{"threshold": 0.9, "nextTo": "todo"},
		},
		"phrases": map[string]interface{}{
			"data_nasc":   map[string]interface{}{"text": "data [SKP] nascimento", "threshold": 0.8},
			"valida_todo": map[string]interface{}{"text": "valida [SKP] todo territorio", "threshold": 0.8},
		},
	}

	queryFrente := map[string]interface{}{
		"keywords": map[string]interface{}{
			"republica":  map[string]interface{}{"threshold": 0.9},
			"federativa": map[string]interface{}{"threshold": 0.9, "nextTo": "republica"},
			"brasil":     map[string]interface{}{"threshold": 0.9, "nextTo": "federativa"},

			"departamento": map[string]interface{}{"threshold": 0.9},
			"nacional":     map[string]interface{}{"threshold": 0.9, "nextTo": "departamento"},
			"transito":     map[string]interface{}{"threshold": 0.9, "nextTo": "nacional"},
		},
		"phrases": map[string]interface{}{
			"republica": map[string]interface{}{"text": "republica federativa [SKP] brasil", "threshold": 0.8},
			"transito":  map[string]interface{}{"text": "departamento nacional [SKP] transito", "threshold": 0.8},
		},
	}
	runAiRG(pathFrente, queryFrente, fmt.Sprintf("%d_frente", index))
	runAiRG(pathVerso, queryVerso, fmt.Sprintf("%d_verso", index))

}

func runAiRG(path string, query map[string]interface{}, id string) {
	vis := analizeVision(path, query)

	docKind := "null"

	//Check for words
	for key, word := range vis.Keywords {
		if word.Pass == true && key == "transito" {
			docKind = "wrong_cpf"
			break
		}
		if word.Pass {
			docKind = "rg"
			break
		}
	}

	//Check for phrases
	for key, phrase := range vis.Phrases {
		if phrase.Pass == true && key == "transito" {
			docKind = "WRONG_CPF"
			break
		}
		if phrase.Pass {
			docKind = "RG"
			break
		}
	}

	// Check Spec of doc to set quality
	ver, hor := checkBorder(vis.BorderCtrl, vis.ImageSpec)

	docQual := "ok"

	if ver < 0.3 || hor < 0.3 {
		docQual = "dist"
	}

	if ver > 0.8 && hor > 0.8 {
		docQual = "OK"
	}

	if ver > 0.98 || hor > 0.98 {
		docQual = "crop"
	}

	img := downloadImage(path)
	ioutil.WriteFile(fmt.Sprintf("src/rg_test/%s_%s_%s.jpg", id, docKind, docQual), []byte(img), 0777)
}

func checkBorder(border borderCtrl, spec imageSpec) (ver float64, hor float64) {
	verRatio := float64(border.Right-border.Left) / float64(spec.Width)
	horRatio := float64(border.Down-border.Up) / float64(spec.Height)

	return verRatio, horRatio

}

func downloadImage(path string) (r []byte) {
	query := map[string]interface{}{}

	resp, err := resty.R().SetHeader("Content-Type", "application/json").
		SetBody(query).
		Post("http://localhost:8085/download/" + path)

	if err != nil {
		panic(err)
	}

	return resp.Body()
}

func analizeVision(path string, query map[string]interface{}) (r responseVision) {

	resp, err := resty.R().SetHeader("Content-Type", "application/json").
		SetBody(query).
		Post("http://localhost:8085/vision/" + path)

	if err != nil {
		panic(err)
	}

	json.Unmarshal(resp.Body(), &r)
	return
}
