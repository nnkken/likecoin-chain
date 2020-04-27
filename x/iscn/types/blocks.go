package types

import (
	"log"

	goblocks "github.com/ipfs/go-block-format"
	cid "github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	node "github.com/ipfs/go-ipld-format"
	iscnblocks "github.com/likecoin/iscn-ipld/plugin/block"
	mh "github.com/multiformats/go-multihash"
)

type IscnBlock struct {
	IscnRecord
	cid     *cid.Cid
	rawdata []byte
}

// Cid returns the CID of the block header
func (b *IscnBlock) Cid() cid.Cid {
	return *b.cid
}

// RawData returns the binary of the CBOR encode of the block header
func (b *IscnBlock) RawData() []byte {
	return b.rawdata
}

// Resolve resolves a path through this node, stopping at any link boundary
// and returning the object found as well as the remaining path to traverse
func (*IscnBlock) Resolve(path []string) (interface{}, []string, error) {
	return nil, nil, nil
}

// Tree lists all paths within the object under 'path', and up to the given depth.
// To list the entire object (similar to `find .`) pass "" and -1
func (*IscnBlock) Tree(path string, depth int) []string {
	return nil
}

// node.Node interface

// Copy will go away. It is here to comply with the Node interface.
func (*IscnBlock) Copy() node.Node {
	panic("dont use this yet")
}

// Links is a helper function that returns all links within this object
// HINT: Use `ipfs refs <cid>`
func (*IscnBlock) Links() []*node.Link {
	return nil
}

// ResolveLink is a helper function that allows easier traversal of links through blocks
func (*IscnBlock) ResolveLink(path []string) (*node.Link, []string, error) {
	return nil, []string{}, nil
}

// Size will go away. It is here to comply with the Node interface.
func (*IscnBlock) Size() (uint64, error) {
	return 0, nil
}

// Stat will go away. It is here to comply with the Node interface.
func (*IscnBlock) Stat() (*node.NodeStat, error) {
	return &node.NodeStat{}, nil
}

// github.com/ipfs/go-ipld-format.DecodeBlockFunc

// BlockDecoder takes care of the iscn-block IPLD objects (ISCN block headers)
func BlockDecoder(block goblocks.Block) (node.Node, error) {
	n, err := cbor.DecodeBlock(block)
	if err != nil {
		log.Printf("Cannot not decode block: %s", err)
		return nil, err
	}
	return n, nil
}

// Package function

// NewISCNBlock creates a iscn-block IPLD object
func NewISCNBlock(m map[string]interface{}) (*IscnBlock, error) {
	//TODO: validation code go here

	rawdata, err := cbor.DumpObject(m)
	if err != nil {
		log.Printf("Fail to marshal object: %s", err)
		return nil, err
	}

	c, err := cid.Prefix{
		Codec:    iscnblocks.CodecISCN,
		Version:  1,
		MhType:   mh.SHA2_256,
		MhLength: -1,
	}.Sum(rawdata)
	if err != nil {
		log.Printf("Fail to create CID: %s", err)
		return nil, err
	}

	block := IscnBlock{
		cid:     &c,
		rawdata: rawdata,
	}
	return &block, nil
}
