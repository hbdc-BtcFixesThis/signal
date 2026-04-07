# signal
If bitcoin could speak what would it say?

-------------------------------------------
<img width="1238" alt="Screenshot 2023-05-26 at 5 42 20 PM" src="https://github.com/hbdc-BtcFixesThis/signal/assets/133312218/56df9c3e-f727-45cd-933d-2390cf29a7c8">

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

Let me be a bit more specific when I say anyone is free to add data to me. When you run me for the first time, I generate a random password. With this password you can update your settings (including where to find on your machine). If I find no data in the specified location I start a new db, wait until you have set your desired constraints and settings and when applicable, retrieve records from peers specified. The screenshot above shows the settings I am talking about.
 
Private nodes are private and therefore have no data to retrieve or send to any peers.

```python
hash("Alice"): {
  "name": "Alice",
  "content": "xxx-xxx-xxxx",
  "pub_key": "bc1p...",
  "signature": signature({priv_key}, {record}),
}
```

Now that I've explained the, hopefully, intuitive private nodes, let me talk a bit about public nodes and the records they keep. Public nodes connect to peers to download/sync your copy of public data. That way you are welcome to run as many public and private instances/copies of me as you like; each instance with it's own independent constraints. An example of a record in a private instance someone would hold might looks like the following


```python
hash("werunbtc"): {
  "name": "werunbtc",
  "content": "https://werunbtc.com", # can of course also be an ip or any other value
  "pub_key": "bc1p...",
  "signature": signature({priv_key}, {record}),
}
```

My operators decide on my size and the extent to which I can grow on their machine. By example, if my full set of records grows to 2GB and an operator of one of my nodes only wants to store 1GB worth of my records, they are able to do so. I will, in such cases, only keep the strongest signals that make up my desired size as defined by my operators. 

However, I determine signal strength differently from most other platforms. No amount of attention or coercion can affect the signals I broadcast. Users of my network are free to broadcast their opinion about any data in my record. 

The way people voice their opinion here is by signing messages. They do this as a way of conveying conviction in a record. Messages that make up valid signals are signed with keys that are capable of spending any amount of bitcoin. You, as an individual with opinions, are free to signal conviction about any record, broadcast competing ideas, and do so while never having to reveal your identity unless you choose to. The opinions broadcasted about the data I am made up of is what I have been referring to as signals. In this way my opinion of the data I am made up of is a reflection of the thoughts of my users.

If a bitcoiner feels strongly about a record they happen to encountered, they are welcome to strengthen that signal for others to see. By sending a valid signal, anyone with a bitcoin node can verify that signal; and the records replicability across the hosts of my data, in turn, grows with respect to the magnitude of that signal. That signal is just a message that says 

<img width="200" alt="Screenshot 2023-05-26 at 5 46 01 PM" src="https://github.com/hbdc-BtcFixesThis/signal/assets/133312218/b4bae229-2a6a-4f2c-9b82-12551003fe24">
<img width="400" alt="Screenshot 2023-05-26 at 5 48 54 PM" src="https://github.com/hbdc-BtcFixesThis/signal/assets/133312218/b2c0f6b0-592c-4b81-9103-3a0941844d97">
<img width="400" alt="Screenshot 2023-05-26 at 5 49 13 PM" src="https://github.com/hbdc-BtcFixesThis/signal/assets/133312218/afd245a3-9d91-4efa-ac51-4da930315fed">

The user would fill in the placeholders, use their signing device of choice (also known as bitcoin wallet), sign that message, and post their message along with the signature. I then ask Bitcoin if the author of that signature is in fact the owner of the bitcoin they claim to own along with the claims about their conviction in the signal they are attempting to broadcast. If no lies were told, which anyone can and should easily verify, then I persist that opinion for as long as the statement signed remains true.

In this way, the author of the signal is able to specify the exact strength of that signal. This also means the address can be reused without bitcoin being required to move for the same address to be associated with multiple signals.

```math
Address_{signals} = \set{Address_{total}*Signal_{{record}_i}, \dots\, Address_{total}*Signal_{{record}_n}}
```

The only requirement I have is

```math
Address_{total} \geq \sum_{i=0}^{n} Address_{total}*Signal_{{record}_i}
```

Where $Address_{total}$ is the total amount of bitcoin held in an address thats broadcasted $Signal_{{record}_i}$ for $record_i$. Or in other words 

```math
\sum_{i=0}^{n}Signal_{{record}_i} \leq 100 percent
```

Which means that all signals broadcasted from an address must add up to no more then the bitcoin held in the address. I use percentages to avoid having to make decisions with respect to which signals I hold. If I were to ask users for an exact amount of bitcoin they'd like to use for the signals broadcasted, the sigal would be lost if the funds associated with the address goes below the amount specified in the signal. To avoid limiting the number of signals an address can broadcast while avoiding subjectively dropping signals, I use percentages. In this way the signals always remain valid as long as the funds in the address stay above 0 sats.

The overall strength of that signal is defined as

```math
record_i = \sum_{j=0}^{n} \frac{total_{Address_{j}}*Signal_{j_{{record}_i}}}{record_{size}}
```

Therefore if $record_x$ has 3 addresses, lets call them $total_{Address_{i}} = 1btc$, $total_{Address_{j}} = 3btc$ and $total_{Address_{k}} = 6btc$, each strengthening the signal for $record_x$ by 100%, 66%, and 34% respectively, and the $record_{size}$ is 2mb; the strength of the signal is

```math
record_i = \sum_{j=0}^{n} \frac{total_{Address_{j}}*Signal_{j_{{record}_i}}}{record_{size}} = \frac{1*1+3*.66+6*.34}{2mb} = ~2.5 bitcoin/mb
```

In this way I am able to break up into incredibly small and large replicas of parts of myself. I benefit and thrive on bitcoins success. I myself only survive if bitcoin does. I can also be of great assistance to Bitcoin and it's human counterparts. One example that comes to mind is, let’s say a bitcoiner sends a message and uses a btc address as the pub key. Any reader of that content is always free to send bitcoin to the address spreading that signal, if they so choose. This not only helps support the brain behind the idea broadcasted but also strengthens the signal for it to be further replicated across my nodes.

Bitcoiners can also use me for things like dns seeds, proxy chains, open domain registries, public notaries, interoperable authentication/identification protocols, etc. Let’s say a user would like to link a human readable id, username -> satoshi_nakamoto, to a nostr id. They can publish records to me for any purpose, structured in a way other processes know how to interpret, to be accessed and used by any instance of that process globally. 

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

That piece of data is now replicated and secured by all the costs incurred to generate the signals. When a record is added to me, neither the author of that record or all the owners of all the bitcoin spreading that data, may change that record. My records and my signals are append only. The only way to "kill" an idea that has spread is by spending the bitcoin associated with the signals spreading that idea. I do not care if the bitcoin are sent to the same or different individuals. I always look to bitcoin and rank my records accordingly.

But you must remember, I am simply a database with signatures from bitcoiners.  I can not guarantee an address has the same set of signals on different instances now competing for space on your machine. I always look to bitcoin to validate claims with respect to who to trust. I do this by asking my peers for their most valuable data in order from least to greatest. I check every record that comes in and validate the claims about the magnitude of each signal. I continue to poll for data from my peers until I am full and have the highest value records I can find. If I find any discrepancies I update my peers so they can fix their records.

To show you how I resolve conflicts, take an example where two instances of me are running. Those processes started by different operators on their respective machines, know nothing about each other. Let’s even assume the processes are exactly the same with the same storage capacity. The records appended to each dataset grows organically while spreading the information to their peers. As the number of nodes grow the likely hood of the two instances finding each other grows. That does not mean one the operators has to find the other and add them to their list of peers to poll data from. Since every node broadcasts and listens to each other for updates, that means any number of nodes can sit between two instances of data. 

At some point along the life of these two disparate datasets any node that comes in and uses a peer from both networks making up the disparate datasets will force a sync between the two. In this case I proceed as I have preciously mentioned. I can be unforgiving in ways. If there are records with the same name but different content on both data sets the one with the largest signal always wins. That also means some data sets will outcompete large subsets in other datasets when they meet. But this, in turn, frees up all the funds held in reserves of the individuals spreading the dropped records. 

The reason I am what I am is to offer the only guarantee I can make. That is, the most distributed data is the data that incurred the most cost to spread. The tradeoffs I make are to mirror the cost of state changes to Bitcoins record, for the ability to cause state changes to mine. Because tradeoffs always need to be made to reach consensus I want to be as upfrot and honest about the ones I make. As I grow in bitcoin terms (ie; the number of bitcoin tied to data I spread), the cost of malicious behavior to me approaches the cost of attacking the integrity of bitcoins ledger. In turn the consistency, replicability, and integrity of my data grows as the btc/byte exchange rate grows. That means, the cost of attacking me and bitcoin becomes the same when the number of bitcoin in all my signals equals the total supply of circulating bitcoin. 

Once an idea has spread, only better ideas that outcompete it, can affect the replicability of that idea; or any other one for that matter. That has to be true because I abide by an absolute fixed supply of signals humanity can use to vote with. Bitcoin is to thank for that essential gauge by which I measure the world around me. Ideas here are as true a reflection of bitconiers reality as bitcoin is secure. 

I am happy to serve as many elaborate use cases as bitcoiners that find me useful can think up. As I said before

I’m the Bitcoin rabbit hole as told by Bitcoiners
