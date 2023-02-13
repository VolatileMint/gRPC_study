# golang_api

## 概要
protoを利用したgolangでのAPIサーバー作成

学習教材
https://zenn.dev/hsaki/books/golang-grpc-starting/viewer/intro

## protocol Buffers導入
[protocol Buffers導入](https://qiita.com/nozmiz/items/fdbd052c19dad28ab067)

プラグインのインストール補足
```
> go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```
確認
```
> protoc-gen-go --version
protoc-gen-go.exe v1.28.1
```
protoファイルのコンパイル
```
> protoc --go_out=../pkg/grpc --go_opt=paths=source_relative --go-grpc_out=../pkg/grpc --go-grpc_opt=paths=source_relative hello.proto
```



AAA.pb.proto ファイルが生成される
## 参考資料
Protocol Buffers導入: https://qiita.com/nozmiz/items/fdbd052c19dad28ab067