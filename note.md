UDPMux
```txt
                                       data            c.WriteToDst
Client connection        +-------------+  +-------------+  +-------------+    Send to dst Portal
          ---------------> net.UDPConn +--> FakeUDPConn +--> Portal.Conn +--------->
                         +-----------+-+  +-^-----------+  +-------------+
                                     |  Get |
                         +-----------v------+----------------------------+
                         |         connMap   map[tag]FakeUDPConn         |
                         |  tag = addr.Port          tag = data[0:2]     |
                         +----------------------------+------^-----------+
                                c.WriteToSrc          | Get  |
                                          +-----------v-+    |                UDPMux.ReadFromPortal
          <-------------------------------+ FakeUDPConn <----+---------------------
                                          +-------------+
```