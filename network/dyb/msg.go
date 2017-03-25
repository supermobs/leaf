package dyb

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"runtime/debug"

	"github.com/name5566/leaf/log"
)

func READ_int8(bin []byte) (*int8, []byte) {
	data := int8(bin[0])
	return &data, bin[1:]
}
func READ_int16(bin []byte) (*int16, []byte) {
	data := int16(binary.LittleEndian.Uint16(bin[:2]))
	return &data, bin[2:]
}
func READ_int32(bin []byte) (*int32, []byte) {
	data := int32(binary.LittleEndian.Uint32(bin[:4]))
	return &data, bin[4:]
}
func READ_int64(bin []byte) (*int64, []byte) {
	data := int64(binary.LittleEndian.Uint64(bin[:8]))
	return &data, bin[8:]
}
func READ_string(bin []byte) (*string, []byte) {
	length, bin := READ_int16(bin)
	data := ""
	if *length > 1 {
		data = string(bin[:*length])
	}
	return &data, bin[*length:]
}
func WRITE_int8(data *int8) []byte {
	bin := make([]byte, 1)
	bin[0] = byte(*data)
	return bin
}
func WRITE_int16(data *int16) []byte {
	bin := make([]byte, 2)
	binary.LittleEndian.PutUint16(bin, uint16(*data))
	return bin
}
func WRITE_int32(data *int32) []byte {
	bin := make([]byte, 4)
	binary.LittleEndian.PutUint32(bin, uint32(*data))
	return bin
}
func WRITE_int64(data *int64) []byte {
	bin := make([]byte, 8)
	binary.LittleEndian.PutUint64(bin, uint64(*data))
	return bin
}
func WRITE_string(data *string) []byte {
	length := int16(len(*data))
	return append(WRITE_int16(&length), []byte(*data)...)
}

type BaseProcessor struct {
	call_back func(interface{}, interface{})
	is_aes    bool
}

func NewBaseProcessor(call_back func(interface{}, interface{}), is_aes bool) *BaseProcessor {
	p := new(BaseProcessor)
	p.call_back = call_back
	p.is_aes = is_aes
	return p
}

func (p *BaseProcessor) Route(msg interface{}, userData interface{}) error {
	//这里捕获下异常吧,万一操作错误,免得玩家断线啥的
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			stack := debug.Stack()
			log.Error(string(stack))
		}
	}()
	if p.is_aes {
		msg = msg_aes_decrypt(msg.([]byte))
	}
	p.call_back(msg, userData)
	return nil
}

// goroutine safe
func (p *BaseProcessor) Unmarshal(data []byte) (interface{}, error) {
	return data, nil
}

// goroutine safe
func (p *BaseProcessor) Marshal(msg interface{}) ([][]byte, error) {
	if p.is_aes {
		data := msg_aes_encrypt(msg.([]byte))
		return [][]byte{data}, nil
	}
	return [][]byte{msg.([]byte)}, nil
}

func msg_aes_encrypt(bin []byte) []byte {
	length := len(bin)
	for i := 0; i < 16-length%16; i++ {
		bin = append(bin, 0)
	}
	cmdbin := bin[:4]
	bin = bin[4:]
	return aesEncrypt(append(bin, cmdbin...))
}
func msg_aes_decrypt(bin []byte) []byte {
	bin = aesDecrypt(bin)
	length := len(bin)
	body, cmd := bin[:length-4], bin[length-4:]
	return append(cmd, body...)
}

var g_aeskey []byte = []byte("jin_tian_ni_chi_le_mei_you?chi_l")

func aesEncrypt(origData []byte) []byte {
	if len(origData)%16 > 0 {
		fmt.Println("AesEncrypt error:", len(origData), "%16 != 0")
		return origData
	}
	block, err := aes.NewCipher(g_aeskey[:16])
	if err != nil {
		return nil
	}
	//blockSize := block.BlockSize()
	blockMode := cipher.NewCBCEncrypter(block, g_aeskey[16:])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted
}
func aesDecrypt(crypted []byte) []byte {
	if len(crypted)%16 > 0 {
		fmt.Println("AesDecrypt error:", len(crypted), "%16 != 0")
		return crypted
	}
	block, err := aes.NewCipher(g_aeskey[:16])
	if err != nil {
		return nil
	}
	//blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, g_aeskey[16:])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	return origData
}
