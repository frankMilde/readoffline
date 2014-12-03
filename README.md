Website 2 LaTex
===============

The aim is to have a program that takes an URL (of a website containing and
article/blogpost) and generates a two-column article pdf to read, work and
store offline.

Program UML
-----------

1. [x] Take as input URL
2. [x] Download website (and save temporarily)
3. [ ] Parse HTML to identify
  - [x] Title
	- [x] Subtitle
	- [ ] Author
	- [ ] Text body
	- [ ] Images
4. Transform to LaTeX
  - [ ] Convert <p></p> into Absatz
	- [ ] Convert quotes and apostrophes
	- [ ] Convert UTF-8 Sonderzeichen
	- [ ] Lists
	- [ ] Images and their caption
	- [ ] find `section/subsection` headings
	- [ ] latex formulas as images
	- [ ] Images
	- [ ] If no subtitle is present comment out `\subtitle` field in template
	- [ ] and send to online Latex compiler
5. [ ] Get PDF

Resources
---------

http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/
https://gist.github.com/bemasher/1224702
http://www.nokogiri.org/tutorials/parsing_an_html_xml_document.html
https://stackoverflow.com/questions/22655457/golang-xml-parse

Testing articles:
-----------------

http://science.ksc.nasa.gov/shuttle/missions/51-l/docs/rogers-commission/Appendix-F.txt
http://nautil.us/issue/9/time/haunted-by-his-brother-he-revolutionized-physics
http://blog.ycombinator.com/we-make-mistakes
http://slatestarcodex.com/2014/02/23/in-favor-of-niceness-community-and-civilization/
http://h14s.p5r.org/2012/09/0x5f3759df.html

HTML basics
-----------

The general form of an HTML element is 
```
<tag attribute1="value1" attribute2="value2">content</tag>
```
The `go` `net/http` html parser is 
```
func Parse(r io.Reader) (*Node, error)
```
`Parse` returns the parse tree for the HTML from the given Reader. The nodes
of the parse tree
```
type Node struct {
    Parent, FirstChild, LastChild, PrevSibling, NextSibling *Node

    Type      NodeType
    DataAtom  atom.Atom
    Data      string
    Namespace string
    Attr      []Attribute
}
```
consist of a `NodeType` and some `Data` (`tag name` for element nodes,
`content` for text) and are part of a tree of Nodes. Element nodes may also
have a `Namespace` and contain a slice of `Attributes`. `Data` is unescaped,
so that it looks like `"a<b"` rather than `"a&lt;b"`. For element nodes,
`DataAtom` is the atom for Data, or zero if Data is not a known tag name. 

The different `NodeType` are
```
type NodeType uint32
const (
    ErrorNode NodeType = iota
    TextNode
    DocumentNode
    ElementNode
    CommentNode
    DoctypeNode
)
```
An `Attribute` is an attribute namespace-key-value triple. `Namespace` is
non-empty for foreign attributes like `xlink`, `Key` is alphabetic (and hence
does not contain escapable characters like `'&'`, `'<'` or `'>'`), and `Val` is
unescaped (it looks like `"a<b"` rather than `"a&lt;b"`).
```
type Attribute struct {
	    Namespace, Key, Val string
}
```
