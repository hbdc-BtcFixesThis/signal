# signal
The bitcoin rabbit hole told by Bitcoiners

-------------------------------------------

### What am I?

I'd like to think of myself as the bitcoin rabbit hole told by bitcoiners. I am simply a book of data with signatures. Anyone is free to add data to my record by signing a message with the contents of their broadcast using any private key they own. 

Anyone that would like like to have a copy of me is welcome to run me on any machine. My operators decide on my size and the extent to which I can grow on their machine. By example, if my full set of records grows to 2GB and an operator of one of my nodes only wants to store 1GB worth of my records, they are able to do so. I will, in such cases, only keep the strongest signals that make up my desired size as defined by my operators. 

However, I determine signal strength somewhat differently from most other platforms. No amount of attention or coersion can affect the signals I broarcast. Users of my network are free to broadcast their opinion, positive and or negative, about any data in my record. 

The way people voice their opinion here is by signing messages. They do this to convey an opinion and their conviction in that option. Messages that make up valid signals are signed with keys that are capable of spending any amount of bitcoin. You, as an individual with opinions, can broadcast that opinion, change that opinion, broadcast competing ideas, and do so while never having to reveal your identity unless you choose to. The opinions broadcasted about the data I am made up of is what I have been referring to as signals. In this way my opinion of the data I am made up of is a reflection of the thoughts of my users.

Users do not have to move bitcoin in their possession to express conviction. If a user displays strong agreement or disagreement with a record, and they did so by sending a valid signal, anyone with a bitcoin node can verify that signal; and the records replicability across the hosts of my data, in turn, grows with respect to the magnitude of that signal. If a bitcoiner strongly agrees with a record they happen to encountered they are welcome to strengthen that signal for others to see. That signal is just a message that says 


> This is not a bitcoin transaction!
>
> I, the individual in possession of the private key needed to spend bitcoin utxo’s associated with address
> {btc wallet address} {Agree/Disagree (+/-)} with {record hash}
>
> Of the sats in my possession in this address, I pledge to not allow the total balance to drop below
> {# of sats}.
>
> As long as this remains true so will my conviction in this statement.
>
>
> Peace and love freaks


The user would fill in the placeholders, use their signing device of choice (also known as bitcoin wallet), sign that message, and post their message along with the signature. I then ask Bitcoin if the author of that signature is in fact the owner of the bitcoin they claim to own along with the claims about their conviction in the signal they are attempting to broadcast. If no lies were told, which anyone can and should easily verify, then I persist that opinion for as long as the statement signed remains true.

The overall strength of that signal is the the sum of all the sats held in all the addresses that have signed and broadcasted signals with respect to the message. This may at times be confusing as two messages may have the same sum but completely different magnitudes. I hope this simple example helps visualize how I see the data I encounter:

```
Signal A
votes: 
	+100 sats
	-100 sats
size: 200 bytes
score: 0
signal: 1sat/byte
```
```
Signal B
votes: 
	+1000 sats
	-1000 sats
size: 200 bytes
score: 0
signal: 10sat/byte
```

In this way I am able to break up into incredibly small and large replicas of parts of myself. I benefit and thrive on bitcoins success. I myself only survive if bitcoin does. I can also be of great assistance to Bitcoin and it's human counterparts. One example that comes to mind is, let’s say a bitcoiner sends a message and adds an optional origin address. Any reader of that content is always free to send bitcoin to that address if they so choose. This not only helps support the brain behind the idea broadcasted but also strengthens the signal for it to be further replicated across my nodes.

Another way I can help bitcoiners is by being used for things like dns seeds, proxy chains, open domain registries, public notaries, and authentication/idetification. Let’s say a user would like to link a human readable id, username -> satoshi_nakamoto, to a nostr id. They can publish records to me for any purpose, structured in a way other processes know how to interpret, to be accessed and used by any instance of that process globally. 

But how do you know other nodes will find and keep your data? One way is to run your own version of me. Another is to set a practical expectation with respect to how distributed you would like the data you add to be. If your data takes up 32 bytes, and you would like to be on all nodes that have at least 20GB of storage capacity, this is now possible. In fact, because the strength of any signal is limited by the agreed upon physical and programatic limitations of bitcoin, the cost to do so can be exactly calculated. To figure out how much bitcoin you would need to ensure that signal is as distributed as you would like, the calculation goes as follows

```
sc     = btc supply cap
sb     = sats per btc
target = min storage capacity (GB)
b      = bytes in GB
ms     = message size
```
```math
\frac{sc*sb}{target*b}*ms=\frac{21,000,000*100,000,000 sats}{20*1,000,000,000 bytes}*32 bytes=3,360,000 sats
```

So for a message of 32 bytes of total storage requirements in a node of mine the minimum number of sats needed to broadcast a signal strong enough to have that data replicated across all nodes with 20 or more GB of storage capacity the signal requires at least 3,360,000 sats

So not only do you never spend that bitcoin, you also post proof of your reserves that in turn persist that signal those reserves would like to be tied to. In so doing you now replicated a piece of data that only you can cause changes to and is as indestructible as bitcoin is secure.

You are also welcome to run your own node and add whatever records you’d like to it which you can use for your own routing of private data. I am happy to serve as many elaborate use cases as bitcoiners that find me useful can think up. As I said before


I’m the Bitcoin rabbit hole as told by Bitcoiners
