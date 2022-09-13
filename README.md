# tianxun-lite



## **参数**



```bash
Usage of ./tianxun-lite:
  # 自动播报的定时任务 配合add_cron使用
  -cron string
    	cron expr (default "* */1 * * * sh tcs-tools.sh check -mod all")
  # 模式默认终端大盘, add_cron增加自动播报, del_cron删除自动播放
  -m string
    	default/get_data/add_cron/del_cron (default "default")

```



## **示例**



```bash
# 打开终端大盘
./tianxun-lite
# 增加自动播报
./tianxun-lite -m add_cron -cron "*/5 * * * * sh tcs-tools.sh check -mod docker"
# 删除自动播报
./tianxun-lite -m del_cron
```
