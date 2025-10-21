# brimocomp/sftp
module for brimo being sftp client.

## General Purpose
improve reusability on the activity of file transfer using Secure File Transfer Protocol, between BRImo'S backend to remote SFTP server


## About The Repo
### Library Used
Libraries used here should be modular, easy to replace or upgrade, not a framework that tight couple with lot of components. no external library used in this module, everything native by go.

## How we use it?
Here are some example on how to use some of method on transferring file via SFTP. 


```
package main

import (
	"context"
	"fmt"
	"time"

	sftpbrimo "github.com/xanderxampp-be/franco/sftp"
)

func main() {
	done := make(chan bool)
	ctx := context.Background()

	// source and source2 are the fullpath of local file. Make sure its exist in your application root
	source := "brimo-estatement/aziz215.pdf"
	source2 := "brimo-estatement/ridwan215.pdf"

	// configuration for create SFTP connection
	user := "bdt"
	pass := "Bre@kthrough2312!"
	host := "103.183.74.243"
	port := "22"

	// dest and dest2 are full path of destination path on remote server
	dest := "/home/bdt/attachment/jea215.pdf"
	dest2 := "/home/bdt/attachment/ridwan215.pdf"

	sftpBrimo, err := sftpbrimo.New(user, pass, host, port)
	if err != nil {
		fmt.Println("error sftp connect : ", err)
		return
	}

	written, err := sftpBrimo.SendLocalFileToRemoteWithDelete(ctx, source, dest)
	if err != nil {
		fmt.Println("error send local file to remote : ", err)
		return
	}

	fmt.Println("written file : ", written)

	// simulate the disconnection case
	connCurrent := sftpBrimo.GetConnection(ctx)
	connCurrent.Close()
	time.Sleep(30 * time.Second)

	// reconnection should be sucess. then object sftpBrimo can send file to remote FTP again.
	written, err = sftpBrimo.SendLocalFileToRemoteWithDelete(ctx, source2, dest2)
	if err != nil {
		fmt.Println("error send local file to remote after reconnect : ", err)
		return
	}

	fmt.Println("written file after reconnect success : ", written)

	<-done
}
```


