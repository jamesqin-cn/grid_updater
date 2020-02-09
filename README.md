# grid_updater

网格更新器，解析basefile，逐行逐列更新MySQL

## 设计目标
提供一种数据从异构容器快速迁移到MySQL的工具，关键点：
- 能很方便接入异构容器，解决方案是basefile文件格式
- 能做到快速局部更新MySQL，解决方案是grid_update算法，把basefile转化为逐行逐列的MySQL更新操作

## 框架
定义Reader和Writer，数据从Reader转移到Writer

## 原理
### basefile文件格式
基本逻辑
- 文件名包含的入库的基本信息，文件名格式：{dbName}#{tbName}#{statDate}#{pkNum}#{eraseBeforeUpdate}#{desc}.{ext}
- 文件内容以明文表示
- 指令行：从第一行开始，以#或//开头的，是为设置行，注释符为#或//，后面是key=val的格式
- 字段行：文件第一行非指令行为入库目标的字段，每一列按指定分隔符隔开，一列就是一个字段，字段名语法：[!]{colName}:{colType}，!表示该字段需要建立索引
- 数据行：紧随字段行至末尾，一行为一条入库记录，一列为一个字段

```
base file format:
+--------------------------+
|          指令行           |
+--------------------------+
|          字段行           |
+--------------------------+
|                          |
|                          |
|                          |
|          数据行           |
|                          |
|                          |
|                          |
+--------------------------+


```


### basefile文件示例

- 文件名

test#x_user_summary#20171214#1#0#fund_withdraw.txt

- 内容

```
#  before_command = UPDATE `tzj_user_summary` SET `fund_last_withdraw_time` = NOW() WHERE `tzj_user_summary`.`tzj_username` = 'u008';
// after_command  = UPDATE `tzj_user_summary` SET `fund_withdraw_count` = '8848' WHERE `tzj_user_summary`.`tzj_username` = 'u008';
tzj_username	fund_first_withdraw_time:DATETIME	!fund_last_withdraw_time:DATETIME	fund_withdraw_count:INT(10)	fund_withdraw_amount:DECIMAL(15,2)
u001	2017-04-07 14:18:34	2017-12-01 09:14:15	10	201494.57
u002	2015-06-09 11:04:18	2016-04-27 16:18:15	6	61540.24
u003	2016-01-25 15:34:46	2016-01-25 15:34:46	1	6010.70
u004	2015-08-05 10:55:23	2017-05-04 15:57:10	4	7140.65
u005	2015-09-05 00:55:37	2017-09-16 19:53:53	2	12162.82
u006	2015-11-14 12:31:31	2015-11-26 08:27:03	2	385.26
u007	2016-07-01 13:14:44	2016-07-25 12:38:17	2	16202.57
u008	2017-02-21 08:49:54	2017-02-21 08:49:54	1	61261.60
u009	2015-05-14 09:15:30	2016-09-11 08:20:09	22	856837.80
u010	2017-02-15 23:54:20	2017-04-18 08:05:44	3	6712.82
```

### 如何生成basefile

```
mysql -hxxx -Pxxx -uxxx -pxxx -e "<your sql here>" > your_base_file.txt
```


### grid_update工作逻辑
算法步骤
- 如果database不存在，则自动创建
- 如果table不存在，则自动创建
- 如果column不存在，则自动创建
- 如果index不存在，则自动创建
- 从basefile中的指令行中提取 before_command 的赋值内容，作为SQL执行，该功能一般用于预置处理
- 把records更新插入到MySQL指定表
- 从basefile中的指令行中提取 after_command 的赋值内容，作为SQL执行，该功能一般用于字段的关联更新

其中，更新插入是使用`INSERT INTO ... VALUES(...) ON DUPLICATE KEY UPDATE ...` 的方法

### 创建table的逻辑
- 创建一个只包含primary key、_created_at、_updated_at、_deleted_at这些字段的表，其中primary key要结合basefile文件名的key_col_cnt字段及basefile文件的column行解析得到
- 扫描basefile的column行，通过ALTER...ADD COLUMN语法完善所建之表
- 扫描basefile的column行，通过ALTER...ADD INDEX语法完善所建之表

