package usecase

import (
	"fmt"
	"main/internal/logger"
	"math"
	"strconv"
	"strings"
)

type UseCase struct {
	logger logger.ILogger
}

func New(logger logger.ILogger) UseCase {
	return UseCase{logger: logger}
}

type IPAddress struct {
	Bit       string
	Count     int
	AddrUnder string
	AddrBroad string
	FirstAddr string
	LastAddr  string
	Error     error
}

func (us UseCase) GetInfoByIP(addr string) (IPAddress, error) {
	var ipa IPAddress
	ipa.Bit = us.bitView(addr)[:35]
	h := strings.Split(addr, "/")
	mask := us.Atoi(h[1])

	ipa.Count = us.countHosts(mask)
	ipa.AddrUnder = us.find(ipa.Bit, mask, "0")
	ipa.AddrBroad = us.find(ipa.Bit, mask, "1")

	lastAddr, err := us.GetLastAddr(ipa.AddrBroad)
	if err != nil {
		return IPAddress{}, err
	}
	ipa.LastAddr = lastAddr

	ipa.FirstAddr, err = us.GetFirstAddr(ipa.AddrUnder)
	if err != nil {
		return IPAddress{}, err
	}

	return ipa, nil
}

func (us UseCase) bitView(addr string) string {
	h := strings.Split(addr, "/")
	arr := strings.Split(h[0], ".")

	var sb strings.Builder
	for _, el := range arr {
		digit, err := strconv.Atoi(el)
		if err != nil {
			us.logger.Warn(err.Error())
		}

		var b = make([]string, 8)
		j := 0
		for i := 128; j < 8; i /= 2 {
			if digit >= i {
				b[j] = "1"
				digit -= i
			} else {
				b[j] = "0"
			}
			j++
		}

		sv := strings.Join(b, "") + "."
		sb.WriteString(sv)
	}

	return sb.String()
}

func (us UseCase) Atoi(s string) int {
	digit, err := strconv.Atoi(s)
	if err != nil {
		us.logger.Warn(err.Error())
	}

	return digit
}

func (us UseCase) countHosts(mask int) int {
	const maxMask = 32
	res := math.Pow(2.0, float64(maxMask-mask))

	return int(res) - 2
}

func (us UseCase) find(bit string, mask int, el string) string {
	bit = bit[:35]
	var addr = make([]string, 35)
	zero := 0
	val := 32 - (32 - mask)
	for i := 0; i < 35; i++ {
		addr[i] = string(bit[i])
		if bit[i] == '.' {
			zero++
			continue
		}

		if i < val+zero {
			continue
		}

		addr[i] = el
	}

	bit = strings.Join(addr, "")

	var sb strings.Builder
	arr := strings.Split(bit, ".")
	for _, el := range arr {
		sb.WriteString(fmt.Sprint(us.fromBinToTen(el)) + ".")
	}

	return strings.Trim(sb.String(), ".")
}

func (us UseCase) fromBinToTen(s string) float64 {
	num := 0.0
	for i, c := range s {
		if c == '0' {
			continue
		}
		num += math.Pow(2.0, float64(8-1-i))
	}

	return num
}

func (us UseCase) GetLastAddr(addr string) (string, error) {
	digitsArr := strings.Split(addr, ".")
	lastDigit := digitsArr[len(digitsArr)-1]
	digit, err := strconv.Atoi(lastDigit)
	if err != nil {
		us.logger.Warn(err.Error())
		return "", err
	}

	digitsArr[len(digitsArr)-1] = strconv.Itoa(digit - 1)
	return strings.Join(digitsArr, "."), nil
}

func (us UseCase) GetFirstAddr(addr string) (string, error) {
	digitsArr := strings.Split(addr, ".")
	lastDigit := digitsArr[len(digitsArr)-1]
	digit, err := strconv.Atoi(lastDigit)
	if err != nil {
		us.logger.Warn(err.Error())
		return "", err
	}

	digitsArr[len(digitsArr)-1] = strconv.Itoa(digit + 1)
	return strings.Join(digitsArr, "."), nil
}
