package main

import (
	"C"
	"unsafe"
)

//export Stop
func Stop() {
	mu.Lock()
	defer mu.Unlock()
	if cancelFunc != nil {
		cancelFunc()
	}
}

//export RunShared
func RunShared(argc C.int, argv **C.char) {
	// 将 C 的 argv 转换为 Go 的 string slice
	length := int(argc)
	tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(argv))[:length:length]
	args := make([]string, length)
	for i, s := range tmpslice {
		args[i] = C.GoString(s)
	}

	// 如果传递了参数，且第一个参数是程序名（通常是），cobra 需要剩下的参数
	// 但 cobra 的 SetArgs 会完全覆盖 os.Args[1:]。
	// 如果 args 是 ["mcp-liner", "-v"]，SetArgs 应该只接受 ["-v"]?
	// cobra.Command.SetArgs: "sets arguments for the command. It is set to os.Args[1:] by default."
	// 所以如果 args 包括程序名，我们应该传 args[1:]
	if len(args) > 0 {
		rootCmd.SetArgs(args[1:])
	}

	main()
}

// 确保 main_cgo.go 仅在 CGO 启用时参与构建（如果不显示指定 build tags，只要 import "C" 就会自动处理）
// 但为了明确，我们保留它与 main.go 在同一个包中，共享全局变量。
// 注意：main.go 中的全局变量 appVersion, cancelFunc, mu, stdinReader 必须是可见的。
