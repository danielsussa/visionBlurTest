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
	fmt.Println("Running analizes:", index)

	docKind := "null"

	query := map[string]interface{}{
		"keywords": map[string]interface{}{
			"nome": map[string]interface{}{},
			"mae":  map[string]interface{}{},
		},
		"phrases": map[string]interface{}{
			"nome_mae": map[string]interface{}{"text": "nome [SKP] mae", "threshold": 0.8},
		},
	}

	vis := analizeVision(pathFrente, query)

	//Check for words
	for _, word := range vis.Keywords {
		if word.Pass {
			docKind = "rg"
			break
		}
	}

	//Check for phrases
	for _, phrase := range vis.Phrases {
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
		docQual = "dist"
	}

	img := downloadImage(pathFrente)
	ioutil.WriteFile(fmt.Sprintf("src/rg_test/%d_%s_%s.jpg", index, docKind, docQual), []byte(img), 0777)
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
