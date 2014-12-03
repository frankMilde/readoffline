//
// =========================================================================
//
//       Filename:  main.go
//
//    Description: The aim is to download websites containing text articles
//                 and transform them into tow column pdfs.
//
//        Version:  1.0
//        Created:  12/01/2014 09:30:32 PM
//       Revision:  none
//       Compiler:  g++
//
//          Usage:  <+USAGE+>
//
//         Output:  <+OUTPUT+>
//
//         Author:  Frank Milde (fm), frank.milde (at) posteo.de
//        Company:
//
// =========================================================================
//

package main

//--------------------------------------------------------------------------
//  Includes
//--------------------------------------------------------------------------
import (
	//	"bytes"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	//	"strconv"
	"strings"
	"time"
)

//--------------------------------------------------------------------------
//  Globals
//--------------------------------------------------------------------------
var URL string
var FILE string

type websiteSource struct {
	fileType string
	content  string
}

// ===  FUNCTION  ==========================================================
//         Name:  main
//  Description:
// =========================================================================
func main() {
	ClearTerminalScreen()
	start := time.Now()

	flag.Parse()

	source := GetSourceFrom(URL)
	SaveToFile(source.content)
	//	PrintParseTree(source)
	GetTitle(source)

	contentParseTree, _ := html.Parse(strings.NewReader(source.content))

	//	var titleAttr = []html.Attribute{html.Attribute{Key: "class", Val: "article-title"}}
	//	var subtitleAttr = []html.Attribute{html.Attribute{Key: "class", Val: "article-subtitle"}}
	//	title1 := ExtractAttribute(contentParseTree, titleAttr)
	//	fmt.Println("Title1: ", title1)

	var titleNode = html.Node{
		Data: "h1",
		Attr: []html.Attribute{
			html.Attribute{
				Key: "class",
				Val: "article-title",
			},
		},
	}
	var subtitleNode = html.Node{
		Data: "h2",
		Attr: []html.Attribute{
			html.Attribute{
				Key: "class",
				Val: "article-subtitle",
			},
		},
	}
	title := ExtractItem(contentParseTree, titleNode)
	subtitle := ExtractItem(contentParseTree, subtitleNode)

	fmt.Println("Title   : ", title)
	fmt.Println("Subtitle: ", subtitle)

	end := time.Now()
	fmt.Println("\n\nrun time:", end.Sub(start))
}

// ===  FUNCTION  ==========================================================
//         Name:  init
//  Description:  Needed by the flag package to define the variables that
//                are parsed from the command line.
// =========================================================================
func init() {
	const (
		default_URL  = "science.ksc.nasa.gov/shuttle/missions/51-l/docs/rogers-commission/Appendix-F.txt"
		default_FILE = "output.html"
		usage        = "convert text of URL to latex pdf"
	)
	flag.StringVar(&URL, "url", default_URL, usage)
	flag.StringVar(&URL, "u", default_URL, usage+" (shorthand)")

	flag.StringVar(&FILE, "file", default_FILE, usage)
	flag.StringVar(&FILE, "f", default_FILE, usage+" (shorthand)")
}

// ===  FUNCTION  ==========================================================
//         Name:  GetSourceFrom
//  Description:  Retrieves website source and file type.
// =========================================================================
func GetSourceFrom(url string) websiteSource {
	fmt.Println("Retrieving website from:  ", url)

	resp, err := http.Get("http://" + Trim(url))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fileType := http.DetectContentType(content)
	return websiteSource{fileType, string(content)}
} // -----  end of function GetSourceFrom  -----

// ===  FUNCTION  ==========================================================
//         Name:  Trim
//  Description:  Removes "http://" from url
// =========================================================================
func Trim(url string) string {
	return strings.TrimPrefix(url, "http://")
} // -----  end of function Trim  -----

// ===  FUNCTION  ==========================================================
//         Name:  SaveToFile
//  Description:
// =========================================================================
func SaveToFile(text string) {
	fmt.Println("Save website content to:  ", FILE)
	outfile, err := os.Create(FILE)
	if err != nil {
		panic(err)
	}
	// f.Close will run when we're finished.
	// checkin on error when closing:
	// http://stackoverflow.com/questions/1821811/how-to-read-write-from-to-file#comment20239340_9739903
	defer func() {
		if outfile.Close() != nil {
			panic(err)
		}
	}()
	io.WriteString(outfile, text)
} // -----  end of function SaveToFile  -----

// ===  FUNCTION  ==========================================================
//         Name:  GetTitle
// =========================================================================
func GetTitle(w websiteSource) string {

	//	w := websiteSource{
	//		fileType: "http",
	//		content: `<span class="article-tag">Biology<span class="article-tag-pipe">&nbsp;|&nbsp;</span><span class="article-tag-focus">Neurology & Entomology</span></span>
	//		<hgroup>
	//		<h1 class="article-title">Ants Swarm Like Brains Think</h1>
	//		<h2 class="article-subtitle">A neuroscientist studies ant colonies to understand feedback in the brain.</h2>		</hgroup>`}

	contentParseTree, _ := html.Parse(strings.NewReader(w.content))

	var f func(*html.Node, *string)
	f = func(n *html.Node, title *string) {
		if n.Attr != nil {
			if n.Attr[0].Key == "class" && n.Attr[0].Val == "article-title" {
				n = n.FirstChild
				if n.Type == html.TextNode {
					*title = n.Data
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, title)
			if *title != "" {
				break
			}
		}
	}
	var title string
	f(contentParseTree, &title)
	return title
} // -----  end of function GetTitle  -----

// ===  FUNCTION  ==========================================================
//         Name:  ExtractItem
//  Description:
// =========================================================================
func ExtractItem(source *html.Node, item html.Node) string {

	var f func(*html.Node, *string)
	f = func(n *html.Node, result *string) {
		if n.Data == item.Data {
			if n.Attr != nil {
				if n.Attr[0].Key == item.Attr[0].Key && n.Attr[0].Val == item.Attr[0].Val {
					n = n.FirstChild
					if n.Type == html.TextNode {
						*result = n.Data
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, result)
		}
	}
	var result string
	f(source, &result)
	return result
} // -----  end of function ExtractItem  -----

// ===  FUNCTION  ==========================================================
//         Name:  ExtractSearchedAttr
//  Description:
// =========================================================================
func ExtractAttribute(source *html.Node, searchedAttr html.Attribute) string {

	var f func(*html.Node, *string)
	f = func(n *html.Node, result *string) {
		if n.Attr != nil {
			if n.Attr[0].Key == searchedAttr.Key && n.Attr[0].Val == searchedAttr.Val {
				n = n.FirstChild
				if n.Type == html.TextNode {
					*result = n.Data
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, result)
		}
	}
	var result string
	f(source, &result)
	return result
} // -----  end of function ExtractSearchedAttr  -----

// ===  FUNCTION  ==========================================================
//         Name:  DisplayHtmlNode
//  Description:
// =========================================================================
func DisplayHtmlNode(n *html.Node) {
	fmt.Println("Type     : ", n.Type)
	fmt.Println("Data Atom: ", n.DataAtom)
	fmt.Println("Data     : ", n.Data)
	fmt.Println("Namespace: ", n.Namespace)
	fmt.Println("Attribute: ", n.Attr)
	if n.Attr != nil {
		fmt.Println("      Key: ", n.Attr[0].Key)
		fmt.Println("      Val: ", n.Attr[0].Val)
	}
	fmt.Println("")
} // -----  end of function DisplayHtmlNode  -----

// ===  FUNCTION  ==========================================================
//         Name:  GetParseTree
//  Description:
// =========================================================================
func PrintParseTree(w websiteSource) {
	//	w := websiteSource{
	//		fileType: "http",
	//		content: `<span class="article-tag">Biology<span class="article-tag-pipe">&nbsp;|&nbsp;</span><span class="article-tag-focus">Neurology & Entomology</span></span>
	//		<hgroup>
	//		<h1 class="article-title">Ants Swarm Like Brains Think</h1>
	//		<h2 class="article-subtitle">A neuroscientist studies ant colonies to understand feedback in the brain.</h2>		</hgroup>`}
	doc, err := html.Parse(strings.NewReader(w.content))
	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		DisplayHtmlNode(n)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
} // -----  end of function GetParseTree  -----

// ===  FUNCTION  ==========================================================
//         Name:  ClearTerminalScreen
//  Description:  Clears the terminal screen to have nice output
// =========================================================================
func ClearTerminalScreen() {
	os.Stdout.Write([]byte("\033[2J"))
	os.Stdout.Write([]byte("\033[H"))
	os.Stdout.Write([]byte("\n"))
}
