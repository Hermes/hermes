# Protocol
===========

## DHT

Hermes underpinnings rely on a fault tolerant, and powerful distributed hash table (DHT). The DHT of choice is (for now) the Chord DHT (However ZHT(link?) is also being investigated). Chord DHT is similar to most DHT(link?) in how it stores date (key/value model), however it differs greatly in its routing overlay. Each node has a unique ID, and is responsible for a portion of the entire hash ring(link). Messages are routed to nodes in O(logN) time, with the MessageID closest matching the NodeID. All ID's must be made with consistent hashs(link) and large enough to avoid collision.

Chord is a highly fault tolerant due to its routing overlay, and it stabilization routines. This makes it very accomadating to nodes constantly joining and leaving the network (As would be expected in a global public cloud such as Hermes).

## UDP Transport Protocol

All communication in the Hermes network will be done via UDP. UDP has a lower overhead than TCP, and more importantly it allows for UDP Hole Punching(link?). UDP Hole Punching is a requirement for Hermes to work with minimal configuration on the users part. However for the actual peer-to-peer communication there should be minimal 3rd Party ICE/STUN Servers for communicaiton, instead direct p2p communication will be investigated (Which will require a UDP Transport). The reason we have chosen not to use 3rd party ICE/STUN servers is because we want the Hermes network to be as resiliant and efficient as possible. Since a vast majoraty will be behind NAT's (without port forwarding), routing 80% (guesstamite) of the networks traffic through ICE servers would be less then ideal, even if it is only used for initially establishing a p2p connection.

However with the advantages of UDP, there are some draw backs, and for that reason the libutp protocol is a possible avenue to investigate to use ontop of UDP for more reliable communication

## Peer-To-Peer NAT'd Communication

Since most of the network will be behind NAT's, a reliable method is needed to communicate with peers behind them. Unfortuently there is no single reliable method, so Hermes will employ a couple. 

### UPnP

The UPnp (Universal Plug'n Play) protocol is used to easily configure, and communicate to devices on a network, IE Computers/Routers on a LAN. UPnP will allow users to easily configure their routers to allow communication from outside their LAN, aka port forwarding. Port forwarding is used on routers to allow outside traffic to communicate to a computer behind a NAT while still using a public facing IP and PORT. Unfortuently not all routers either A) Support UPnP or B) Have it disabled by default, and sticking the minial config work on the users end we don't expect them to know how to enable it.

### ICMP Hole Punching

ICMP Hole Punching is similar to that of UDP Hole Punching, where you can essentially "Punch" a hole in your NAT to allow outside traffic to use. However due to security restrictions UDP Hole Punching is limited to traffic orginating from only a single host, which isnt ideal for Hermes, since at all times a reliable communication from a multitude of hosts is required. Alternativley ICMP Hole Punching (Autonmous NAT Traversal Paper) allows for any host to connect through the "Punched hole". This is achieved similarily to how traceroute requests can always be routed back to the originating request Address (Even behind a NAT). Basically in ICMP Hole Punching you periodically send ICMP ECHO requests to a un-occupied IP adress (ie. 1.2.3.4 or 3.3.3.3), which basically is just a ping. Communication outbound from inside a LAN is allowed via a NAT. This initailly opens a port on the NAT to the internal computer. Communication back is done by sending a forged ICMP TIMEOUT packet back to the computer. The NAT will receive the packet and assume its a response from the original ICMP ECHO REQUEST, so it lets it pass through.

### ICE/STUN Servers

Even though this isn't an ideal means of communication, since it requires quite an elaborate dance of back and forth before an actual connection is established, and is required for every peer you want to connect to

