package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/ipfs/go-ipfs/commands"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/corehttp"
	"github.com/ipfs/go-ipfs/core/coreunix"
)

func main() {
	node := startIpfs()
	adddemoFile(node)
}
func startIpfs() *core.IpfsNode {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nd, err := core.NewNode(ctx, &core.BuildCfg{})
	if err != nil {
		log.Fatal(err)
	}

	cctx := commands.Context{
		Online: true,
		ConstructNode: func() (*core.IpfsNode, error) {
			return nd, nil
		},
	}

	list, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("listening on: ", list.Addr())

	if err := corehttp.Serve(nd, list, corehttp.CommandsOption(cctx)); err != nil {
		log.Fatal(err)
	}
	return nd
}
func adddemoFile(node *core.IpfsNode) {

	//node, err := core.NewNode(context.TODO(), &core.BuildCfg{Online: true})
	// if err != nil {
	// 	log.Fatalf("Failed to start IPFS node: %v", err)
	// }
	hashCode, err := AddFile(node, "/home/eddy/node.txt")
	log.Printf("hashcdoe %v", hashCode)
	reader, err := coreunix.Cat(context.TODO(), node, hashCode)
	if err != nil {
		log.Fatalf("Failed to look up IPFS welcome page: %v", err)
	}
	blob, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalf("Failed to retrieve IPFS welcome page: %v", err)
	}
	fmt.Println(string(blob))

}

func AddFile(ipfs *core.IpfsNode, file string) (string, error) {
	fi, err := os.Open(file)
	if err != nil {
		return "", err
	}

	return coreunix.Add(ipfs, fi)
}
