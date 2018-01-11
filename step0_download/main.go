package main

import (

	"github.com/go-resty/resty"
	"fmt"
	"io/ioutil"
	"acesso.io/acessorh/datasource/srcdoc"

	"acesso.io/acessorh/lib/uuid"
	"github.com/danielsussa/visionBlurTest/shared"
)

func main(){
	//Load CSV
	docs := shared.LoadSource()

	for i,doc := range docs {
		if i == 0 {
			continue
		}
		uid,_ := uuid.Parse(doc.Doc)
		d := srcdoc.NewEmpty(uid)
		if doc.Foto != ""{
			run(doc.Foto,fmt.Sprintf("%d_f_%s",i,d.Name()))
		}
		if doc.Foto1 != ""{
			run(doc.Foto1,fmt.Sprintf("%d_f1_%s",i,d.Name()))
		}
		if doc.Foto2 != ""{
			run(doc.Foto2,fmt.Sprintf("%d_f2_%s",i,d.Name()))
		}
		if doc.Foto3 != ""{
			run(doc.Foto3,fmt.Sprintf("%d_f3_%s",i,d.Name()))
		}
		if doc.Foto4 != ""{
			run(doc.Foto4,fmt.Sprintf("%d_f4_%s",i,d.Name()))
		}
		if doc.Foto5 != ""{
			run(doc.Foto5,fmt.Sprintf("%d_f5_%s",i,d.Name()))
		}
		if doc.Foto6 != ""{
			run(doc.Foto6,fmt.Sprintf("%d_f6_%s",i,d.Name()))
		}

	}
}

func run(path string,name string){
	img := downloadImage(path)
	ioutil.WriteFile(fmt.Sprintf("src/all_images/%s.jpg",name), []byte(img), 0777)
}

func downloadImage(path string)(r []byte){
	query := map[string]interface{}{
		"scale":map[string]interface{}{"kind":"size","size":1000},
		"cropped":false,
	}

	resp, err := resty.R().SetHeader("Content-Type", "application/json").
		SetBody(query).
		Post("http://localhost:8085/download/"+path)

	if err != nil {
		panic(err)
	}

	return resp.Body()
}