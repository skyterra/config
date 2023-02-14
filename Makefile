.PHONY: test cover cover-func cover-html clean bench

# 执行单元测试
test:
	go test ./...

# 执行benchmark测试
bench:
	go test -bench . ./... -run=none

# 统计覆盖率
cover:
	go test ./... -coverprofile cover.profile

# 打开浏览器显示覆盖统计信息
cover-html:
	go tool cover -html=cover.profile

# 显示函数覆盖统计
cover-func:
	go tool cover -func=cover.profile

# 清理
clean:
	rm cover.profile