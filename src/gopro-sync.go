package main

import (
	"fmt"
	"net/http"
	"os"
	"io/ioutil"
	"code.google.com/p/go.net/html"
	"bytes"
	"strings"
	"io"
	"github.com/cheggaaa/pb"
	"strconv"
)

const (
	//GO_PRO_URL = "http://kevoree.org/"
	GO_PRO_URL = "http://10.5.5.9:8080/DCIM/100GOPRO/"
)

func main() {
	dwl := new(GoProDwl)
	remoteFiles := dwl.filterNonExisting(dwl.listFiles())
	if(len(remoteFiles)==0){
		fmt.Println("GoPro in sync, no files downloaded")
		os.Exit(0)
	}
	bar := pb.StartNew(len(remoteFiles));
	var total int64 = 0
	for i := range remoteFiles {
		bar.Increment()
		var remoteFile = remoteFiles[i];
		output, error := os.Create(remoteFile)
		if error != nil {
			fmt.Println("Error while creating", remoteFile, "-", error)
			return
		}
		defer output.Close()

		response, error := http.Get(GO_PRO_URL + remoteFile)
		if error != nil {
			fmt.Println("Error while downloading", GO_PRO_URL+remoteFile, "-", error)
			return
		}
		defer response.Body.Close()
		n, error := io.Copy(output, response.Body)
		total = total + n
		if error != nil {
			fmt.Println("Error while downloading", GO_PRO_URL+remoteFile, "-", error)
			return
		}
	}
	bar.FinishPrint("GoPro in sync "+(strconv.FormatInt(total/1000000,10))+" Mbytes downloaded")
}

func (dwl *GoProDwl) filterNonExisting(remoteFiles []string) []string {
	var listFiles []string
	for i := range remoteFiles {
		var remoteFile = remoteFiles[i];
		localFile, err := os.OpenFile(remoteFile, os.O_RDONLY, os.ModeDevice)
		if err != nil {
			listFiles = append(listFiles, remoteFile)
		} else {
			localFile.Close()
		}
	}
	return listFiles;
}

func (dwl *GoProDwl) listFiles() []string {
	var listFiles []string
	response, err := http.Get(GO_PRO_URL)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		d := html.NewTokenizer(bytes.NewReader(contents))
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		for {
			tokenType := d.Next()
			if tokenType == html.ErrorToken {
				return listFiles;
			}
			token := d.Token()
			for i := range token.Attr {
				att := token.Attr[i]
				if (att.Key == "href") {
					if (strings.HasSuffix(att.Val, "JPG") || strings.HasSuffix(att.Val, "MP4")) {
						listFiles = append(listFiles, att.Val)
					}
				}
			}
		}


	}
	return listFiles;
}

type GoProDwl struct {

}
