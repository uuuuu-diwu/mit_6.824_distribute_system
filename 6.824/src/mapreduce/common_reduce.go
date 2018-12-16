package mapreduce

import (
	"encoding/json"
	"log"
	"os"
)

// doReduce does the job of a reduce worker: it reads the intermediate
// key/value pairs (produced by the map phase) for this task, sorts the
// intermediate key/value pairs by key, calls the user-defined reduce function
// (reduceF) for each key, and writes the output to disk.
func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTaskNumber int, // which reduce task this is
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	// TODO:
	// You will need to write this function.
	// You can find the intermediate file for this reduce task from map task number
	// m using reduceName(jobName, m, reduceTaskNumber).
	// Remember that you've encoded the values in the intermediate files, so you
	// will need to decode them. If you chose to use JSON, you can read out
	// multiple decoded values by creating a decoder, and then repeatedly calling
	// .Decode() on it until Decode() returns an error.
	//
	// You should write the reduced output in as JSON encoded KeyValue
	// objects to a file named mergeName(jobName, reduceTaskNumber). We require
	// you to use JSON here because that is what the merger than combines the
	// output from all the reduce tasks expects. There is nothing "special" about
	// JSON -- it is just the marshalling format we chose to use. It will look
	// something like this:
	//
	// enc := json.NewEncoder(mergeFile)
	// for key in ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()
	interResult := make(map[string][]string)
	mergeFileName := mergeName(jobName,reduceTaskNumber)
	//mutex := sync.Mutex{}
	//var wg sync.WaitGroup
	//wg.Add(nMap)
	interFileList := make([]*os.File,0)
	defer closeAllFile(interFileList)
	for i := 0;i < nMap;i++{
		log.Printf("process %d-%d",i,reduceTaskNumber)
		fileName := reduceName(jobName,i,reduceTaskNumber)
		intermediateFile ,err := os.Open(fileName);if err != nil{
			//todo
		}
		interFileList = append(interFileList,intermediateFile)
		dec := json.NewDecoder(intermediateFile)
		/*go*/ subDecode(dec,interResult/*,&mutex,wg*/)
		log.Printf("after decode %s,map size = %d",fileName,len(interResult))
	}
	//wg.Wait()
	//wait go work done
	log.Printf("process reduce %d done",reduceTaskNumber)
	mergeFile,err := os.OpenFile(mergeFileName,os.O_CREATE|os.O_RDWR,0644);if err != nil{
		log.Fatalf("open file %s failed with msg: %v",mergeFileName,err)
		return
	}
	enc := json.NewEncoder(mergeFile)
	for key,val := range interResult{
		thisCount := reduceF(key,val)
		enc.Encode(KeyValue{key,thisCount})
	}
	mergeFile.Close()
}

func subDecode(deCodeFile *json.Decoder,interResult map[string][]string,/*mutex *sync.Mutex,wg sync.WaitGroup*/){
	log.Printf("now map size = %d",len(interResult))
	for{
		var tmp KeyValue
		if err := deCodeFile.Decode(&tmp);err != nil{
			break
		}
		//mutex.Lock()
		interResult[tmp.Key] = append(interResult[tmp.Key], tmp.Value)
		//mutex.Unlock()
	}
	//wg.Done()
}

func closeAllFile(fileList []*os.File){
	for _,file:=range fileList{
		file.Close()
	}
}