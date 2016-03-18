package main

//import "fmt"
import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

//文件操作
func GetLatestNumberFromFile() (oldNumber [6]int) {
	userFile := "Number.txt"
	fin, err := os.OpenFile(userFile, os.O_RDWR|os.O_CREATE, 0666)
	defer fin.Close()
	if err != nil {
		fmt.Println(userFile, err)
		os.Exit(0)
	}

	br := bufio.NewReader(fin)
	var isp error
	var a []byte
	var ret []byte
	for isp == nil {
		a, _, isp = br.ReadLine()
		if len(a) != 0 {
			ret = a
		}

	}
	fmt.Println("Get latest number :", string(ret))
	ss := strings.Split(string(ret), " ")
	for i, v := range ss {
		oldNumber[i], _ = strconv.Atoi(v)

	}
	return
}
func WriteLatestNumbers(Numbers [6]int) {
	userFile := "Number.txt"
	fin, err := os.OpenFile(userFile, os.O_RDWR|os.O_CREATE, 0666)
	defer fin.Close()
	if err != nil {
		fmt.Println(userFile, err)
		os.Exit(0)
	}
	bw := bufio.NewWriter(fin)
	var str string
	for i, v := range Numbers {
		if i != 5 {
			str = str + strconv.Itoa(v) + " "
		} else {
			str = str + strconv.Itoa(v)
		}

	}
	//str := strconv.Itoa(Numbers[0]) + strconv.Itoa(Numbers[1]) + strconv.Itoa(Numbers[2]) + strconv.Itoa(Numbers[3]) + strconv.Itoa(Numbers[4]) + strconv.Itoa(Numbers[5])
	bw.WriteString(str + "\n")
	bw.Flush()
}

//三区原则，1-11,12-22,23-33，每个区至少1，至多3个,123,132,213,231,312,321
//2+2+2 15% 3个区各生成1个再随机生成3个然后判断某个区是否大于4个
func CreateAndSplitBlock() []int {
	rand.Seed(time.Now().UnixNano())
	Number := make([]int, 6)
	Number[0] = GenBlockNumber(1, 11)
	//fmt.Println("a is: ", Number[0])
	Number[1] = GenBlockNumber(12, 22)
	//fmt.Println("b is: ", Number[1])
	Number[2] = GenBlockNumber(23, 33)
	//fmt.Println("c is: ", Number[2])
	Number = GenLast3Number(Number)
	sort.Ints(Number)

	return Number

}
func CreateAndSplitBlockFor3(oldnumber [6]int) (a, b, c []int) {
	var Number [6]int
	var isPass bool
	for {
		//var threeNumber := make([][]int 3)
		Number[0] = GenBlockNumber(1, 11)
		//fmt.Println("a is: ", Number[0])
		Number[1] = GenBlockNumber(12, 22)
		//fmt.Println("b is: ", Number[1])
		Number[2] = GenBlockNumber(23, 33)
		//fmt.Println("c is: ", Number[2])
		threeNumber := Gen3BlockNumber(Number)
		a = threeNumber[0][:]
		//fmt.Println("a is: ", a)
		sort.Ints(a)
		isPass, _ = ValidateNumbers(a, oldnumber)
		if !isPass {
			continue
		}
		b = threeNumber[1][:]
		//fmt.Println("b is: ", b)
		sort.Ints(b)
		isPass, _ = ValidateNumbers(b, oldnumber)
		if !isPass {
			continue
		}
		c = threeNumber[2][:]
		//fmt.Println("c is: ", c)
		sort.Ints(c)
		isPass, _ = ValidateNumbers(c, oldnumber)
		if !isPass {
			continue
		}
		break
	}

	// for _, v := range threeNumber {
	// 	sort.Ints(v)
	// }
	//sort.Ints(Number)

	return

}

//统计号码分区
//return "2xx" "x2x" "xx2" "222" "other"
func GetBlockType(number []int) (blocktype string) {
	var g, s, b int
	for _, i := range number {
		switch {
		case i <= 11:
			g++
		case i <= 22:
			s++
		case i <= 33:
			b++
		}
	}
	if g >= 4 || s >= 4 || b >= 4 {
		blocktype = "other"
	} else if g == 2 && s == 2 && b == 2 {
		blocktype = "222"
	} else if g == 2 && s != 2 {
		blocktype = "2xx"
	} else if g != 2 && s == 2 {
		blocktype = "x2x"
	} else if g != 2 && b == 2 {
		blocktype = "xx2"
	}
	return
}
func Gen3BlockNumber(number [6]int) [3][6]int {
	var ret [3][6]int
	var a, b, c, d int
	for {
		a = GenBlockNumber(1, 11)
		if a != number[0] {
			break
		}
	}
	//fmt.Println("a is: ", a)
	for {
		b = GenBlockNumber(12, 22)
		if b != number[1] {
			break
		}
	}
	//fmt.Println("b is: ", b)
	for {
		c = GenBlockNumber(23, 33)
		if c != number[2] {
			break
		}
	}
	//fmt.Println("c is: ", c)
	for {
		d = GenNumber(33)
		switch d {
		case a, b, c, number[0], number[1], number[2]:
			continue
		default:
			break
		}
		break
	}
	ret[0][0] = number[0]
	ret[0][1] = number[1]
	ret[0][2] = number[2]
	ret[0][3] = a
	ret[0][4] = b
	ret[0][5] = c

	ret[1][0] = number[0]
	ret[1][1] = number[1]
	ret[1][2] = number[2]
	ret[1][3] = a
	ret[1][4] = c
	ret[1][5] = d

	ret[2][0] = number[0]
	ret[2][1] = number[1]
	ret[2][2] = number[2]
	ret[2][3] = b
	ret[2][4] = d
	ret[2][5] = c
	//fmt.Println(ret[0], ret[1], ret[2])
	return ret
}
func GenLast3Number(number []int) []int {
	var a, b, c int
	for {

		for i := 0; i != 1; {
			a = GenNumber(33)
			switch a {
			case number[0], number[1], number[2]:

				//fmt.Println("a is: ", a)
				continue
			default:
				//fmt.Println("aa is: ", a)
				number[3] = a
				i = 1
				break
			}

		}

		for i := 0; i != 1; {
			b = GenNumber(33)
			switch b {
			case number[0], number[1], number[2], number[3]:

				//fmt.Println("b is: ", b)
				continue
			default:
				//fmt.Println("bb is: ", b)
				number[4] = b
				i = 1
				break
			}
		}

		for i := 0; i != 1; {
			//fmt.Println("c is: ", c)
			c = GenNumber(33)
			switch c {
			case number[0], number[1], number[2], number[3], number[4]:

				continue
			default:
				number[5] = c
				i = 1
				break
			}
		}
		//fmt.Println("a is: ", a, " b is: ", b, " c is: ", c)
		if 11 >= a && a > 0 {
			if (11 >= b && b > 0) && (11 >= c && c > 0) {
				continue
			} else {
				break
			}
		}
		if 22 >= a && a > 11 {
			if (22 >= b && b > 11) && (22 >= c && c > 11) {
				continue
			} else {
				break
			}
		}
		if 33 >= a && a > 22 {
			if (33 >= b && b > 22) && (33 >= c && c > 22) {
				continue
			} else {
				break
			}
		}
	}
	return number
}
func CalcAbs(a int) (ret int) {
	ret = (a ^ a>>31) - a>>31
	return
}

// func GenBlockNumber(max int) int {
// 	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	//rand.Seed(time.Now().Unix())
// 	num := rand.Intn(max)
// 	//fmt.Println("num is: ", num)
// 	if num <= max-11 {
// 		num = GenBlockNumber(max)
// 	}
// 	return num
// }
func GenBlockNumber(min, max int) (ret int) {
	response, ok := http.Get("https://www.random.org/integers/?num=1&min=" + strconv.Itoa(min) + "&max=" + strconv.Itoa(max) + "&col=1&base=10&format=plain&rnd=new")
	if ok != nil {
		fmt.Println("get failed!")
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	//fmt.Println(string(body))
	list := strings.Split(string(body), "\n")
	as, error := strconv.Atoi(list[0])
	if error != nil {
		fmt.Println("字符串转换成整数失败")
	}
	ret = as

	return
}

// func GenNumber(max int) int {
// 	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	//rand.Seed(time.Now().Unix())
// 	num := rand.Intn(max)
// 	for num == 0 {
// 		num = rand.Intn(max)
// 	}

// 	return num
// }
func GenNumber(max int) (ret int) {

	response, ok := http.Get("https://www.random.org/integers/?num=1&min=1&max=" + strconv.Itoa(max) + "&col=1&base=10&format=plain&rnd=new")
	if ok != nil {
		fmt.Println("get failed!")
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	//fmt.Println(string(body))
	list := strings.Split(string(body), "\n")
	as, error := strconv.Atoi(list[0])
	if error != nil {
		fmt.Println("字符串转换成整数失败")
	}
	ret = as

	return
}

//连号2-3个
func IsConsequent(Number []int) bool {
	csCounter := 0
	if Number[1]-Number[0] == 1 {
		csCounter++
	}
	if Number[2]-Number[1] == 1 {
		csCounter++
	}
	if Number[3]-Number[2] == 1 {
		csCounter++
	}
	if Number[4]-Number[3] == 1 {
		csCounter++
	}
	if Number[5]-Number[4] == 1 {
		csCounter++
	}
	if csCounter > 2 {
		return false
	} else {
		return true
	}
}

//首号和末号不重复
func FirstAndLastNotDitto(newNumber []int, oldNumber [6]int) bool {
	if newNumber[0] == oldNumber[0] || newNumber[5] == oldNumber[5] {
		return false
	} else {
		return true
	}
}

//1-2个重号
func DittoNumber(newNumber []int, oldNumber [6]int) bool {
	dittoCounter := 0
	for _, i := range newNumber {
		if OneInSix(i, oldNumber) {
			dittoCounter++
		}
	}
	if dittoCounter > 2 {
		return false
	} else {
		return true
	}

}
func OneInSix(r1 int, oldNumber [6]int) bool {
	switch r1 {
	case oldNumber[0], oldNumber[1], oldNumber[2], oldNumber[3], oldNumber[4], oldNumber[5]:
		return true

	default:
		return false
	}

}

//和值，79-125,best:102-111
func AllSum(number []int) bool {
	suma := number[0] + number[1] + number[2] + number[3] + number[4] + number[5]
	if suma >= 102 && suma <= 116 {
		return true
	} else {
		return false
	}
}

//尾和23-28
func TailSum(number []int) bool {
	unit1, _ := SplitTail(number[0])
	unit2, _ := SplitTail(number[1])
	unit3, _ := SplitTail(number[2])
	unit4, _ := SplitTail(number[3])
	unit5, _ := SplitTail(number[4])
	unit6, _ := SplitTail(number[5])
	sumt := unit1 + unit2 + unit3 + unit4 + unit5 + unit6
	if sumt >= 23 && sumt <= 28 {
		return true
	} else {
		return false
	}
}

//拆分十位和个位
func SplitTail(bollnumber int) (unit, ten int) {
	if bollnumber > 9 {
		return bollnumber % 10, bollnumber / 10

	} else {
		return bollnumber, 0

	}
	return 0, 0

}

//奇偶比 2:4,4:2,3:3 ,3种比例1:1:1
func Parity(Number []int) bool {
	parCounter := 0
	for _, i := range Number {
		if i%2 != 0 {
			parCounter++
		}
	}
	if parCounter < 2 || parCounter > 4 {
		return false
	} else {
		return true
	}
}
func ValidateNumbers(newNumber []int, oldNumber [6]int) (pass bool, reason string) {
	//var Number []int
	var newBlockType string
	oldBlockType := GetBlockType(oldNumber[:])
	//for {
	//Number = CreateAndSplitBlock()
	newBlockType = GetBlockType(newNumber)
	if oldBlockType != newBlockType {
		if FirstAndLastNotDitto(newNumber, oldNumber) {

			if IsConsequent(newNumber) {
				//fmt.Println("连号 OK")
				//fmt.Println(Number)
				if AllSum(newNumber) {
					//fmt.Println("和值 OK")
					if TailSum(newNumber) {
						//fmt.Println("尾和 OK")
						if Parity(newNumber) {
							//fmt.Println("奇偶比 OK")
							if DittoNumber(newNumber, oldNumber) {
								//fmt.Println("重号 OK")
								//fmt.Println("ALL OK")
								b := GenBlockNumber(1, 16)
								fmt.Println(newNumber, b, " ok...")

							} else {

								return false, "重号<2"
							}
						} else {
							return false, "奇偶比24,42,33"
						}
					} else {
						return false, "尾和 23-28"
					}
				} else {

					return false, "和值102-116"
				}
			} else {
				//fmt.Println("连号 NOK")
				//fmt.Println(Number)

				return false, "连号2"
			}
		} else {
			return false, "第一个号最后一个号重复"
		}
	} else {
		return false, "分区重复"
	}

	//}
	return true, ""
}

// func main() {
// 	fmt.Println(CreateAndSplitBlock())
// }
func ProcessGenerate(oldNumbers [6]int) {
	var numberOftimes int

	fmt.Println("How many Numbers are you want to gen?")
	_, scanerr := fmt.Scanln(&numberOftimes)
	if scanerr != nil || numberOftimes <= 0 {
		fmt.Println("Error input, process exit...")
		return
	}

	//oldnumber := []int{9, 11, 13, 22, 24, 26}
	fmt.Println("oldNumbers is : ", oldNumbers)
	for numberOftimes > 0 {
		//for i := 0; i != 3; i++ {
		//InteGrate(oldNumbers)

		for {
			Number := CreateAndSplitBlock()
			pass, _ := ValidateNumbers(Number, oldNumbers)
			if pass {
				break
			}
		}
		time.Sleep(time.Duration(500))
		//}

		numberOftimes--
	}

}
func VerifyNumbers(oldNumber [6]int) {
	var n1, n2, n3, n4, n5, n6 int
	fmt.Println("Please input the numbers(6) what need to verify as: 1 2 3 4 5 6 and press Enter...")
	_, err := fmt.Scanln(&n1, &n2, &n3, &n4, &n5, &n6)
	if err != nil {
		fmt.Println("Error input, process exit...")
		return
	}
	CheckingNumbers := []int{n1, n2, n3, n4, n5, n6}
	pass, reason := ValidateNumbers(CheckingNumbers, oldNumber)
	if !pass {
		fmt.Println(CheckingNumbers, " are not follow the rule...", reason)
	}

}
func InputLatestNumbers() (old [6]int, isInput bool) {
REPEAT:
	fmt.Println("Please input the lastest numbers(6) as: 1 2 3 4 5 6 and press Enter...")

	_, err := fmt.Scanln(&old[0], &old[1], &old[2], &old[3], &old[4], &old[5])
	if err != nil {
		goto REPEAT
	}

	return old, true
}
func GenOneNumber() {

}
func main() {
	var oldNumbers [6]int
	//var LatestExist bool
	var ControlKey int
	for {

		fmt.Println("\nPress digital key(1-4) to continue...\n1.Input latest numbers;\n2.Verify Numbers;\n3.Generate Numbers;\n4.Gen blue Number;\n5.Gen3numbers;\n6.Exit...")
		fmt.Scanln(&ControlKey)
		switch ControlKey {
		case 1:
			oldNumbers, _ = InputLatestNumbers()
			WriteLatestNumbers(oldNumbers)
		case 2:
			oldNumbers = GetLatestNumberFromFile()
			// if LatestExist {
			VerifyNumbers(oldNumbers)
			// } else {
			// 	fmt.Print("Please input the latest Numbers first.")
			// 	time.Sleep(500 * time.Millisecond)
			// 	fmt.Print(".")
			// 	time.Sleep(500 * time.Millisecond)
			// 	fmt.Print(".\n")
			// 	continue
			// }

		case 3:
			oldNumbers = GetLatestNumberFromFile()
			//if LatestExist {
			ProcessGenerate(oldNumbers)
			// } else {

			// 	fmt.Print("Please input the latest Numbers first.")
			// 	time.Sleep(500 * time.Millisecond)
			// 	fmt.Print(".")
			// 	time.Sleep(500 * time.Millisecond)
			// 	fmt.Print(".\n")
			// 	continue
			// }
		case 4:
			var min, max, late int
			fmt.Println("Please input the MinNumer , MaxNumber , latest Number and press Enter...")
			fmt.Scanln(&min, &max, &late)
			for i := GenBlockNumber(min, max); ; {
				if i != late {
					fmt.Println(i)
					break
				} else {
					i = GenBlockNumber(min, max)
				}

			}

		case 5:
			oldNumbers = GetLatestNumberFromFile()
			//if LatestExist {
			a, b, c := CreateAndSplitBlockFor3(oldNumbers)

			fmt.Println(a, "\n", b, "\n", c)
			// } else {
			// 	fmt.Print("Please input the latest Numbers first.")
			// 	time.Sleep(500 * time.Millisecond)
			// 	fmt.Print(".")
			// 	time.Sleep(500 * time.Millisecond)
			// 	fmt.Print(".\n")
			// 	continue
			// }

		case 6:
			return
		default:
			fmt.Println("Error input")
			break

		}
		//var oldnumber [6]int

	}

}
