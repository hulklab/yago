## 已知问题及解决方案

1.  unknown import path "github.com/ugorji/go/codec": ambiguous import: found github.com/ugorji/go/codec in multiple modules 模块冲突问题

	原因参考：
	
	> https://cloud.tencent.com/developer/article/1417112
	
	解决方案：
	在go.mod文件最下面添加如下代码
	```go
	replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
	```
