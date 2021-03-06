package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/resumic/schema/docs/schema"
	"github.com/tidwall/gjson"
)

func main() {
	var conf, doc, tmpl string
	flag.StringVar(&conf, "conf", "../../schema.json", "Schema file")
	flag.StringVar(&tmpl, "tmpl", "../../docs/dev/doc.tmpl", "Template file")
	flag.StringVar(&doc, "doc", "../../docs/dev/doc.md", "Output file")
	content, err := ioutil.ReadFile(conf)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(tmpl, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	value := gjson.GetBytes(content, "properties")
	value.ForEach(func(key, value gjson.Result) bool {
		_, err = f.WriteString("\n# " + strings.Title(key.String()) + "\nThis section is of " + "{{.Properties." + strings.Title(key.String()) + ".Type}} type that tells about the basic information of the user and consists of the following sub-sections:\n") //print name of sections
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(key.String())
		result := gjson.GetBytes(content, "properties."+key.String())
		result.ForEach(func(key1, value gjson.Result) bool {
			_, err = f.WriteString("##" + spaceGen(3) + key1.String() + "\n") //print name of sub-sections
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(spaceGen(2))
			fmt.Println(key1.String())

			if key1.String() == "properties" {
				result1 := gjson.GetBytes(content, "properties."+key.String()+"."+key1.String())
				result1.ForEach(func(key2, value gjson.Result) bool {
					_, err := f.WriteString("###" + spaceGen(8) + key2.String() + "\nSub-section of type {{.Properties." + strings.Title(key.String()) + ".Properties." + strings.Title(key2.String()) + ".Type}}, used to specify the " + key2.String() + " of the person\nThe schema snippet can be shown below:\n\n       " + key2.String() + ":\n        " + value.String() + "\n\n") //print properties of subsections
					if err != nil {
						log.Fatal(err)
					}
					fmt.Print(spaceGen(9))
					fmt.Println(key2.String())
					return true
				})
			}
			if key1.String() == "items" {
				result1 := gjson.GetBytes(content, "properties."+key.String()+".items.properties")
				result1.ForEach(func(key2, value gjson.Result) bool {
					_, err := f.WriteString("###" + spaceGen(8) + key2.String() + "\nSub-section of type {{.Properties." + strings.Title(key.String()) + ".Items.Properties." + strings.Title(key2.String()) + ".Type}}, used to specify the " + key2.String() + " of the person\nThe schema snippet can be shown below:\n\n       " + key2.String() + ":\n        " + value.String() + "\n\n") //print items of subsections
					if err != nil {
						log.Fatal(err)
					}
					fmt.Print(spaceGen(9))
					fmt.Println(key2.String())
					return true
				})
			}
			return true
		})
		return true
	})
	f.Close()
	var schema schema.Data
	if err := json.Unmarshal(content, &schema); err != nil {
		log.Fatal(err)
	}
	tpl, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create(doc)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	if err = tpl.Execute(file, schema); err != nil {
		log.Fatal(err)
	}
	err = os.Remove(tmpl)
	if err != nil {
		log.Fatal(err)
		return
	}
}
func spaceGen(x int) string {
	var s string
	for i := 0; i < x; i++ {
		s = s + " "
	}
	return s
}
