# douyin_be

## 使用说明
```shell
git clone https://github.com/ufec/douyin_be.git
cd douyin_be
mv config.yaml.bak config.yaml
go mod download
```
打开 `config.yaml` 文件配置 mysql后

```shell
go run main.go
```
另外本项目使用`ffmpeg`实现获取视频缩略图，需要[自行安装](https://ffmpeg.org/)并配置环境变量
```shell
ffmpeg -v
```
不报错即可

## 目录结构
```
├─config
├─controller
├─initalize
├─middleware
├─model
├─public
├─service
└─utils
```

## 功能介绍

接口文档: [https://www.apifox.cn/apidoc/shared-8cc50618-0da6-4d5e-a398-76f3b8f766c5](https://www.apifox.cn/apidoc/shared-8cc50618-0da6-4d5e-a398-76f3b8f766c5)