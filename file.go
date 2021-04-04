package main

import (
	"archive/zip"
	"fmt"
	"github.com/kataras/iris/v12/sessions"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
)

const (
	maxSize   = 1 * iris.GB
	uploadDir = "./userdata/userfiles"
)

func init() {
	os.Mkdir("./userdata/userfiles", 0700)
}

type FileDetail struct {
	FileName     string    `json:"fileName-js"`
	LastModiTime string `json:"lastModiTime-js"`
	Size         int64     `json:"size-js"`
	IsDir        bool      `json:"isDir-js"`
}

type FileListData struct {
	FileDetails []FileDetail `json:"fileDetails-js"`
	FilePath    string       `json:"filePath-js"`
	ParentPath  string       `json:"parentPath-js"`
}

type MDStream struct {
	FileName   string `json:"fileName-js"`
	DataStream string `json:"dataStream-js"`
}

type FileResponse struct {
	IsNewFlag     bool      `json:"isNewFile"`
	ResponseData  MDStream  `json:"JSONData"`
}

func getFile(ctx iris.Context){
	if auth, _ := sessions.Get(ctx).GetBoolean("authenticated"); !auth{
		ctx.JSON(iris.Map{
			"redirect": true,
		})
		return
	}
	var responseMessage FileResponse
	getFileName := ctx.URLParam("FileName")
	if getFileName == ""{
		responseMessage = FileResponse{
			IsNewFlag: true,
			ResponseData: MDStream{
				FileName:   "",
				DataStream: "",
			},
		}
	} else{
		checkFileName := getFileName
		if string(getFileName[0]) != "." {
			checkFileName = uploadDir + "/" + getFileName
		}

		if _, err := os.Stat(checkFileName); os.IsNotExist(err) {
			responseMessage = FileResponse{
				IsNewFlag: true,
				ResponseData: MDStream{
					FileName:   "",
					DataStream: "",
				},
			}
		}else{
			mdData, _ := ioutil.ReadFile(checkFileName)
			checkFileName = strings.TrimSuffix(path.Base(checkFileName), path.Ext(path.Base(checkFileName)))
			responseMessage = FileResponse{
				IsNewFlag: false,
				ResponseData: MDStream{
					FileName: checkFileName,
					DataStream: string(mdData),
				},
			}
		}
	}
	ctx.JSON(responseMessage)
}
func saveFile(ctx iris.Context){
	if auth, _ := sessions.Get(ctx).GetBoolean("authenticated"); !auth{
		ctx.JSON(iris.Map{
			"redirect": true,
		})
		return
	}
	var MDData MDStream
	err := ctx.ReadJSON(&MDData)
	if err != nil{
		fmt.Println("JSON Format Error")
		return
	} else {
		pathWithFile := MDData.FileName
		pathWithFile = uploadDir + "/" + MDData.FileName
		if _, err := os.Stat(pathWithFile); os.IsNotExist(err) {
			basePath := path.Dir(pathWithFile)
			//filename := path.Base(pathWithFile)
			os.Mkdir(basePath, 0755)
		}
		err2 := ioutil.WriteFile(pathWithFile, []byte(MDData.DataStream), 0755)
		if err2 != nil {
			return
		}
	}
	ctx.JSON(iris.Map{
		"result":true,
	})
}

func getFilePath(ctx iris.Context) {
	if auth, _ := sessions.Get(ctx).GetBoolean("authenticated"); !auth{
		ctx.JSON(iris.Map{
			"redirect": true,
		})
		return
	}
	ctx.JSON(iris.Map{
		"presentPath": uploadDir,
	})
}

func showAllFiles(ctx iris.Context) {
	if auth, _ := sessions.Get(ctx).GetBoolean("authenticated"); !auth{
		ctx.JSON(iris.Map{
			"redirect": true,
		})
		return
	}
	changedPath:= ctx.URLParam("changePathRequest")

	filepathNames,err := filepath.Glob(filepath.Join(changedPath,"*"))
	if err != nil {
		log.Fatal(err)
	}

	var fileDetailData []FileDetail
	for i := range filepathNames {
		info, err := os.Stat(filepathNames[i])
		if err != nil {
			fmt.Println("os.Stat err =",err)
			return
		}
		fileDetailData = append(fileDetailData,FileDetail{
			FileName:     info.Name(),
			LastModiTime: info.ModTime().Format(time.RFC822),
			Size:         info.Size(),
			IsDir:        info.IsDir(),
		})
	}

	test := FileListData{
		FileDetails: fileDetailData,
		FilePath:    uploadDir,
		ParentPath:  uploadDir,
	}
	ctx.JSON(test)
}

type FileList struct {
	SelectList []string `json:"Request"`
	SelectPath string   `json:"Path"`
}

func downloadFile(ctx iris.Context){
	if auth, _ := sessions.Get(ctx).GetBoolean("authenticated"); !auth{
		ctx.JSON(iris.Map{
			"redirect": true,
		})
		return
	}
	var selectList FileList
	output := uploadDir + "/" + "Download.zip"
	if err := os.RemoveAll(output); err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}
	err := ctx.ReadJSON(&selectList)
	if err != nil {
		fmt.Println("JSON Format error")
		return
	}
	for i:=0; i < len(selectList.SelectList); i++{
		selectList.SelectList[i] = selectList.SelectPath + "/" + selectList.SelectList[i]
	}
	if err := ZipFiles(output, selectList.SelectList); err != nil {
		panic(err)
	}
	err = os.Chmod(output, 0755)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Zipped File:", output)
	dest := "" /* optionally, keep it empty to resolve the filename based on the "src" */

	// Limit download speed to ~50Kb/s with a burst of 100KB.
	limit := 50.0 * iris.KB
	burst := 100 * iris.KB
	ctx.SendFileWithRate("./userdata/userfiles/userfiles.zip", dest, limit, burst)
}

func deleteFile(ctx iris.Context) {
	if auth, _ := sessions.Get(ctx).GetBoolean("authenticated"); !auth{
		ctx.JSON(iris.Map{
			"redirect": true,
		})
		return
	}
	// It does not contain the system path,
	// as we are not exposing it to the user.
	var selectList FileList
	err := ctx.ReadJSON(&selectList)
	fmt.Println(selectList)
	if err != nil {
		fmt.Println("JSON Format error")
		return
	}
	for i := 0; i < len(selectList.SelectList); i++{
		filePath := path.Join(selectList.SelectPath, selectList.SelectList[i])
		if err := os.RemoveAll(filePath); err != nil {
			ctx.StopWithError(iris.StatusInternalServerError, err)
			return
		}
	}
	ctx.JSON(iris.Map{
		"result":true,
	})
}

// ZipFiles compresses one or many files into a single zip archive file.
// Param 1: filename is the output zip file's name.
// Param 2: files is a list of files to add to the zip.
func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	fmt.Println(fileToZip)
	fmt.Println(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}