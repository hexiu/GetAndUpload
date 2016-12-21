package main

 

import (

    "code.google.com/p/mahonia"

    "fmt"

)

 

func main() {

    //"你好，世界！"的GBK编码

    testBytes := []byte{0xC4, 0xE3, 0xBA, 0xC3, 0xA3, 0xAC, 0xCA, 0xC0, 0xBD, 0xE7, 0xA3, 0xA1}

    var testStr string

    utfStr := "你好，世界！"

    var dec mahonia.Decoder

    var enc mahonia.Encoder

 

    testStr = string(testBytes)

 

    dec = mahonia.NewDecoder("gbk")

    if ret, ok := dec.ConvertStringOK(testStr); ok {

        fmt.Println("GBK to UTF-8: ", ret, " bytes:", testBytes)

    }

 

    enc = mahonia.NewEncoder("gbk")

    if ret, ok := enc.ConvertStringOK(utfStr); ok {

        fmt.Println("UTF-8 to GBK: ", ret, " bytes: ", []byte(ret))

    }

    return

}