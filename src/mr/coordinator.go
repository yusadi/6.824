/*
 * @Description:
 * @User: Snaper <532990528@qq.com>
 * @Date: 2021-06-16 12:25:17
 * @LastEditTime: 2021-06-18 00:43:25
 */

package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

//程序master，协调器，负责分发委派任务
type Coordinator struct {
	nFile               int             //要统计的文本数量
	completedMapTask    int             //完成的Map任务数量
	completedReduceTask int             //完成的reduce任务数量
	QMapTask            chan MapTask    //保存Map任务，即文件路径，因为存在并发所以使用chan保存
	QReduceTask         chan ReduceTask //保存reduce任务，即文件路径，因为存在并发所以使用chan保存
	MapOutputFile       []string        //map任务输出的文件
	ReduceOutputFile    []string        //reduce任务输出的文件
	Lock                sync.Mutex      //互斥锁，保证协调器在工作时的线程安全性
}

//map任务类
type MapTask struct {
	TaskSeqNum int    //任务序号
	Filename   string //任务路径
}

//reduce任务类
type ReduceTask struct {
	TaskSeqNum int    //任务序号
	Filename   string //任务路径
}

/**
 * @name: SendTask
 * @desc: 请求任务
 * @param {*MrRpcArgs} args
 * @param {*MrRpcReply} reply
 * @return {*}
 */
func (c *Coordinator) SendTask(args *MrRpcArgs, reply *MrRpcReply) error {

	return nil
}

/**
 * @name: ReportTask
 * @desc: 任务完成后，汇报任务
 * @param {*}
 * @return {*}
 */
func (c *Coordinator) CompleteTask(args *MrRpcArgs, reply *MrRpcReply) error {
	c.Lock.Lock()
	switch args.TaskType {
	case MAP_TASK:
		c.completedMapTask++
		c.MapOutputFile = append(c.MapOutputFile, args.FilePaths...)
		if c.completedMapTask == c.nFile {
			reply.TaskType = REDUCE_TASK

		}
	case REDUCE_TASK:

	}
	c.Lock.Unlock()
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.

	return ret
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}

	// Your code here.

	c.server()
	return &c
}