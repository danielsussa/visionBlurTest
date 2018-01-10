package shared

import (
	"encoding/csv"
	"bufio"
	"io"
	"os"
	"log"
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

func LoadSource()(docs []Docs){
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