package main

import (
	"gopkg.in/mgo.v2"
	"acesso.io/acessorh/context/ctxcfg"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"io/ioutil"
	"github.com/acesso-io/uuid"
)

var mongo *mgo.Database


func main(){
	ctxcfg.New("acessorh")

	session, _ := mgo.Dial(ctxcfg.Env.MongoHost)
	mongo = session.DB(ctxcfg.Env.MongoDB)

	docs := getDocs()

	fotoArray := []string{"foto","foto1","foto2","foto3","foto4","foto5","foto6"}

	line := fmt.Sprintf("doc,foto,foto1,foto2,foto3,foto4,foto5,foto6\n")

	for _, doc := range docs {
		for _,foto := range fotoArray {
			path := ""
			if doc[foto] != nil {
				path = doc[foto].(map[string]interface{})["path"].(string)
			}
			line += path + ","
		}
		line += "\n"
	}

	fmt.Println("Saving report on:", "export.csv")
	ioutil.WriteFile("export.csv", []byte(line), 0777)

}

func getDocs()(docs []map[string]interface{}){
	mongo.C("srcdoc").Find(bson.M{}).All(&docs)
	return
}
