用于自动备份的脚本，不能直接使用

# nga_backup_create.sh

生成 db 文件的备份，修改 DB_FILE 为实际的文件路劲

# nga_bak_s3.service

具体调用的s3进行备份的 service

EnvironmentFile 中可以配置 AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY 这些变量

## nga_bak_s3.sh

BACKUP_FILE 修改为到 nga_backup_create.sh 文件的绝对路径

url bucketname 都修改为实际的名字


# nga_bak_s3.timer

定时触发器 可以修改 OnCalendar 来调整触发时间