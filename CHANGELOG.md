# v1.1.1
2017.12.20

## [function]

* 消息忽略大小写匹配
* 每日消息上限回复
* 数据库表默认charset-utf8 

# v1.1.0
2017.12.19

## [function]

* 新增微信消息回复接口(POST) 
* 新增关注信息回复
* 添加默认controller
* 新增数据库存储特定消息回复

## [db]

* 新增robot_msg表

# v1.0.5
2017.12.1

## [fix]

* 添加AppID作为验签字符串

# v1.0.4
2017.12.1

## [function]

* 修改js接口地址

# v1.0.3
2017.12.1

## [function]

* 自动化部署
* 新增微信公众号token验证

## [conf]

* AppID
* AESKey

## [fix]

* 修复v1.0.2的自动化部署问题

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
