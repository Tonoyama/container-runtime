package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main(){
	// コマンドライン引数が足りない場合は使用方法を表示
	if len(os.Args) <= 1 {
		Usage()
	}
	// 引数に応じて関数を実行
	switch os.Args[1] {
	case "run":
		run()
	case "initContainer":
		// 通常は外部から直接呼び出されることはない
		initContainer()
	default:
		// 引数が不明な場合は使用方法を表示
		Usage()
	}
}

// コンテナを実行する
func run(){
	fmt.Printf("Running %v as user id %d in process %d\n", os.Args[2:], os.Getuid(), os.Getpid())

    // cgroupを設定してリソース制限を適用
    mainCgroup()

    // Dockerイメージの情報を取得
    // mainImage()

	// 現在の実行ファイル自身を再実行することで、コンテナ環境を初期化
	cmd := exec.Command("/proc/self/exe", append([]string{"initContainer"}, os.Args[2:]...)...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    // 名前空間の分離とUID/GIDのマッピング設定
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID,
        UidMappings: []syscall.SysProcIDMap{
            {
                ContainerID: 0,
                HostID: 1000,
                Size: 1,
            },
        },
        GidMappings: []syscall.SysProcIDMap{
            {
                ContainerID: 0,
                HostID: 1000,
                Size: 1,
            },
        },
    }
    must(cmd.Run())
}

// 使用方法を表示する
func Usage(){
	fmt.Printf("Usage: run <command>\n For example run /bin/bash")
	os.Exit(1)
}

// コンテナ環境を初期化する
func initContainer() {
    fmt.Printf("Running %v as user id %d in process %d\n", os.Args[2:], os.Getuid(), os.Getpid())

    // ルートディレクトリを変更し、プロセスを隔離
    must(syscall.Chroot("/"))
    must(os.Chdir("/"))
    // procファイルシステムをマウントしてプロセス情報を隔離
    must(syscall.Mount("proc", "proc", "proc", 0, ""))

    // 指定されたコマンドを実行
    cmd := exec.Command(os.Args[2], os.Args[3:]...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    must(cmd.Run())
    // コマンド実行後、procファイルシステムのアンマウント
    must(syscall.Unmount("proc", 0))
}

// エラーチェックを行い、エラーがあればpanicする
func must(err error) {
    if err != nil {
        panic(err)
    }
}
