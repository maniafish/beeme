# v1.0.3
2018.07.31

## [function]

* 新增文件流上传，保存功能

## [fix]

* no bug fix

## [conf]

* 新增FileDir = ./data/ # 上传文件路径, 必须以/结尾, ./开头的则为相对路径

## [db]

* no database change
# v1.0.2
2018.07.30

## [function]

* 最小化程序功能
* 新增部分公共模块
* 更新mylog模块

## [fix]

* no bug fix

## [conf]

* no config change

## [db]

* no database change

# v1.0.1
2017.11.13

## [function]

* 统一使用beego内置config
* 添加错误日志打印

## [fix]

* 修复错误返回码没有及时return的问题

## [conf]

* config.toml -> app.conf.[apps]

# v1.0.0
2017.11.10

## [function]

* 基本框架
* 添加robot功能
* 添加mysql orm
* 集成travis ci单元测试

## [conf]

* TulingURL
* TulingKeys
* DBMaxOpenConns
* DBMaxIdleConns
* UserDB
