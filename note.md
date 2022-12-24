# note
https://asciiflow.com/
## TODO

UDPMux

- ~~add tag to packet between portals. ~~
- remove conns which are timeout.

## Portal(Client)



## UDPMux

```go
//                                        data            c.WriteToDst
// Client connection        +-------------+  +-------------+  +-------------+    Send to dst Portal
//           ---------------> net.UDPConn +--> FakeUDPConn +--> Portal.Conn +--------->
//                          +-----------+-+  +-^-----------+  +-------------+
//                                      |  Get |
//                          +-----------v------+----------------------------+
//                          |         connMap   map[tag]FakeUDPConn         |
//                          |  tag = addr.Port          tag = data[0:2]     |
//                          +----------------------------+------^-----------+
//                                 c.WriteToSrc          | Get  |
//                                           +-----------v-+    |                UDPMux.ReadFromPortal
//           <-------------------------------+ FakeUDPConn <----+---------------------
//                                           +-------------+
```

## iranai

```go
// 
//                   udpMux
// 
//                  +---------------------------------------------------------+
//                  |                                                         |
//                  |                                                         |
//          Client  |     +-----------+                                       |
//          Connection    |net.UDPConn|                                       |
//            ------+---->|listening  |                                       |
//                  |     +--+--------+                                       |
//                  |        |   v Data                                       |
//                  |        | +----------------+      +-----------------+    |
//                  |        | | FakeUDPConn    |      | Portal.Conn     |    |  Send to destination Portal
//                  |        | |                +----->|                 +----+----->
//                  |        | +----------------+      +-----------------+    |
//                  |   Tag  |  ^                                             |
//                  |        |  | GetOrDefault                                |
//                  |        v  |                                             |
//                  |    +------+--------                                     |  UDPMux.ReadFromPortal
//                  |    | connMap        |                                   |
//                  |    | get FakeUDPConn|                     <-------------+---------
//                  |    |                |                                   |
//                  |    +----------------+                                   |
//                  |                                                         |
//                  +---------------------------------------------------------+
// 
// 
// 
//                                        data            c.WriteToDst
// Client connection        +-------------+  +-------------+  +-------------+    Send to dst Portal
//           ---------------> net.UDPConn +--> FakeUDPConn +--> Portal.Conn +--------->
//                          +------------++  +-^-----------+  +-------------+
//                                       | Get |
//                          +------------v-----+----------------------------+
//                          |         connMap   map[tag]FakeUDPConn         |
//                          |  tag = addr.Port          tag = data[0:2]     |
//                          +----------------------------+------^-----------+
//                                 c.WriteToSrc          | Get  |
//                                           +-----------v-+    |                UDPMux.ReadFromPortal
//           <-------------------------------+ FakeUDPConn <----+---------------------
//                                           +-------------+
// 

```