package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
/*sliceForGraffiti - слайс, который требуется для корректного, красивого вывода графики.
По лежащим в нем булевым значениям функция printGraffiti определяет, где выводить "│	", а где просто "	",
где выводить "├───", а где выводить "└───" .
В слайсе хранится true, если текущий рассматриваемый файл или папка является последним в текущей директории
т.е. является последним элементом слайса dirs или files.
Логика этого механизма тяжело продумана и не рекомендуется для внесения изменений.*/

//выводит дерево каталога в отсортированном по имени виде (сначала выводятся подпапки, а потом подфайлы)
func dirTree(out io.Writer, path string, printFiles bool, sliceForGraffiti []bool) (err error) {
	dirsAndFiles, err := ioutil.ReadDir(path) //инфа по содержимому в текущей папке (получаемый слайс - уже в отсортированном по имени виде)
	if err != nil {
		return err
	}
	var dirs []string //сюда заносим названия папок в директории Path
	var files []string //сюда заносим названия файлов в директории Path
	for _, file := range dirsAndFiles { //перебор содержимого текущей папки
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		} else { //если это файл, а не папка
			files = append(files, file.Name()+" ("+strconv.Itoa(int(file.Size()))+"b)")
		}
	}
	//обработка печати папок (не файлов)
	for i, dirName := range dirs {
		if i == len(dirs)-1 {
			sliceForGraffiti = append(sliceForGraffiti, true)
		} else {
			sliceForGraffiti = append(sliceForGraffiti, false)
		}
		if printFiles { //нужно для ключа -f
			if len(files) > 0 {
				sliceForGraffiti[len(sliceForGraffiti)-1] = false
			}
		}
		printGraffiti(out, dirName, sliceForGraffiti)
		if err := dirTree(out, path+"/"+dirName, printFiles, sliceForGraffiti); err != nil { //рекурсивный вызов (переход в следующую папку текущей директории)
			return err
		}
		sliceForGraffiti = sliceForGraffiti[:len(sliceForGraffiti)-1]
	}
	//обработка печати файлов и их размеров (не папок)
	if printFiles {
		for i, filesName := range files {
			if i == len(files) - 1 {
				sliceForGraffiti = append(sliceForGraffiti, true)
			} else {
				sliceForGraffiti = append(sliceForGraffiti, false)
			}
			printGraffiti(out, filesName, sliceForGraffiti)
			sliceForGraffiti = sliceForGraffiti[:len(sliceForGraffiti)-1]
		}
	}
	return nil
}

//отвечает за печать графики
func printGraffiti(out io.Writer, dirOrFileName string, sliceForGraffiti []bool) {
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
	if _, err := fmt.Fprintf(out, "%s\n", str+dirOrFileName); err != nil {
		log.Fatalf("Can't drow graffiti ", err.Error())
	}
}
