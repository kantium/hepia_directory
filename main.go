/*
	Author  : Stéphane Küng
	Comment : Command line Directory for hepia
	Date    : 2 Juin 2017
	Version : 0.0.1
*/

//Package
package main

//Imports
import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"golang.org/x/net/html/charset"
)

//URL of the directory
const url string = "http://hepia.hesge.ch/fr/accueil/annuaire/annuaire-detaille/?listeNom=tous&listeFonction=tous&listePole=tous&listeFiliere=tous&listeInstitut=tous&efRecherche=$1&pbEnvoyer=envoyer"

//Main function
func main() {

	if len(os.Args) < 2 {
		fmt.Println("Hepia Directory 0.0.1 | by Stéphane Küng")
		fmt.Println("Usage : " + path.Base(os.Args[0]) + " name ")
		os.Exit(1)
	}

	resp, err := http.Get(strings.Replace(url, "$1", os.Args[1], 1))
	if err != nil {
		fmt.Printf("ERROR: Failed to crawl url %v\n", err)
		os.Exit(2)
	}
	defer resp.Body.Close()

	var reader io.Reader
	reader, err = charset.NewReader(resp.Body, "")
	if err != nil {
		fmt.Printf("ERROR: Creating new charset reader %v\n", err)
		os.Exit(3)
	}

	root, err := html.Parse(reader)
	if err != nil {
		fmt.Printf("ERROR: parsing HTML %v\n", err)
		os.Exit(3)
	}

	contacts, ok := getContacts(root)
	if !ok {
		log.Fatal("could not find table")
	}

	if len(contacts) > 0 {
		printContacts(contacts)
	} else {
		fmt.Println("No result")
	}

	return
}

//Print a list of contacts
func printContacts(contacts [][]string) {
	for _, element := range contacts {
		fmt.Printf("%-7s %-12s %s\n", element[2][13:], element[0], element[1])
	}
}

//Extract all contact from the hepia web page
func getContacts(node *html.Node) ([][]string, bool) {

	contacts := [][]string{}

	if node.DataAtom == atom.Table {
		node = node.FirstChild

		for c := node.FirstChild; c != nil; c = c.NextSibling {

			contact := []string{}

			for j := c.FirstChild; j != nil; j = j.NextSibling {

				x := j.FirstChild
				if x != nil && x.DataAtom != atom.Br && x.DataAtom != atom.A {
					contact = append(contact, x.Data)
				}
			}

			if len(contact) > 2 {
				contacts = append(contacts, contact)
			}
		}
		return contacts, true
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if cs, ok := getContacts(c); ok {
			return cs, true
		}
	}
	return contacts, false
}
