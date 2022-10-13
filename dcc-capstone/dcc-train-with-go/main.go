// package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"strings"
// )

// func main() {

//     content, err := ioutil.ReadFile("instruction.txt")

//      if err != nil {
//           log.Fatal(err)
//      }

// 	str_content := string(content)
// 	split_res := strings.Split(str_content, " ")
// 	cmd := split_res[0]
// 	arg1 := split_res[1]
// 	arg2 := split_res[3]
// }

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type train struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address int    `json:"address"`
}

var trains = []train{
	{ID: "1", Name: "train_1", Address: 3},
	{ID: "2", Name: "train_2", Address: 4},
	{ID: "3", Name: "train_3", Address: 5},
}

type instruction struct {
	Cmd  string `json:"cmd"`
	Arg1 string `json:"arg1"`
	Arg2 string `json:"arg2"`
}

var instructions = instruction{Cmd: "", Arg1: "", Arg2: ""}

type status struct {
	Power bool
}

func getTrains(context *gin.Context) {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	path = filepath.Join(path, "..")
	path = filepath.Join(path, "train-controller-dcc")
	print(path)

	jsonFile, err := os.Open(path + "/.dccpi")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}

	json.Unmarshal([]byte(byteValue), &result)

	context.IndentedJSON(http.StatusOK, result["locomotives"])
}

type Locomotive struct {
	Name      string `json:"name"`
	Address   int    `json:"address"`
	Direction int    `json:"direction"`
	Speed     int    `json:"speed"`
	Fl        bool   `json:"fl"`
	F1        bool   `json:"f1"`
	F2        bool   `json:"f2"`
	F3        bool   `json:"f3"`
	F4        bool   `json:"f4"`
}
type Locomotives struct {
	Locomotives []Locomotive `json:"locomotives"`
}

func commandTest(context *gin.Context) {
	var newInstruct instruction

	if err := context.BindJSON(&newInstruct); err != nil {
		return
	}

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	path = filepath.Join(path, "..")
	path = filepath.Join(path, "train-controller-dcc")
	// print(path)

	jsonFile, err := os.Open(path + "/.dccpi")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result Locomotives

	json.Unmarshal([]byte(byteValue), &result)

	// fmt.Println("\n")
	// for i := 0; i < len(result.Locomotives); i++ {
	// 	fmt.Println("Name Train : " + result.Locomotives[i].Name)
	// 	fmt.Println("Train Address: " + strconv.Itoa(result.Locomotives[i].Address))
	// 	fmt.Println("Direction Train: " + strconv.Itoa(result.Locomotives[i].Direction))
	// 	fmt.Println("Front Light: " + strconv.FormatBool(result.Locomotives[i].Fl))
	// 	fmt.Println("Speed: " + strconv.Itoa(result.Locomotives[i].Speed))
	// }

	// fmt.Println(result.Locomotives)
	fmt.Println(newInstruct)

	if newInstruct.Cmd == "register" {
		var newData Locomotive
		newData.Name = newInstruct.Arg1
		newData.Address, _ = strconv.Atoi(newInstruct.Arg2)
		newData.Speed = 0
		newData.Direction = 0
		newData.Fl = true
		newData.F1 = false
		newData.F2 = false
		newData.F3 = false
		newData.F4 = false
		newArray := append(result.Locomotives, newData)
		result.Locomotives = newArray
		payloadJSON, _ := json.Marshal(result)
		fmt.Println(payloadJSON)
		ioutil.WriteFile(path+"/.dccpi", payloadJSON, 0664)
	} else if newInstruct.Cmd == "unregister" {
		deleteIndex := 0
		for i, data := range result.Locomotives {
			if data.Name == newInstruct.Arg1 {
				deleteIndex = i
			} else {
				continue
			}
		}
		fmt.Println(deleteIndex)

		var unregArray Locomotives

		unregArray1 := result.Locomotives[:deleteIndex]
		unregArray2 := result.Locomotives[deleteIndex+1:]

		newArray := append(unregArray.Locomotives, unregArray1...)
		newArray = append(newArray, unregArray2...)
		fmt.Println(newArray)
		result.Locomotives = newArray
		payloadJSON, _ := json.Marshal(result)
		fmt.Println(payloadJSON)
		ioutil.WriteFile(path+"/.dccpi", payloadJSON, 0664)
	} else if newInstruct.Cmd == "speed" {
		tempSpeed, _ := strconv.Atoi(newInstruct.Arg2)
		indexSpeed := 0
		for i, data := range result.Locomotives {
			if data.Name == newInstruct.Arg1 {
				indexSpeed = i
			} else {
				continue
			}
		}
		fmt.Println(tempSpeed)
		fmt.Println(result.Locomotives)
		result.Locomotives[indexSpeed].Speed = tempSpeed
		payloadJSON, _ := json.Marshal(result)
		fmt.Println(payloadJSON)
		ioutil.WriteFile(path+"/.dccpi", payloadJSON, 0664)
	} else if newInstruct.Cmd == "fl" {
		lightCmd := newInstruct.Arg2
		indexFl := 0
		stateFl := false
		for i, data := range result.Locomotives {
			if data.Name == newInstruct.Arg1 {
				indexFl = i
			} else {
				continue
			}
		}
		if lightCmd == "on" {
			stateFl = true
		} else if lightCmd == "off" {
			stateFl = false
		}
		result.Locomotives[indexFl].Fl = stateFl
		payloadJSON, _ := json.Marshal(result)
		fmt.Println(payloadJSON)
		ioutil.WriteFile(path+"/.dccpi", payloadJSON, 0664)
	} else if newInstruct.Cmd == "power" {
		cmd := newInstruct.Cmd
		arg1 := newInstruct.Arg1
		arg2 := newInstruct.Arg2
		fullCommand := ""
		if arg1 == "" {
			fullCommand = "send," + cmd
		} else if arg2 == "" {
			fullCommand = "send," + cmd + "," + arg1
		} else {
			fullCommand = "send," + cmd + "," + arg1 + "," + arg2
		}

		byteCommand := []byte(fullCommand)
		fmt.Println(cmd)
		fmt.Println(arg1)
		fmt.Println(arg2)

		path, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		path = filepath.Join(path, "..")
		fmt.Println(path)

		ioutil.WriteFile(path+"/instruction.txt", byteCommand, 0644)

	} else if newInstruct.Cmd == "direction" {

		payloadJSON, _ := json.Marshal(result)
		fmt.Println(payloadJSON)
		ioutil.WriteFile(path+"/.dccpi", payloadJSON, 0664)
	} else {

	}

	context.IndentedJSON(http.StatusAccepted, newInstruct)
}

func addTrain(context *gin.Context) {
	var newTrain train

	if err := context.BindJSON(&newTrain); err != nil {
		return
	}

	trains = append(trains, newTrain)

	context.IndentedJSON(http.StatusCreated, newTrain)
}

func getTrain(context *gin.Context) {
	id := context.Param("id")
	train, err := getTrainById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": err})
	}
	context.IndentedJSON(http.StatusOK, train)
}

func getTrainById(id string) (*train, error) {
	for i, t := range trains {
		if t.ID == id {
			return &trains[i], nil
		}
	}
	return nil, errors.New("Train not found")
}

func getInstruction(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, instructions)
}

func editInstruction(context *gin.Context) {
	var newInstruct instruction

	if err := context.BindJSON(&newInstruct); err != nil {
		return
	}

	instruction := &instructions
	instruction.Cmd = newInstruct.Cmd
	instruction.Arg1 = newInstruct.Arg1
	instruction.Arg2 = newInstruct.Arg2

	context.IndentedJSON(http.StatusOK, instruction)
}

func sendIntruction(context *gin.Context) {
	instruction := &instructions
	cmd := instruction.Cmd
	arg1 := instruction.Arg1
	arg2 := instruction.Arg2
	fullCommand := ""
	if arg1 == "" {
		fullCommand = "send," + cmd
	} else if arg2 == "" {
		fullCommand = "send," + cmd + "," + arg1
	} else {
		fullCommand = "send," + cmd + "," + arg1 + "," + arg2
	}

	byteCommand := []byte(fullCommand)
	fmt.Println(cmd)
	fmt.Println(arg1)
	fmt.Println(arg2)

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	path = filepath.Join(path, "..")
	fmt.Println(path)

	ioutil.WriteFile(path+"/instruction.txt", byteCommand, 0)
	// ioutil.WriteFile("instruction.txt",byteCommand, 0)
	// ioutil.WriteFile("instruction.txt", byteArg2, 0)
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func GetInternalIP() string {
	itf, _ := net.InterfaceByName("enp1s0") //here your interface
	item, _ := itf.Addrs()
	var ip net.IP
	for _, addr := range item {
		switch v := addr.(type) {
		case *net.IPNet:
			if !v.IP.IsLoopback() {
				if v.IP.To4() != nil { //Verify if IP is IPV4
					ip = v.IP
				}
			}
		}
	}
	if ip != nil {
		return ip.String()
	} else {
		return ""
	}
}

func main() {
	router := gin.Default()
	router.GET("/train", getTrains)
	router.GET("/train/:id", getTrain)
	router.POST("/command", commandTest)
	router.POST("/train", addTrain)
	router.GET("/instruct", getInstruction)
	router.PATCH("/instruct", editInstruction)
	router.GET("/send", sendIntruction)
	fmt.Println("Serving the server with IP : ", GetLocalIP(), ":8000")
	router.Run(":8000")
}
