package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

/**
cgroupは"Control Group"の略。
プロセスIDをグループに追加すると、共通の設定を適用できる。
ホストOSが持つCPUやメモリなどのリソースの制限をかけることができる。
**/

var pid int = os.Getpid()
var cgroup_dir string = "/sys/fs/cgroup/"
var cgroup_mem_dir string
var cgroup_mem_file string

var cgroup_cpu_dir string
var cgroup_cpu_file string


func mainCgroup() {
	// 使用メモリー数(定数)
	var mem_limit int
	var cpu_limit_num int

	mem_limit = 10
	// メモリ制限を設定する
	memory_limit(mem_limit)

	cpu_limit_num = 10
	cpu_limit(cpu_limit_num)
}

// メモリを制限する関数
func memory_limit(mem_limit int) {
	// cgroupのメモリのディレクトリを定義する
	cgroup_mem_dir = cgroup_dir + "memory/" + strconv.Itoa(pid)
	// 設定ファイルを定義する
	cgroup_mem_file = cgroup_mem_dir + "/memory.limit_in_bytes"
	// ディレクトリが存在する場合、メモリ制限を設定する
	if err := os.MkdirAll(cgroup_mem_dir, 0700); err != nil {
		fmt.Println(err)
	}
	// ファイルにメモリ制限を設定する
	ioutil.WriteFile(cgroup_mem_file, []byte(strconv.Itoa(mem_limit*1024*1024)), 0700)
	// プロセスIDを追加する
	ioutil.WriteFile(cgroup_mem_dir + "/cgroup.procs", []byte(strconv.Itoa(pid)), 0700)
}

// CPUを制限する関数
func cpu_limit(cpu_limit_num int) {
	// cgroupのCPUのディレクトリを定義する
	cgroup_cpu_dir = cgroup_dir + "cpu/" + strconv.Itoa(pid)
	// 設定ファイルを定義する
	cgroup_cpu_file = cgroup_cpu_dir + "/cpu.cfs_quota_us"
	// ディレクトリが存在する場合、CPU制限を設定する
	if err := os.MkdirAll(cgroup_cpu_dir, 0700); err != nil {
		fmt.Println(err)
	}
	// ファイルにCPU制限を設定する
	ioutil.WriteFile(cgroup_cpu_file, []byte(strconv.Itoa(cpu_limit_num * 1000)), 0700)
	// プロセスIDを追加する
	ioutil.WriteFile(cgroup_cpu_dir + "/cgroup.procs", []byte(strconv.Itoa(pid)), 0700)
}