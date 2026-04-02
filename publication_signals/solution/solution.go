package main

import "fmt"

type Snapshot struct {
	Value string
}

type Catalog struct {
	Value string
}

func startSnapshotPublisher(start <-chan struct{}) (*Snapshot, <-chan struct{}) {
	snapshot := &Snapshot{}
	ready := make(chan struct{})

	go func() {
		<-start
		snapshot.Value = "snapshot published"
		ready <- struct{}{}
	}()

	return snapshot, ready
}

func startCatalogPublisher(start <-chan struct{}) (*Catalog, <-chan struct{}) {
	catalog := &Catalog{}
	ready := make(chan struct{})

	go func() {
		<-start
		catalog.Value = "catalog published"
		close(ready)
	}()

	return catalog, ready
}

func main() {
	startSnapshot := make(chan struct{})
	snapshot, snapshotReady := startSnapshotPublisher(startSnapshot)
	close(startSnapshot)
	<-snapshotReady
	fmt.Println(snapshot.Value)

	startCatalog := make(chan struct{})
	catalog, catalogReady := startCatalogPublisher(startCatalog)
	close(startCatalog)
	<-catalogReady
	fmt.Println(catalog.Value)
}
