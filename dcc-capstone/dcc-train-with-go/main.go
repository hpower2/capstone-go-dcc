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
	path = filepath.Join(path, "train-controller-dcc-main")

	jsonFile, err := os.Open(path + "/.dccpi")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}

	json.Unmarshal([]byte(byteValue), &result)

	context.IndentedJSON(http.StatusOK, result["locomotives"])
}

func commandTest(context *gin.Context) {
	var newInstruct instruction

	if err := context.BindJSON(&newInstruct); err != nil {
		return
	}

	fmt.Println(newInstruct)

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
