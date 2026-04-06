package keycloak

import (
	"LMSSystem/entitiesDTO"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Nerzal/gocloak/v13"
)

type KeycloakService interface {
	CreateUser(ctx context.Context, signUpReq entitiesDTO.SignupReq) (string, error)
	LoginUser(ctx context.Context, loginReq entitiesDTO.LoginReq) (entitiesDTO.LoginResponse, error)
	KeycloakClient(ctx context.Context) (string, error)
	GetClientSecret(ctx context.Context, internalID string) (string, error)
	GetSecret() string
	RefreshToken(ctx context.Context, refreshToken string) (entitiesDTO.LoginResponse, error)
	GetUser(ctx context.Context, userID string) (*gocloak.User, error)
	DeleteUser(ctx context.Context, userID string) error
	UpdateUser(ctx context.Context, user entitiesDTO.User) error
	SetPassword(ctx context.Context, userID string, password string) error
	SaveUsersRole(ctx context.Context, role string, userID string) error
}

type keycloakService struct {
	client       *gocloak.GoCloak
	token        string
	realm        string
	clientSecret string
}

func NewKeycloakService(ctx context.Context) (KeycloakService, error) {
	client := gocloak.NewClient(os.Getenv("KEYCLOAK_URL"))

	var token *gocloak.JWT
	var err error

	for i := 0; i < 24; i++ {
		token, err = client.LoginAdmin(ctx, os.Getenv("KEYCLOAK_ADMIN"), os.Getenv("KEYCLOAK_PASSWORD"), os.Getenv("KEYCLOAK_REALM"))
		if err == nil {
			break
		}
		fmt.Println("Waiting for Keycloak to be ready...", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("Keycloak login failed: %w", err)
	}

	service := keycloakService{
		client: client,
		token:  token.AccessToken,
		realm:  os.Getenv("KEYCLOAK_REALM"),
	}

	keycloakClient, errTwo := service.KeycloakClient(ctx)
	if errTwo != nil {
		fmt.Println("errK", errTwo)
		return nil, errTwo
	}
	fmt.Println("id", keycloakClient)

	secret, errThree := service.GetClientSecret(ctx, keycloakClient)
	if errThree != nil {
		fmt.Println("errThreeK", errThree)
		return nil, errThree
	}

	service.clientSecret = secret

	return &service, nil
}

func (k *keycloakService) CreateUser(ctx context.Context, signUpReq entitiesDTO.SignupReq) (string, error) {
	fmt.Println("user", signUpReq)
	user := gocloak.User{
		Username:  gocloak.StringP(signUpReq.Username),
		Email:     gocloak.StringP(signUpReq.Email),
		FirstName: gocloak.StringP(signUpReq.FirstName),
		LastName:  gocloak.StringP(signUpReq.LastName),
		Enabled:   gocloak.BoolP(true),
	}

	fmt.Println("token", k.token)

	userID, err := k.client.CreateUser(ctx, k.token, k.realm, user)

	if err != nil {
		fmt.Println("err", err)
		return "", err
	}

	fmt.Println("userID", userID)

	errTwo := k.SetPassword(ctx, userID, signUpReq.Password)
	if errTwo != nil {
		return "", errTwo
	}

	errThree := k.SaveUsersRole(ctx, signUpReq.Role, userID)
	if errThree != nil {
		return "", errThree
	}

	return userID, nil
}

func (k *keycloakService) SetPassword(ctx context.Context, userID string, password string) error {
	errTwo := k.client.SetPassword(ctx, k.token, userID, k.realm, password, false)

	if errTwo != nil {
		return errTwo
	}

	return nil
}

func (k *keycloakService) LoginUser(ctx context.Context, loginReq entitiesDTO.LoginReq) (entitiesDTO.LoginResponse, error) {
	access, err := k.client.Login(ctx, os.Getenv("KEYCLOAK_CLIENT_ID"), k.clientSecret, k.realm, loginReq.Username, loginReq.Password)
	if err != nil {
		return entitiesDTO.LoginResponse{}, err
	}

	token := entitiesDTO.LoginResponse{
		AccessToken:  access.AccessToken,
		RefreshToken: access.RefreshToken,
	}

	return token, nil
}

func (k *keycloakService) KeycloakClient(ctx context.Context) (string, error) {
	clientID := os.Getenv("KEYCLOAK_CLIENT_ID")
	clients, err := k.client.GetClients(ctx, k.token, k.realm, gocloak.GetClientsParams{
		ClientID: &clientID,
	})
	if err != nil {
		return "", err
	}

	if len(clients) == 0 {

		newClient := gocloak.Client{
			ClientID:               gocloak.StringP(clientID),
			Name:                   gocloak.StringP(clientID),
			Enabled:                gocloak.BoolP(true),
			Protocol:               gocloak.StringP("openid-connect"),
			PublicClient:           gocloak.BoolP(false),
			ServiceAccountsEnabled: gocloak.BoolP(true),
			Attributes: &map[string]string{
				"access.token.lifespan":       "300",
				"client.session.idle.timeout": "1800",
			},
			StandardFlowEnabled:       gocloak.BoolP(true),
			DirectAccessGrantsEnabled: gocloak.BoolP(true),
		}
		fmt.Println("newClient", newClient)

		client, errTwo := k.client.CreateClient(ctx, k.token, k.realm, newClient)
		if errTwo != nil {
			fmt.Println("errCl", errTwo)
			return "", errTwo
		}

		time.Sleep(1 * time.Second)
		return client, nil
	}

	return *clients[0].ID, nil
}

func (k *keycloakService) GetClientSecret(ctx context.Context, internalID string) (string, error) {
	var secret *gocloak.CredentialRepresentation
	var err error
	for i := 0; i < 5; i++ {
		secret, err = k.client.GetClientSecret(ctx, k.token, k.realm, internalID)
		if err == nil && secret != nil && *secret.Value != "" {
			return *secret.Value, nil
		}
		time.Sleep(1 * time.Second)
	}
	return "", fmt.Errorf("client secret empty or unavailable: %w", err)
}

func (k *keycloakService) RefreshToken(ctx context.Context, refreshToken string) (entitiesDTO.LoginResponse, error) {
	access, err := k.client.RefreshToken(ctx, refreshToken, os.Getenv("KEYCLOAK_CLIENT_ID"), k.clientSecret, k.realm)
	if err != nil {
		return entitiesDTO.LoginResponse{}, fmt.Errorf("failed to refresh token: %w", err)
	}

	return entitiesDTO.LoginResponse{
		AccessToken:  access.AccessToken,
		RefreshToken: access.RefreshToken,
	}, nil
}

func (k *keycloakService) GetUser(ctx context.Context, userID string) (*gocloak.User, error) {
	user, err := k.client.GetUserByID(ctx, k.token, k.realm, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (k *keycloakService) DeleteUser(ctx context.Context, userID string) error {
	err := k.client.DeleteUser(ctx, k.token, k.realm, userID)
	if err != nil {
		return err
	}

	return nil
}

func (k *keycloakService) UpdateUser(ctx context.Context, user entitiesDTO.User) error {
	var updUser gocloak.User
	updUser.ID = gocloak.StringP(user.ID)

	if user.FirstName != "" {
		updUser.FirstName = gocloak.StringP(user.FirstName)
	}
	if user.LastName != "" {
		updUser.LastName = gocloak.StringP(user.LastName)
	}
	if user.Email != "" {
		updUser.Email = gocloak.StringP(user.Email)
	}

	fmt.Println("realm", k.realm)
	fmt.Println("token", k.token)
	err := k.client.UpdateUser(ctx, k.token, k.realm, updUser)
	if err != nil {
		fmt.Println("updErr", err)
		return err
	}

	if user.Password != "" {
		errTwo := k.SetPassword(ctx, user.ID, user.Password)
		if errTwo != nil {
			return errTwo
		}
	}

	return nil
}

func (k *keycloakService) GetSecret() string {
	return k.clientSecret
}

func (k *keycloakService) SaveUsersRole(ctx context.Context, role string, userID string) error {
	realmRole, err := k.client.GetRealmRole(ctx, k.token, k.realm, role)
	if err != nil {
		return err
	}

	errTwo := k.client.AddRealmRoleToUser(ctx, k.token, k.realm, userID, []gocloak.Role{{ID: realmRole.ID, Name: realmRole.Name}})
	if errTwo != nil {
		return errTwo
	}

	return nil
}
