## SQL To Pb 数据库到PB转换工具
> 一个可以快速的将数据库转换为PB文件的工具，使用灵活，可以根据模板高度定制

本工具旨在不用改动代码的情况下，完成sql到文件的迁移，所以提供了简单高效可以定制化的yml文件以及高度定制的proto.tpl模板。帮助你快速生成grpc-message文件。

## 使用方法
你可以看到配置文件中，简单的文件配置
```
config:
  in_tpl: proto.tpl
  out_dir: data/message
  out_pb: data/pb;v1;
db:
  db_name: data
  db_tables: []
  dsn: root:123456@tcp(127.0.0.1:3306)/data?charset=utf8mb4&parseTime=true  
```

### 数据源配置
1. 首先是DB连接，配置一下DSN `db.dsn`
    > 单独的配置`db_name` 是为了更好的获取到数据库schema信息
2. 配置了需要导出那些表，如果`db_tables`为空，默认是导出所有

### 数据输出配置
- `config.in_tpl` 指定了需要生成的pb文件模板。
    > 这使得我们可以灵活的定制生成pb文件的内容,针对不同的项目，可以使用不同的模板
- `config.out_dir` 指定了输出pb文件的路径
- `config.out_pb` 指定了protoc生成后的路径

项目参考
- https://github.com/guyan0319/mysql-to-proto
- https://github.com/zhangkai8048/mysql-to-proto