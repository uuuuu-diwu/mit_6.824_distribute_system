package main

import (
	"../mapreduce"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// The mapping function is called once for each piece of the input.
// In this framework, the key is the name of the file that is being processed,
// and the value is the file's contents. The return value should be a slice of
// key/value pairs, each represented by a mapreduce.KeyValue.
func mapF(document string, value string) (res []mapreduce.KeyValue) {
	// 在wordcount的例子中mapF的功能应该是string中获取到单词（关注下strings.FieldsFunc打用法吧），
	// 返回的结构应该类似KeyValue{w, "1"}
	contentSlice := strings.FieldsFunc(value,unicode.IsSpace)
	for _,word:= range contentSlice{
		res = append(res, mapreduce.KeyValue{word,"1"})
	}
	log.Printf("documetn:%s map done",document)
	return res
}

// The reduce function is called once for each key generated by Map, with a
// list of that key's string value (merged across all inputs). The return value
// should be a single output value for that key.
func reduceF(key string, values []string) string {
	// TODO: you also have to write this function
	// reduceF对每个key调用，然后处理values,在这个例子中，相加全部的１就是单词出现打次数来
	var res int64
	for _,val := range values{
		l,err := strconv.Atoi(val)
		if err!= nil{
			log.Fatalf("read kv error ,key = %s,val=%s,err=%v",key,val,err)
			continue
		}
		res += int64(l)
	}
	return strconv.FormatInt(res,10)
}

// Can be run in 3 ways:
// 1) Sequential (e.g., go run wc.go master sequential x1.txt .. xN.txt)
// 2) Master (e.g., go run wc.go master localhost:7777 x1.txt .. xN.txt)
// 3) Worker (e.g., go run wc.go worker localhost:7777 localhost:7778 &)
func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	if len(os.Args) < 4 {
		fmt.Printf("%s: see usage comments in file\n", os.Args[0])
	} else if os.Args[1] == "master" {
		var mr *mapreduce.Master
		if os.Args[2] == "sequential" {
			mr = mapreduce.Sequential("wcseq", os.Args[3:], 3, mapF, reduceF)
		} else {
			mr = mapreduce.Distributed("wcseq", os.Args[3:], 3, os.Args[2])
		}
		mr.Wait()
	} else {
		mapreduce.RunWorker(os.Args[2], os.Args[3], mapF, reduceF, 100)
	}
}
