package main

type NodeConf struct {
	*DBWithCache
}

func (nc *NodeConf) gop(k NodeConfKey, v []byte) []byte {
	return nc.GetOrPut(nc.ConfBucketName(), k.Bytes(), v, k.DefaultBytes())
}

func (nc *NodeConf) ConfBucketName() []byte         { return NodeConfBucket.DefaultBytes() }
func (nc *NodeConf) datapath(v []byte) []byte       { return nc.gop(Path, v) }
func (nc *NodeConf) Name(v []byte) []byte           { return nc.gop(Name, v) }
func (nc *NodeConf) Type(v []byte) []byte           { return nc.gop(Type, v) }
func (nc *NodeConf) PeersB(v []byte) []byte         { return nc.gop(Peers, v) }
func (nc *NodeConf) MaxRecordSize(v []byte) []byte  { return nc.gop(MaxRecordSize, v) }
func (nc *NodeConf) MaxStorageSize(v []byte) []byte { return nc.gop(MaxStorageSize, v) }

/*func (nc *NodeConf) NodeType(bucket []byte) NodeType {
	return NodeTypeFromBytes(nc.Type(bucket, nil))
}

func (nc *NodeConf) Peers(bucket []byte) ([]string, error) {
	var p []string
	return p, json.Unmarshal(nc.PeersB(bucket, nil), &p)
}*/
