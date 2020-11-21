# 思路
一. 执行命令

1. 检查是否连接，如果没有连接则进行连接。
2. 连接建立后，先将管道中的命令发送到该连接，然后执行当前命令。
3. 如果当前命令执行失败，如果是超时导致，则重连。重连后，重复步骤 2。
4. 返回执行结果。


二. 执行管道中的命令
1. 检查是否连接，如果没有连接则进行连接。
2. 先将管道中的命令发送到该连接。
3. 刷新连接。
4. 在该连接上获取接收结果。如果获取超时，则重连。重连后，重复步骤 2、3。
5. 返回执行结果。


三. 事务

1. 开始事务。
2. 向管道发送命令。
3. 进行步骤一。



四. 重连

- 将当前连接释放掉，然后重新连接。

# 参考
- https://www.jianshu.com/p/fb498f30dff2
- https://github.com/goinbox/redis
