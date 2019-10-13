package main

import (
	firestore "cloud.google.com/go/firestore"
	"context"
	"fmt"
	"gopkg.in/urfave/cli.v1"
)

func addData(
	client *firestore.Client,
	collection string,
	documentId string,
	data string,
	timestampify bool) (string, error) {

	object, err := unmarshallData(data)
	if err != nil {
		return "", err
	}

	if timestampify {
		timestampifyMap(object)
	}

	if documentId == "auto-id" {
		doc, _, err := client.
			Collection(collection).
			Add(context.Background(), object)
		if err != nil {
			return "", err
		}
		return doc.ID, nil
	}else{
		_, err := client.
			Collection(collection).
			Doc(documentId).
			Set(context.Background(), object)
		if err != nil {
			return "", err
		}
		return documentId, nil
	}
}

func addCommandAction(c *cli.Context) error {
	collectionPath := c.Args().First()
	timestampify := c.Bool("timestamp")
	documentId := c.String("documentid")

	data := c.Args().Get(1)

	client, err := createClient(credentials)
	if err != nil {
		return cliClientError(err)
	}
	id, err := addData(client, collectionPath,documentId, data, timestampify)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to add data. \n%v", err), 81)
	}
	fmt.Fprintf(c.App.Writer, "%v\n", id)
	defer client.Close()
	return nil
}
