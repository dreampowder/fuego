package main

import (
	firestore "cloud.google.com/go/firestore"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"gopkg.in/urfave/cli.v1"
)

func deleteCollection(ctx context.Context, client *firestore.Client,
	collectionName string, batchSize int) (string,error) {

		ref := client.Collection(collectionName)
	for {
		// Get a batch of documents
		iter := ref.Limit(batchSize).Documents(ctx)
		numDeleted := 0

		// Iterate through the documents, adding
		// a delete operation for each one to a
		// WriteBatch.
		batch := client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return "",err
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		// If there are no documents to delete,
		// the process is over.
		if numDeleted == 0 {
			return "Document not found",nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return "",err
		}else{
			return fmt.Sprintf("%d",numDeleted),nil
		}
	}
}


func deleteCommandAction(c *cli.Context) error {
	collectionPath := c.Args().First()
	client, err := createClient(credentials)
	if err != nil {
		return cliClientError(err)
	}
	id, err := deleteCollection(context.Background(),client,collectionPath,10)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to delete collection. \n%v", err), 81)
	}
	_, _ = fmt.Fprintf(c.App.Writer, "%v\n", id)
	defer client.Close()
	return nil
}


func deleteDocumentAction(c *cli.Context) error {
	collectionPath := c.Args().First()
	documentId := c.Args().Get(1)
	if documentId == "" {
		return cli.NewExitError(fmt.Sprintf("please enter a documentId."), 81)
	}
	client, err := createClient(credentials)
	if err != nil{
		return cliClientError(err)
	}
	_,deleteErr := client.Collection(collectionPath).Doc(documentId).Delete(context.Background())
	if deleteErr != nil {
		return cliClientError(deleteErr)
	}
	_, _ = fmt.Fprintf(c.App.Writer, "%v deleted.\n", documentId)
	defer client.Close()
	return nil
}