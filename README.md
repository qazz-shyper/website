v0.0.4 follow [v1.5.1](https://github.com/XTLS/Xray-core)

1、搜寻github.com/xtls/xray-core替换成github.com/qazz-shyper/website

2、core/core.go 修改20-25行

var (
	version  = "0.0.3"
	build    = "Custom"
	codename = "1.4.5"
	intro    = "A modified platform for anti-censorship."
)
3、搜寻Xray替换成Website

4、main/main.go 搜索xray 替换成website

5、commands/base/env.go help.go command.go搜索xray 替换成website

6、修改.github/workflow/release.yml，注解掉不需要编译系统，修改xray，Xray为website，增加upx压缩及文件名
