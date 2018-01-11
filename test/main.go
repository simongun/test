package main

import (
    "fmt"
    "os"
    // "io"
    "bytes"
    "bufio"
    "net/http"
    "strconv"
    "github.com/rickar/props"
    "text/template"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main () {
    name := os.Args[1]

    file, e1 := os.Open(name)
    check (e1)
    defer file.Close();

    reader := bufio.NewReader(file)

    var props, e2 = props.Read(reader)
    check (e2)

    res := find(props)

    fmt.Println(res)
}

func find(p *props.Properties) (*Versions)  {
    patch := p.Get("PATCH")
    minor := p.Get("MINOR")
    major := p.Get("MAJOR")
    urlTemplate := p.Get("url")

    var v Versions
    patI, patE := strconv.Atoi(patch)
    minI, minE := strconv.Atoi(minor)
    majI, majE := strconv.Atoi(major)

    if major != "" {
        check(majE)
        v = Versions {0, 0, majI + 1}
       return findNewest(urlTemplate, v, majorLookupStrategy{})
    } else if minor != "" {
        check(minE)
        v = Versions {0, minI + 1, majI}
        return findNewest(urlTemplate, v, minorLookupStrategy{})
    } else if patch != "" {
        check(patE)
        v = Versions {patI + 1, minI, majI}
        return findNewest(urlTemplate, v, patchLookupStrategy{})
    } else {
     // panic, no patch version given
        panic("No version given")
    }
}

func findNewest(urlTemplate string, v Versions, s lookupStrategy) (*Versions) {
    tmpl, err := template.New("urlTemplate").Parse(urlTemplate)
    check (err)

    buf := new(bytes.Buffer)
    // fmt.Println(v)
    err = tmpl.Execute(buf, v)
    check (err)

    // check whether URL exists
    resp, err := http.Head(buf.String())
    fmt.Println(buf.String())
    fmt.Println(resp.Status)

    if (resp.StatusCode == 200) {
        incrementedV := s.increment(v)
        res := findNewest(urlTemplate, incrementedV, s)
        if res == nil {
            return &v
        } else {
            return res
        }
    } else {
        return nil
    }
}