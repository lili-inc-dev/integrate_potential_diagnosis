package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// https://firebase.google.com/docs/reference/rest/auth#section-create-email-password

func main() {
	var (
		serviceFilePath = flag.String("s", "../../service-account-file.json", "firebase service-account-file path")
		firebaseUID     = flag.String("u", "", "firebase UID")
	)

	flag.Parse()

	if err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", *serviceFilePath); err != nil {
		panic(fmt.Sprintf("failed to set env: %#v", err))
	}

	auth, err := newFirebaseAuthClient()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize firebase auth client: %#v", err))
	}

	fmt.Println("Input custom claim:")
	/* 例:
	{
		"role": "admin",
		"is_valid": true
	}
	*/

	var customClaim map[string]interface{}
	dec := json.NewDecoder(bufio.NewReader(os.Stdin))
	if err := dec.Decode(&customClaim); err != nil {
		panic(fmt.Sprintf("failed to decode custom claim: %#v", err))
	}

	if err := auth.SetCustomUserClaims(context.Background(), *firebaseUID, customClaim); err != nil {
		panic(fmt.Sprintf("failed to SetCustomUserClaims: %#v", err))
	}

	fmt.Println("succeeded")
}

func newFirebaseApp() (*firebase.App, error) {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %w", err)
	}

	return app, nil
}

func newFirebaseAuthClient() (*auth.Client, error) {
	app, err := newFirebaseApp()
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize firebase auth client: %w", err)
	}

	return client, nil
}
