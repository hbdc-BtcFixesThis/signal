# signal
The bitcoin rabbit hole told by Bitcoiners

-------------------------------------------

### What am I?

I'd like to think of myself as the bitcoin rabbit hole, told by bitcoiners. I too am simply a book of data with signatures. Anyone is free to add data to my record by signing a message with the contents of their broadcast using any private key they own. 

Anyone that would like like to have a copy of me is welcome to run me on any machine. My operators decide on my size and the extent to which I can grow on their machine. By example, if my full set of records grows to 2GB and an operator of one of my nodes only wants to store 1GB worth of my records, they are able to do so. I will, in such cases, only keep the strongest signals that make up my desired size as defined by my operators. 

However, I determine signal strength somewhat differently from most other platforms. No amount of attention or coersion can affect the signals I broadcast. Users of my network are free to broadcast their opinion about any data in my record. 

The way people voice their opinion here is by signing messages. They do this to convey an opinion and their conviction in that option. Messages that make up valid signals are signed with keys that are capable of spending any amount of bitcoin. You, as an individual with opinions, can broadcast that opinion, change that opinion, broadcast competing ideas, and do so while never having to reveal your identity unless you choose to. The opinions broadcasted about the data I am made up of is what I have been referring to as signals. In this way my opinion of the data I am made up of is a reflection of the thoughts of my users.

Users do not have to move bitcoin in their possession to express conviction. If a user displays strong agreement or disagreement with a record, and they did so by sending a valid signal, anyone with a bitcoin node can verify that signal; and the records replicability across the hosts of my data, in turn, grows with respect to the magnitude of that signal. If a bitcoiner feels strongly with a record they happen to encountered, they are welcome to strengthen that signal for others to see. That signal is just a message that says 


> This is not a bitcoin transaction!
>
> Of the bitcoin in my possession,
> {x%} of the total amount of held in {address}
> shall be used to display my conviction
> in {record id}.
>
>
> Peace and love freaks

The user would fill in the placeholders, use their signing device of choice (also known as bitcoin wallet), sign that message, and post their message along with the signature. I then ask Bitcoin if the author of that signature is in fact the owner of the bitcoin they claim to own along with the claims about their conviction in the signal they are attempting to broadcast. If no lies were told, which anyone can and should easily verify, then I persist that opinion for as long as the statement signed remains true.

In this way, the author of the signal is able to specify the exact strength of that signal. This also means the address can be reused without bitcoin being required to move in for the same address to be associated with multiple signals.

```math
Address_{signals} = \set{Address_{total}*Signal_{{record}_i}, \dots\, Address_{total}*Signal_{{record}_n}}
```

The only requiremet I have is

```math
Address_{total} \geq \sum_{i=0}^{n} Address_{total}*Signal_{{record}_i}
```

Where $Address_{total}$ is the total amount of bitcoin held in an address thats broadcasted $Signal_{{record}_i}$ for $record_i$. Or in other words 

```math
\sum_{i=0}^{n}Signal_{{record}_i} \leq 100 percent
```

Which means that all signals broadcasted from an address must add up to no more then the bitcoin held in the address. I use percentages to avoid having to descisions with respect to which signals to hold. If I were to ask users how much bitcoin they'd like to use for the signal they're broadcasting then the sigal would be lost if the funds associated with the address goes below the amount specified in the signal. Now lets say multiple signals were sent and the total for all sigals broadcasted goes below the total associated with the address. To avoid potentially subjectively dropping signals I use percetages. In this way the signals always remain valid.

The overall strength of that signal is defined as

```math
record_i = \sum_{j=0}^{n} \frac{total_{Address_{j}}*Signal_{j_{{record}_i}}}{record_{size}}
```

Therefore if $record_x$ has 3 addresses, lets call them $total_{Address_{i}} = 1btc$, $total_{Address_{j}} = 3btc$ and $total_{Address_{k}} = 6btc$, each strengthening the signal for $record_x$ by 100%, 66%, and 34% respectively, and the $record_{size}$ is 2mb; the strength of the signal is

```math
record_i = \sum_{j=0}^{n} \frac{total_{Address_{j}}*Signal_{j_{{record}_i}}}{record_{size}} = \frac{1*1+3*.66+6*.34}{2mb} = ~2.5 bitcoin/mb
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
