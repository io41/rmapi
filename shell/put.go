package shell

import (
	"errors"
	"fmt"

	"github.com/abiosoft/ishell"
	"github.com/juruen/rmapi/util"
)

func updateCmd(ctx *ShellCtxt) *ishell.Cmd {
	return &ishell.Cmd{
		Name:      "update",
		Help:      "update/overwrite an existing document in cloud",
		Completer: createFsEntryCompleter(),
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Err(errors.New("missing source file"))
				return
			}

			srcName := c.Args[0]
			node := ctx.node
			var err error

			if len(c.Args) == 2 {
				node, err = ctx.api.Filetree().NodeByPath(c.Args[1], ctx.node)
				if err != nil || node.IsFile() {
					c.Err(errors.New("directory doesn't exist"))
					return
				}
			}

			c.Printf("updating: [%s]...", srcName)
			dstDir := node.Id()
			document, err := ctx.api.UploadDocument(dstDir, srcName, true)
			if err != nil {
				c.Err(fmt.Errorf("Failed to update file [%s] %v", srcName, err))
				return
			}

			c.Println("OK")
			ctx.api.Filetree().AddDocument(document)
		},
	}
}

func putCmd(ctx *ShellCtxt) *ishell.Cmd {
	return &ishell.Cmd{
		Name:      "put",
		Help:      "copy a local document to cloud",
		Completer: createFsEntryCompleter(),
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Err(errors.New("missing source file"))
				return
			}

			srcName := c.Args[0]
			docName, _ := util.DocPathToName(srcName)

			node := ctx.node
			var err error

			if len(c.Args) == 2 {
				node, err = ctx.api.Filetree().NodeByPath(c.Args[1], ctx.node)

				if err != nil || node.IsFile() {
					c.Err(errors.New("directory doesn't exist"))
					return
				}
			}

			_, err = ctx.api.Filetree().NodeByPath(docName, node)
			//TODO: force flag and overwrite
			if err == nil {
				c.Err(errors.New("entry already exists"))
				return
			}

			c.Printf("uploading: [%s]...", srcName)

			dstDir := node.Id()

			document, err := ctx.api.UploadDocument(dstDir, srcName, true)

			if err != nil {
				c.Err(fmt.Errorf("Failed to upload file [%s] %v", srcName, err))
				return
			}

			c.Println("OK")

			ctx.api.Filetree().AddDocument(document)
		},
	}
}
