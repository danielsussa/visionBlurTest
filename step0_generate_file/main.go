package main

import (
	"fmt"
	"io/ioutil"

	"acesso.io/acessorh/context/ctxcfg"
	"github.com/acesso-io/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var mongo *mgo.Database

type Docs struct {
	Doc   uuid.UUID
	Foto  foto
	Foto1 foto
	Foto2 foto
	Foto3 foto
	Foto4 foto
	Foto5 foto
	Foto6 foto
	Cep   string
}

type foto struct {
	Path string
}

func main() {
	ctxcfg.New("acessorh")

	session, _ := mgo.Dial(ctxcfg.Env.MongoHost)
	mongo = session.DB(ctxcfg.Env.MongoDB)

	docs := getDocs()

	finalLine := fmt.Sprintf("doc,foto,foto1,foto2,foto3,foto4,foto5,foto6,cep\n")
	line := ""

	for k, doc := range docs {
		line += doc.Doc.String() + ","

		i := 0
		for {

			if i == 0 {
				line += doc.Foto.Path + ","
			}

			if i == 1 {
				line += doc.Foto1.Path + ","
			}

			if i == 2 {
				line += doc.Foto2.Path + ","
			}

			if i == 3 {
				line += doc.Foto3.Path + ","
			}

			if i == 4 {
				line += doc.Foto4.Path + ","
			}

			if i == 5 {
				line += doc.Foto5.Path + ","
			}

			if i == 6 {
				line += doc.Foto6.Path + ","
			}

			if i == 6 {
				break
			}
			i++
		}

		//Doc endereço
		//if doc.Doc.String() == "7d5eb9da-b50e-49f7-8c3e-3ee5877eb620" {
		//	line += doc.Cep
		//}
		//line += ","

		line += "\n"

		if k%5000 == 0 {
			//fmt.Println("Saving report on:", "export.csv")
			//ioutil.WriteFile(fmt.Sprintf("export_%d.csv",k), []byte(line), 0777)
			finalLine += line
			line = ""
		}
	}

	fmt.Println("Saving report on:", "import.csv")
	ioutil.WriteFile("import.csv", []byte(finalLine), 0777)

}

func getDocs() (docs []Docs) {

	uid, _ := uuid.Parse("541c9413-2979-419e-b3bc-b13dc189ad60")

	err := mongo.C("srcdoc").Find(bson.M{
		"_deleted": false,
		"doc":      bson.M{"$nin": []bson.Binary{{Kind: 0x04, Data: uid[:]}}},
	}).All(&docs)

	if err != nil {
		panic(err)
	}

	return
}
