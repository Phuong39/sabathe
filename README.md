# sabathe
一个服务端-客户端-植入物的三端式RAT，不完整，仅仅作为学习开源（仅可以命令执行）

服务端：采用rpcx与客户端通信，并接收客户端传输过来的命令
客户端：采用rpcx与服务端通信
植入物：kcp协议接收命令，http协议回传命令执行结果
服务端：
![image](https://user-images.githubusercontent.com/74412075/149697357-3b60019d-78dc-4bbb-9b8d-978f739a21c4.png)
植入物：
![image](https://user-images.githubusercontent.com/74412075/149697523-c80b209f-7175-48c2-8287-74fd4c531991.png)
客户端
![image](https://user-images.githubusercontent.com/74412075/149698174-e4e47b09-87e5-475e-8d28-d415f504e706.png)
后续不在开源，需要请自行修改代码
