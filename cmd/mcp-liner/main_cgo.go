package main

import "C"

//export Stop
func Stop() {
	mu.Lock()
	defer mu.Unlock()
	if cancelFunc != nil {
		cancelFunc()
	}
	// 关闭 Pipe Reader 以打断阻塞的 Read 操作
	// 这不会关闭真实的系统 Stdin
	if stdinReader != nil {
		_ = stdinReader.Close()
	}
}

//export RunShared
func RunShared() {
	// 确保在 CGO 环境下的初始化，如果需要的话
	// 这里直接调用 main 中的 run 逻辑或者复用 main 本身逻辑
	// 由于 main() 是入口，直接调用 main()
	main()
}

// 确保 main_cgo.go 仅在 CGO 启用时参与构建（如果不显示指定 build tags，只要 import "C" 就会自动处理）
// 但为了明确，我们保留它与 main.go 在同一个包中，共享全局变量。
// 注意：main.go 中的全局变量 appVersion, cancelFunc, mu, stdinReader 必须是可见的。
