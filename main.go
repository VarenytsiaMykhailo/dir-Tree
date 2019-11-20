package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	if err := dirTree(out, path, printFiles, nil); err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool, sliceForGraffiti []bool) (err error) {
	dirsAndFiles, err := ioutil.ReadDir(path) //инфа по содержимому в папке
	if err != nil {
		return err
	}
	var dirs []string //сюда заносим названия папок в директории Path
	var files []string //сюда заносим названия файлов в директории Path
	for _, file := range dirsAndFiles { //перебор содержимого папки
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		} else {
			files = append(files, file.Name()+" ("+strconv.Itoa(int(file.Size()))+"b)")
		}
	}
	for i, dirName := range dirs {
		if i == len(dirs)-1 {
			sliceForGraffiti = append(sliceForGraffiti, true)
		} else {
			sliceForGraffiti = append(sliceForGraffiti, false)
		}
		if printFiles {
			if len(files) > 0 {
				sliceForGraffiti[len(sliceForGraffiti)-1] = false
			}
		}
		if err := printGraffiti(out, path, dirName, sliceForGraffiti); err != nil {
			return err
		}
		if err := dirTree(out, path+"/"+dirName, printFiles, sliceForGraffiti); err != nil {
			return err
		}
		sliceForGraffiti = sliceForGraffiti[:len(sliceForGraffiti)-1]
	}
	//обработка печати файлов
	if printFiles {
		for i, filesName := range files {
			if i == len(files) - 1 {
				sliceForGraffiti = append(sliceForGraffiti, true)
			} else {
				sliceForGraffiti = append(sliceForGraffiti, false)
			}
			if err := printGraffiti(out, path, filesName, sliceForGraffiti); err != nil {
				return err
			}
			sliceForGraffiti = sliceForGraffiti[:len(sliceForGraffiti)-1]
		}
	}
	return nil
}

func printGraffiti(out io.Writer, path string, dirName string, sliceForGraffiti []bool) (err error) {
	var str string
	for i := 0; i < len(sliceForGraffiti)-1; i++ {
		if sliceForGraffiti[i] == false {
			str += "│\t"
		} else {
			str += "\t"
		}
	}
	if sliceForGraffiti[len(sliceForGraffiti)-1] == false {
		str += "├─── "
	} else {
		str += "└─── "
	}
	if _, err := fmt.Fprintf(out, "%s\n", str+dirName); err != nil {
		return err
	}
	return nil
}
