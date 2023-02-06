## KubeVela Cue DEMO

```shell
go build -o cue-gen cmd/main.go && ./cue-gen

cat test/struct.cue
```

目前进度：读取 go 文件，转换所有结构体并展开

## TODO

- 部分指针结构体如 `http.Request` 展开无限循环，
