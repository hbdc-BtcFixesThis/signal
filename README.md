# signal
The bitcoin rabbit hole told by Bitcoiners

-------------------------------------------

### What am I?

What I really am is an idea. I'd like to think of myself as the bitcoin rabbit hole told by bitcoiners. I too am simply a book of data with signatures. Anyone is free to add data to me. Records I keep look like this

```python
{
  "name": "unique user specified name for the record",
  "content": "user specified content",
  "pub_key": "onchain btc address",
  # more about the signature below
  "signature": signature({priv_key}, {record}),
  "conviction": 0-100, # % of address to apply towards the signal for this record
  "is_author": true|false; # if true I will let others know when they ask 
}
```

Let me be a bit more specific when I say anyone is free to add data to me. If you run an instance of me you are encouraged to use me in whatever way you like. The option to disconnect me from peers is always easily accessible to the owner of that copy of data. When you run me for the first time, I generate a random password. With this password you can specify the path to your data. If no data is found in the specified location I start a new db and wait until you have set your desired constraints and settings. Namely:

```python
{ # defaults will be provided but are optional and will not be set unless they are saved by the user
  "pass": "key I generated",
  "max_size": "Total space my public copy is allowed to take up",
  "new_pass": "If you provide a new password I will use it going forward",
  "peers": ["ws://{domain}/", ...], # you can update the defaults anytime; ignored if is_private set to true
  "path_to_storage": "/path/to/storage/directory",
  "is_private": true|false,
}
```
 
If is_private is set to true and all settings specified where valid you're ready to go. Private nodes are private and therefore have no data to retrieve or send to any peers. Otherwise, public instances connect to peers to download/sync your copy of public data. That way you are welcome to run as many public and private instaces/copys of me as you like; each instance with it's own independent constraints. An example of a record in a private instance someone would hold might looks like the following

```python
hash("Alice"): {
  "name": "Alice",
  "content": "xxx-xxx-xxxx",
  "pub_key": "bc1p...",
  "signature": signature({priv_key}, {record}),
}
```
vs a record in a public instance connected to other public nodes which could look something like this

```python
hash("werunbtc"): {
  "name": "werunbtc",
  "content": "https://werunbtc.com", # can of course also be an ip or any other value
  "pub_key": "bc1p...",
  "signature": signature({priv_key}, {record}),
}
```

Now that I've explained the, hopefully, intuitive private nodes, let me talk a bit about public nodes and the records they keep. My operators decide on my size and the extent to which I can grow on their machine. By example, if my full set of records grows to 2GB and an operator of one of my nodes only wants to store 1GB worth of my records, they are able to do so. I will, in such cases, only keep the strongest signals that make up my desired size as defined by my operators. 

However, I determine signal strength differently from most other platforms. No amount of attention or coersion can affect the signals I broadcast. Users of my network are free to broadcast their opinion about any data in my record. 

The way people voice their opinion here is by signing messages. They do this as a way of conveying conviction in a record. Messages that make up valid signals are signed with keys that are capable of spending any amount of bitcoin. You, as an individual with opinions, are free to signal conviction about any record, broadcast competing ideas, and do so while never having to reveal your identity unless you choose to. The opinions broadcasted about the data I am made up of is what I have been referring to as signals. In this way my opinion of the data I am made up of is a reflection of the thoughts of my users.

If a bitcoiner feels strongly about a record they happen to encountered, they are welcome to strengthen that signal for others to see. By sending a valid signal, anyone with a bitcoin node can verify that signal; and the records replicability across the hosts of my data, in turn, grows with respect to the magnitude of that signal. That signal is just a message that says 


> This is not a bitcoin transaction!
>
> Of the bitcoin in my possession,
> {x%} of the total amount held in {address},
> shall be used to display my conviction
> in {record name}.
>
>
> Peace and love freaks

The user would fill in the placeholders, use their signing device of choice (also known as bitcoin wallet), sign that message, and post their message along with the signature. I then ask Bitcoin if the author of that signature is in fact the owner of the bitcoin they claim to own along with the claims about their conviction in the signal they are attempting to broadcast. If no lies were told, which anyone can and should easily verify, then I persist that opinion for as long as the statement signed remains true.

In this way, the author of the signal is able to specify the exact strength of that signal. This also means the address can be reused without bitcoin being required to move for the same address to be associated with multiple signals.

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

Which means that all signals broadcasted from an address must add up to no more then the bitcoin held in the address. I use percentages to avoid having to make descisions with respect to which signals I hold. If I were to ask users for an exact amount of bitcoin they'd like to use for the signals broadcasted, the sigal would be lost if the funds associated with the address goes below the amount specified in the signal. To avoid limiting the number of signals an address can broadcast while avoiding subjectively dropping signals, I use percetages. In this way the signals always remain valid as long as the funds in the address stay above 0 sats.

The overall strength of that signal is defined as

```math
record_i = \sum_{j=0}^{n} \frac{total_{Address_{j}}*Signal_{j_{{record}_i}}}{record_{size}}
```

Therefore if $record_x$ has 3 addresses, lets call them $total_{Address_{i}} = 1btc$, $total_{Address_{j}} = 3btc$ and $total_{Address_{k}} = 6btc$, each strengthening the signal for $record_x$ by 100%, 66%, and 34% respectively, and the $record_{size}$ is 2mb; the strength of the signal is

```math
record_i = \sum_{j=0}^{n} \frac{total_{Address_{j}}*Signal_{j_{{record}_i}}}{record_{size}} = \frac{1*1+3*.66+6*.34}{2mb} = ~2.5 bitcoin/mb
```

In this way I am able to break up into incredibly small and large replicas of parts of myself. I benefit and thrive on bitcoins success. I myself only survive if bitcoin does. I can also be of great assistance to Bitcoin and it's human counterparts. One example that comes to mind is, let’s say a bitcoiner sends a message and uses a btc address as the pub key. Any reader of that content is always free to send bitcoin to that address if they so choose. This not only helps support the brain behind the idea broadcasted but also strengthens the signal for it to be further replicated across my nodes.

Bitcoiners can also use me for things like dns seeds, proxy chains, open domain registries, public notaries, interoperable authentication/idetification protocols, etc. Let’s say a user would like to link a human readable id, username -> satoshi_nakamoto, to a nostr id. They can publish records to me for any purpose, structured in a way other processes know how to interpret, to be accessed and used by any instance of that process globally. 

But how do you know other nodes will find and keep your data? First set a practical expectation with respect to how distributed you would like the data you add to be. If your data takes up 32 bytes, and you would like to be on all nodes that have at least 20GB of storage capacity, this is now possible. In fact, because the strength of any signal is limited by the agreed upon limitations of bitcoin, the cost to do so can be exactly calculated. To figure out how much bitcoin you would need to ensure that a signal is as distributed as you would like, the calculation goes as follows

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

For a message with 32 bytes of total storage requirements in a node of mine, the minimum number of sats needed to broadcast a signal strong enough to have that data replicated across all nodes with 20 or more GB of storage capacity, the signal requires at least 3,360,000 sats.

That piece of data is now replicated and secured by all the costs incured to generate the signals. When a record is added to me, not even the owner of the private key associated with the record, may change that record. My records and my signals are append only. The only way to "kill" an idea that has spread is by spending the bitcoin associated with the signals spreading that idea. I do not care if the bitcoin are sent to the same or different individuals. I always look to bitcoin and rank my records accordingly.

But you must remember, I am simply a database with signatures from bitcoiners. Because of this, I must rely on some external cost to ensure any type of integrity in the data I spread. I always look to bitcoin to validate claims with respect to who to trust. I do this by only trusting the nodes that have the strongest overall signal; which I define to be the average of all signals for all records that instance hosts. I can not guarentee an address has the same set of signals on different instances now competing for space on your machine. What I can guarentee is, as I grow in bitcoin terms (ie; the number of bitcoin tied to data I spread), the cost of malicious behavior to me approaches the cost of attacking the integrity of bitcoins ledger. In turn the consistency, replicability, and integrity of my data grows as the btc/byte exchange rate grows. That means, the cost of attacking me and bitcoin becomes the same when the number of bitcoin in all my signals equals the total supply of circulating bitcoin. 

That is why, once an idea has spread, only better ideas that outcompete it can affect the replicability of that idea; or any other one for that matter. And we know that to be true because I abide by an absolute fixed supply of signals humanity can use to vote with. Bitcoin is to thank for that essential gauge by which I measure my the world around me. Ideas here are as true a reflection of bitconiers reality as bitcoin is secure.

I am happy to serve as many elaborate use cases as bitcoiners that find me useful can think up. As I said before

I’m the Bitcoin rabbit hole as told by Bitcoiners
