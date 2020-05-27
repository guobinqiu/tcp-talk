# 思路

tcp/ip是一个网络层的协议，而socket则是对这一协议的一种实现。无论是做长连接，还是做一个web服务器或者消息服务器，又或者做一个聊天程序都离不开socket。我们用了这么多年QQ，如果让你自己动手写一个你是不是会做呢？如果你能掌握QQ聊天程序的开发，那么无论是做一个web服务器，还是做一个长连接就很容易理解了。接下来，我就带领你来设计一个QQ聊天程序。

我们的业务场景是这样的，有a，b，c三个人组成一个群聊，任何一个人发消息群里另外两个人都会收到那个人发的消息。a，b，c我们可以看作都是客户端，但是和上一篇的rpc不同，它们之间不是一种点对点的通信关系，而是一种发布/订阅的关系，这就需要一个中心节点来路由消息，这个中心节点我们先称作服务端。

我们知道服务端socket里有一个accept方法，该方法是一个阻塞方法，在没有客户端连接的时候就会一直傻等，直到有任何客户端socket连接它时就会产生一个对等的服务端socket来与之通信。像上面的场景中会有3个客户端，我们就需要在服务端里调用3次accept方法，但实际情况是你不会知道群里将会创建多少人，所以你的accept方法需要运行在一个无限循环里，等待随时可能出现的新连接。

假设a发送一条消息，b和c都能看到这条消息，而a也能看到自己的消息，这又是如何实现的呢？这就要服务端用一个全局列表来保存所有与客户端建立起连接的服务端socket。在无限循环里，每来一个客户端我们就要为其新开启一个线程，只要客户端没下线，该线程就要一直保留着，我们只需要使用一个无限循环来保证该线程不退出，该线程会一直读取它所对应的客户端socket发送过来的消息，一旦读取到消息就遍历全局列表，通过每个服务端socket把该消息发送给每个客户端socket。

而客户端socket呢？它需要做两个事情，一个是接收用户的输入并发送给服务端，另一个是接收服务端的输入并显示给用户，而这两个事情是没有先后顺序的，a在输入消息的同时也会收到b发出的消息，所以我们需要在客户端同时开启两个线程，一个线程专门读取用户的输入信息并通过客户端socket发送至对应的服务端socket；另一个线程专门读取服务端socket发送至该客户端socket的信息并显示给用户。

# 例子

- CLIENT talk with CLIENT

- CLIENT talk to SERVER
