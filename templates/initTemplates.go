package initTemplate

import (
	"fmt"
	"html/template"
	"os"
)

var Temp *template.Template

func InitTemplate() {
	// Création d'une nouvelle instance de template avec les fonctions personnalisées
	temp := template.New("").Funcs(template.FuncMap{
		"inc": func(a int) int {
			return a + 1
		},
		"dec": func(a int) int {
			return a - 1
		},
	})

	// Parsing des fichiers de template avec les fonctions personnalisées incluses
	temp, err := temp.ParseGlob("./templates/*.html")

	if err != nil {
		fmt.Printf("ERREUR LORS DE L'OUVERTURE DES TEMPLATES: %v", err.Error())
		os.Exit(1)
	}

	Temp = temp
}
