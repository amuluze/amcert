# amcert

简单、轻量、自动化的 SSL 证书管理工具。

## usage

```bash
# 常用命令
$ amcert install  # 安装
$ amcert start    # 启动
$ amcert stop     # 停止
$ amcert status   # 查看状态
$ amcert remove   # 移除

# 生成 amcert 配置
$ amcert setup

# 生成 SSL/TLS 证书
$ amcert generate ssl

# 查看 SSL 证书相关信息
$ amcert ssl keys                     # 查看所有证书
$ amcert ssl detail cert-amuluze.com  # 查看指定证书详情
$ amcert ssl expire cert-amuluze.com  # 查看指定证书过期时间
```

