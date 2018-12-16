package mapreduce

import (
	"fmt"
	"sync"
)

// schedule starts and waits for all tasks in the given phase (Map or Reduce).
func (mr *Master) schedule(phase jobPhase) {
	var ntasks int
	var nios int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mr.files)                         //war_peace.txt pg_some.txt
		nios = mr.nReduce                              //每一个文件都将本map成nios份
	case reducePhase:
		ntasks = mr.nReduce
		nios = len(mr.files)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)
	var wg sync.WaitGroup
	for i,file := range mr.files{
		//freeWorker := <- mr.registerChannel
		var taskArgs DoTaskArgs
		taskArgs.File = file
		taskArgs.JobName = mr.jobName
		taskArgs.NumOtherPhase = nios        //这个赋值可能有问题,bug
		taskArgs.Phase = phase
		taskArgs.TaskNumber = i
		//ok := call(freeWorker, "Worker.DoTask", taskArgs, new(struct{}));if ok == false{
		//	log.Printf("call Dotask err ")
		//}
		wg.Add(1)
		go mr.sendTaskAndWait(&taskArgs,&wg)
	}
	wg.Wait()
	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO
	//

	fmt.Printf("Schedule: %v phase done\n", phase)
}

func (mr *Master)sendTaskAndWait(taskArgs *DoTaskArgs,wg *sync.WaitGroup){
	defer wg.Done()
	freeWorker := <- mr.registerChannel
	for {
		ok := call(freeWorker, "Worker.DoTask", *taskArgs, new(struct{}));
		if ok {
			go func() {
				mr.registerChannel <- freeWorker
			}()
			break
		}
	}
}
