package main

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"unicode"
)

var isDone bool

//Monitor is the main structure of the program
type Monitor struct {
	cond        *sync.Cond
	mutex       *sync.Mutex
	musicians   []Musician
	currentSize int
	from        int
	to          int
	DoneOrNot   bool
}

//Musician is the structure of the musicians
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

	var wg sync.WaitGroup
	var dataMonitor = initializeMonitor(10)
	var resultsMonitor = initializeMonitor(25)
	wg.Add(1)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go MainfunctionOfTheProject(&dataMonitor, &wg, musicians, &resultsMonitor)

	}
	go AddToDataMonitor(&dataMonitor, musicians, &wg)
	wg.Wait()

	WriteFie(musicians, &resultsMonitor)
}

//WriteFie writes data and results to file
func WriteFie(data []Musician, resultsMonitor *Monitor) {
	file, err := os.Create("results.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	file.WriteString("Data:\n")
	file.WriteString("______________________________________________")
	file.WriteString("\n")
	file.WriteString(fmt.Sprintf("|%-10s|%-15s|%-17s|\n", "Name", "career length", "Salary"))
	file.WriteString("______________________________________________")
	file.WriteString("\n")

	for _, musician := range data {
		s := fmt.Sprintf("|%-10s|%-15d|%-17.3f|\n", musician.Name, musician.CareerLength, musician.Salary)
		writtenBytes, err := file.WriteString(s)
		if err != nil {
			fmt.Println(writtenBytes)
			log.Fatal(err)
		}
		file.WriteString("______________________________________________")
		file.WriteString("\n")
	}

	file.WriteString("\n")
	file.WriteString("Results:\n")
	file.WriteString("______________________________________________________________________________________")
	file.WriteString("\n")
	file.WriteString(fmt.Sprintf("|%-10s|%-15s|%-17s|%-40s|\n", "Name", "career length", "salary", "Hash"))
	file.WriteString("______________________________________________________________________________________")
	file.WriteString("\n")
	for _, musician := range resultsMonitor.musicians {
		if musician.CareerLength <= 0 {
			continue
		}
		s := fmt.Sprintf("|%-10s|%-15d|%-17.2f|%30.40s|\n", musician.Name, musician.CareerLength, musician.Salary, musician.Hash)
		writtenBytes, err := file.WriteString(s)
		if err != nil {
			fmt.Println(writtenBytes)
			log.Fatal(err)
		}
		file.WriteString("______________________________________________________________________________________")
		file.WriteString("\n")
	}
	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

//AddToDataMonitor adds data to data monitor
func AddToDataMonitor(dataMonitor *Monitor, musicians []Musician, wg *sync.WaitGroup) {
	defer wg.Done()
	dataMonitor.DoneOrNot = false
	for _, Musician := range musicians {
		dataMonitor.Insert(Musician)
	}
	dataMonitor.mutex.Lock()
	dataMonitor.DoneOrNot = true
	dataMonitor.cond.Broadcast()
	dataMonitor.mutex.Unlock()

}

//MainfunctionOfTheProject is the main function of the project
func MainfunctionOfTheProject(dataMonitor *Monitor, wg *sync.WaitGroup, musicians []Musician, resultsMonitor *Monitor) {

	defer wg.Done()

	for i := 0; i < len(musicians); i++ {

		var Musician, err = dataMonitor.Remove()

		if err != nil {
			break
		}
		YesOrNo, hash := Musician.HashMusician()
		Musician.Hash = hash

		if YesOrNo {
			resultsMonitor.InsertAndSortData(Musician)
		}

	}

}

//HashMusician hashes musicians and filters if the first and last symbol is letter
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
func initializeMonitor(monitorSize int) Monitor {
	var musiciansArray = make([]Musician, monitorSize)
	var mutex = sync.Mutex{}
	var cond = sync.NewCond(&mutex)
	var monitor = Monitor{cond: cond, mutex: &mutex, musicians: musiciansArray}
	return monitor
}

//ReadFile reads file
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

//Insert func inserts data
func (monitor *Monitor) Insert(Musician Musician) {
	monitor.mutex.Lock()
	for monitor.currentSize == len(monitor.musicians) {
		monitor.cond.Wait()
	}
	monitor.musicians[monitor.to] = Musician
	monitor.to = (monitor.to + 1) % len(monitor.musicians)
	monitor.currentSize++
	monitor.cond.Broadcast()
	monitor.mutex.Unlock()
}

//InsertAndSortData func inserts data
func (monitor *Monitor) InsertAndSortData(MusicianToInsert Musician) {
	monitor.mutex.Lock()
	for i := len(monitor.musicians) - 1; i > 0; i-- {
		for monitor.musicians[i].Hash < MusicianToInsert.Hash {
			i--
			if i < 0 {
				break
			}

			monitor.musicians[i+1] = monitor.musicians[i]
		}
		monitor.musicians[i+1] = MusicianToInsert
		monitor.currentSize++
	}

	monitor.cond.Broadcast()
	monitor.mutex.Unlock()
}

//Remove removes the Musician from the monitor
func (monitor *Monitor) Remove() (Musician, error) {
	defer monitor.mutex.Unlock()
	monitor.mutex.Lock()
	for monitor.currentSize == 0 {

		if monitor.DoneOrNot {
			var emptyMusician Musician
			return emptyMusician, errors.New("No Musician")
		}
		monitor.cond.Wait()
	}
	var musician = monitor.musicians[monitor.from]
	var emptyMusician Musician
	monitor.musicians[monitor.from] = emptyMusician
	monitor.from = (monitor.from + 1) % len(monitor.musicians)
	monitor.currentSize--
	monitor.cond.Broadcast()
	return musician, nil
}
