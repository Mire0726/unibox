package firebase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"firebase.google.com/go/auth"

	"github.com/Mire0726/unibox/backend/config"
	"github.com/Mire0726/unibox/backend/internal/cerror"
	"github.com/Mire0726/unibox/backend/pkg/log"
	"github.com/joho/godotenv"
)

type FirebaseAuth interface {
	SignUpWithEmailPassword(ctx context.Context, email, password string) (*SignUpResponse, error)
	SignInWithEmailPassword(ctx context.Context, email, password string) (*SignInResponse, error)
	SendPasswordResetEmail(ctx context.Context, email string) (*SendPasswordResetEmailResponse, error)
	VerifyPasswordResetCode(ctx context.Context, oobCode string) (*VerifyPasswordResetCodeResponse, error)
	ConfirmPasswordReset(ctx context.Context, oobCode, newPassword string) (*ConfirmPasswordResetResponse, error)
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
	GetUser(ctx context.Context, uid string) (*UserRecord, error)
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type FirebaseAPIError struct {
	Error Error `json:"error"`
}

func (c *AuthClient) readErr(res io.ReadCloser) (*FirebaseAPIError, error) {
	body, err := io.ReadAll(res)
	if err != nil {
		return nil, cerror.Wrap(err, "firebase", cerror.WithIOCode())
	}

	firebaseAPIError := &FirebaseAPIError{}
	if err := json.Unmarshal(body, &firebaseAPIError); err != nil {
		return nil, cerror.Wrap(err, "firebase", cerror.WithEncodingJSONCode())
	}

	return firebaseAPIError, nil
}

type AuthClient struct {
	client *auth.Client
	logger *log.Logger
}

func NewClient(ctx context.Context, logger *log.Logger) (*AuthClient, error) {
	app, err := initializeApp(ctx)
	if err != nil {
		return nil, cerror.Wrap(err, "firebase")
	}

	client, err := app.Auth(ctx)
	if err != nil {
		logger.Error("Failed to get auth client", log.Fstring("package", "firebase"), log.Ferror(err))
		return nil, cerror.Wrap(err, "firebase", cerror.WithFirebaseCode())
	}

	return &AuthClient{
		client,
		logger,
	}, nil
}

func NewClientWithoutLogger(ctx context.Context) (*auth.Client, error) {
	app, err := initializeApp(ctx)
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, cerror.Wrap(err, "firebase", cerror.WithFirebaseCode())
	}

	return client, nil
}

type signUpRequest struct {
	ReturnSecureToken bool `json:"returnSecureToken"`
}

type signUpRequestWithEmailPassword struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type SignUpResponse struct {
	IDToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`
}

func (c *AuthClient) SignUpWithEmailPassword(ctx context.Context, email, password string) (*SignUpResponse, error) {
	if err := godotenv.Load("../.env"); err != nil {
		c.logger.Error("Error loading .env file", log.Ferror(err))
	}
	firebaseAPIKey := os.Getenv("FIREBASE_API_KEY")
	reqBody := &signUpRequestWithEmailPassword{
		Email:             email,
		Password:          password,
		ReturnSecureToken: true,
	}

	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signUp?key=%s", firebaseAPIKey)

	signUpResponse := new(SignUpResponse)
	if err := c.callPost(ctx, url, reqBody, signUpResponse); err != nil {
		return nil, err
	}

	return &SignUpResponse{
		IDToken: signUpResponse.IDToken,
		LocalID: signUpResponse.LocalID,
		Email:  signUpResponse.Email,
		RefreshToken: signUpResponse.RefreshToken,
		ExpiresIn: signUpResponse.ExpiresIn,
	}, nil
}

type signInRequestWithEmailPassword struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type SignInResponse struct {
	IDToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`
	Registered   bool   `json:"registered"`
}

func (c *AuthClient) SignInWithEmailPassword(ctx context.Context, email, password string) (*SignInResponse, error) {
	if err := godotenv.Load("../.env"); err != nil {
		c.logger.Error("Error loading .env file", log.Ferror(err))
	}
	firebaseAPIKey := os.Getenv("FIREBASE_API_KEY")

	reqBody := &signInRequestWithEmailPassword{
		Email:             email,
		Password:          password,
		ReturnSecureToken: true,
	}

	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", firebaseAPIKey)

	signInResponse := new(SignInResponse)
	if err := c.callPost(ctx, url, reqBody, signInResponse); err != nil {
		return nil, err
	}

	return &SignInResponse{
		IDToken: signInResponse.IDToken,
		Email:  signInResponse.Email,
		RefreshToken: signInResponse.RefreshToken,
		ExpiresIn: signInResponse.ExpiresIn,
		LocalID: signInResponse.LocalID,
		Registered: signInResponse.Registered,
	}, nil
}

type sendPasswordResetEmailRequest struct {
	Email       string `json:"email"`
	RequestType string `json:"requestType"`
}

type SendPasswordResetEmailResponse struct {
	Email string `json:"email"`
}

func (c *AuthClient) SendPasswordResetEmail(ctx context.Context, email string) (*SendPasswordResetEmailResponse, error) {
	firebaseAPIKey := config.GetEnv().FirebaseAPIKey

	reqBody := &sendPasswordResetEmailRequest{
		Email:       email,
		RequestType: "PASSWORD_RESET",
	}

	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:sendOobCode?key=%s", firebaseAPIKey)

	sendPasswordResetEmailResponse := &SendPasswordResetEmailResponse{}
	if err := c.callPost(ctx, url, reqBody, &sendPasswordResetEmailResponse); err != nil {
		c.logger.Error("Failed to post", log.Fstring("package", "firebase"), log.Ferror(err))

		return nil, err
	}

	return sendPasswordResetEmailResponse, nil
}

type verifyPasswordResetCodeRequest struct {
	OobCode string `json:"oobCode"`
}

type VerifyPasswordResetCodeResponse struct {
	Email       string `json:"email"`
	RequestType string `json:"requestType"`
}

func (c *AuthClient) VerifyPasswordResetCode(ctx context.Context, oobCode string) (*VerifyPasswordResetCodeResponse, error) {
	firebaseAPIKey := config.GetEnv().FirebaseAPIKey

	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:resetPassword?key=%s", firebaseAPIKey)

	reqBody := &verifyPasswordResetCodeRequest{
		OobCode: oobCode,
	}

	verifyPasswordResetCodeResponse := &VerifyPasswordResetCodeResponse{}
	if err := c.callPost(ctx, url, reqBody, &verifyPasswordResetCodeResponse); err != nil {
		c.logger.Error("Failed to post", log.Fstring("package", "firebase"), log.Ferror(err))

		return nil, err
	}

	return verifyPasswordResetCodeResponse, nil
}

type confirmPasswordResetRequest struct {
	OobCode     string `json:"oobCode"`
	NewPassword string `json:"newPassword"`
}

type ConfirmPasswordResetResponse struct {
	Email       string `json:"email"`
	RequestType string `json:"requestType"`
}

func (c *AuthClient) ConfirmPasswordReset(ctx context.Context, oobCode, newPassword string) (*ConfirmPasswordResetResponse, error) {
	firebaseAPIKey := config.GetEnv().FirebaseAPIKey

	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:resetPassword?key=%s", firebaseAPIKey)

	reqBody := confirmPasswordResetRequest{
		OobCode:     oobCode,
		NewPassword: newPassword,
	}

	confirmPasswordResetResponse := &ConfirmPasswordResetResponse{}
	if err := c.callPost(ctx, url, reqBody, &confirmPasswordResetResponse); err != nil {
		c.logger.Error("Failed to post", log.Fstring("package", "firebase"), log.Ferror(err))

		return nil, err
	}

	return confirmPasswordResetResponse, nil
}

func (c *AuthClient) callPost(ctx context.Context, url string, reqBody any, respBody interface{}) error {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		c.logger.Error("Failed to encoding json", log.Fstring("package", "firebase"), log.Ferror(err))

		return cerror.Wrap(err, "firebase", cerror.WithEncodingJSONCode())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error("Failed to create request", log.Fstring("package", "firebase"), log.Ferror(err))

		return cerror.Wrap(err, "firebase", cerror.WithCreateExternalHTTPRequestCode())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to do request", log.Fstring("package", "firebase"), log.Ferror(err))

		return cerror.Wrap(err, "firebase", cerror.WithDoExternalHTTPRequestCode())
	}

	code := cerror.MapHTTPErrorToCode(resp.StatusCode)
	if code != cerror.OK {
		firebaseAPIError, err := c.readErr(resp.Body)
		if err != nil {
			return err
		}
		c.logger.Error(firebaseAPIError.Error.Message, log.Fstring("package", "firebase"), log.Ferror(err))

		return cerror.New(code.String(), cerror.WithCode(code))
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to read body", log.Fstring("package", "firebase"), log.Ferror(err))

		return cerror.Wrap(err, "firebase", cerror.WithIOCode())
	}

	if err := json.Unmarshal(body, &respBody); err != nil {
		c.logger.Error("Failed to unmarshal body", log.Fstring("package", "firebase"), log.Ferror(err))

		return cerror.Wrap(err, "firebase", cerror.WithEncodingJSONCode())
	}

	return nil
}

func (c *AuthClient) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := c.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		c.logger.Error("Failed to verify ID token", log.Fstring("package", "firebase"), log.Ferror(err))
		return nil, cerror.Wrap(err, "firebase", cerror.WithUnauthorizedCode())
	}

	return token, nil
}

type VerifyTokenResponse struct {
	ID string `json:"uid"`
}

type UserRecord struct {
    UID   string `json:"uid"`
    Email string `json:"email"`
}

func (c *AuthClient) GetUser(ctx context.Context, uid string) (*UserRecord, error) {
    user, err := c.client.GetUser(ctx, uid)
    if err != nil {
        c.logger.Error("Failed to get user", log.Fstring("package", "firebase"), log.Ferror(err))
        return nil, cerror.Wrap(err, "firebase", cerror.WithFirebaseCode())
    }

    return &UserRecord{
        UID:   user.UID,
        Email: user.Email,
    }, nil
}
