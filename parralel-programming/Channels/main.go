package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"unicode"
)

//Musician is the struct of the musicians
type Musician struct {
	Name         string
	CareerLength int
	Salary       float64
	Hash         string
}

func main() {
	//Data files
	path := "data-1.json"
	//path := "data-2.json"
	//path := "data-3.json"
	var musicians = ReadFile(path)
	mainToData := make(chan Musician)
	dataToWorker := make(chan Musician)
	workerToResults := make(chan Musician)
	workerToData := make(chan int, 3)
	resultsToMain := make(chan [25]Musician)
	closeWorker := make(chan int)
	go dataProcess(mainToData, dataToWorker, workerToData)
	go resultsProces(workerToResults, resultsToMain)
	go worker(dataToWorker, workerToData, workerToResults, closeWorker)
	go worker(dataToWorker, workerToData, workerToResults, closeWorker)
	go worker(dataToWorker, workerToData, workerToResults, closeWorker)

	for _, mus := range musicians {
		mainToData <- mus
	}
	close(mainToData)
	for i := 0; i < 3; i++ {
		<-closeWorker
	}
	close(workerToResults)
	for {
		msg, open := <-resultsToMain
		if !open {
			break
		}
		WriteFie(musicians, msg)
	}
}

//Data thread, which sends data to workers when gets message "1"
//then closes when job is done
func dataProcess(mainToData chan Musician, dataToWorker chan Musician, workerToData chan int) {
	var musicians [10]Musician
	counter := 0
	from := 0
	to := 0
	doneOrNot := false
	for {
		if counter < 10 {
			msg := <-mainToData
			if msg != (Musician{}) {
				musicians[to] = msg
				to = (to + 1) % len(musicians)
				counter++
			} else {
				doneOrNot = true
			}
		}
		if counter > 0 {
			msg := <-workerToData
			if msg == 1 {
				var musician = musicians[from]
				var emptyMusician Musician
				musicians[from] = emptyMusician
				from = (from + 1) % len(musicians)
				dataToWorker <- musician
				counter--
			}
		}
		if counter == 0 && doneOrNot == true {
			close(dataToWorker)
			break
		}
	}
}

//Worker thread sends message "1" to data thread that he wants data
//then gets the musician from data thread and computed values sends to results thread
func worker(dataToWorker chan Musician, workerToData chan int, workerToResults chan Musician, closeWorker chan int) {
	for {
		workerToData <- 1
		msg := <-dataToWorker
		if msg == (Musician{}) {
			closeWorker <- 1
			break
		}
		Musician := msg
		YesOrNo, hash := Musician.HashMusician()
		Musician.Hash = hash
		if YesOrNo {
			workerToResults <- Musician
		}
	}
}

//Put every musician from workers to struct and sends to main
func resultsProces(workerToResults chan Musician, resultsToMain chan [25]Musician) {

	var musicians [25]Musician
	for {
		msg := <-workerToResults
		if msg == (Musician{}) {
			resultsToMain <- musicians
			break
		}
		i := len(musicians) - 1
		for musicians[i].Hash < msg.Hash {
			i--
			if i < 0 {
				break
			}
			musicians[i+1] = musicians[i]
		}
		musicians[i+1] = msg
	}
	close(resultsToMain)
}
func ReadFile(path string) []Musician {

	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened ", path)
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var musicians []Musician
	json.Unmarshal(byteValue, &musicians)
	return musicians
}
func (s Musician) HashMusician() (bool, string) {

	x := fmt.Sprintf("%v", s)
	h := sha1.New()
	h.Write([]byte(x))
	bs := h.Sum(nil)
	hash := fmt.Sprintf("%x", bs)

	r := []rune(hash)

	if unicode.IsLetter(r[0]) == true && unicode.IsLetter(r[len(r)-1]) == true {

		return true, hash
	}

	return false, hash

}

//WriteFie writes data and results to file
func WriteFie(data []Musician, musicians [25]Musician) {
	file, err := os.Create("results.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	file.WriteString("Data:\n")
	i := 1
	file.WriteString("----------------------------------------------------")
	file.WriteString("\n")
	file.WriteString(fmt.Sprintf("|%-5s|%-10s|%-15s|%-17s|\n", "NR", "Name", "career length", "Salary"))
	file.WriteString("----------------------------------------------------")
	file.WriteString("\n")

	for _, musician := range data {
		s := fmt.Sprintf("|%-5d|%-10s|%-15d|%-17.3f|\n", i, musician.Name, musician.CareerLength, musician.Salary)
		i++
		writtenBytes, err := file.WriteString(s)
		if err != nil {
			fmt.Println(writtenBytes)
			log.Fatal(err)
		}
		file.WriteString("----------------------------------------------------")
		file.WriteString("\n")
	}

	file.WriteString("\n")
	file.WriteString("Results:\n")

	file.WriteString("---------------------------------------------------------------------------------------------")
	file.WriteString("\n")
	file.WriteString(fmt.Sprintf("|%-5s|%-10s|%-15s|%-17s|%-40s|\n", "NR", "Name", "career length", "salary", "Hash"))
	file.WriteString("---------------------------------------------------------------------------------------------")
	file.WriteString("\n")
	i = 1
	for _, musician := range musicians {
		if musician.CareerLength <= 0 {
			continue
		}
		//fmt.Println(i, "+++", musician)
		s := fmt.Sprintf("|%-5d|%-10s|%-15d|%-17.2f|%30.40s|\n", i, musician.Name, musician.CareerLength, musician.Salary, musician.Hash)
		i++
		writtenBytes, err := file.WriteString(s)
		if err != nil {
			fmt.Println(writtenBytes)
			log.Fatal(err)
		}
		file.WriteString("---------------------------------------------------------------------------------------------")
		file.WriteString("\n")
	}
	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
