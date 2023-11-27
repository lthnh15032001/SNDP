package middlewares

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/lthnh15032001/ngrok-impl/internal/api/config"
	"github.com/lthnh15032001/ngrok-impl/internal/constants"

	"github.com/gin-gonic/gin"
)

type Claims struct {
	ResourceAccess client `json:"resource_access,omitempty"`
	JTI            string `json:"jti,omitempty"`
	Email          string `json:"email"`
	Sub            string `json:"sub"`
}

type client struct {
	SNDPServiceClient clientRoles `json:"sndp,omitempty"`
}

type clientRoles struct {
	Roles []string `json:"roles,omitempty"`
}

func authorisationFailed(message string, c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.AbortWithStatus(401)
	c.AbortWithStatusJSON(401, gin.H{
		"message": message,
	})
}

func AuthMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {

		config := config.GetConfig()
		issuerUrl := config.GetString(constants.ENV_OIDC_ISSUER_URL)
		clientID := config.GetString(constants.ENV_OIDC_CLIENT_ID)
		// serverSession := sessions.Default(c)
		rawAccessToken := c.Request.Header.Get("Authorization")
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{
			Timeout:   time.Duration(6000) * time.Second,
			Transport: tr,
		}
		ctx := oidc.ClientContext(c, client)
		provider, err := oidc.NewProvider(ctx, issuerUrl)
		if err != nil {
			authorisationFailed("authorisation failed while getting the provider: "+err.Error(), c)
			return
		}

		oidcConfig := &oidc.Config{
			ClientID: clientID,
			//skip check aud in jwt payload since in keycloak it's always `account`
			SkipClientIDCheck: true,
		}
		verifier := provider.Verifier(oidcConfig)
		idToken, err := verifier.Verify(ctx, rawAccessToken)
		if err != nil {
			authorisationFailed("authorisation failed while verifying the token: "+err.Error(), c)
			return
		}

		var IDTokenClaims Claims
		if err := idToken.Claims(&IDTokenClaims); err != nil {
			authorisationFailed("claims : "+err.Error(), c)
			return
		}

		//checking the roles
		user_access_roles := IDTokenClaims.ResourceAccess.SNDPServiceClient.Roles
		subject := IDTokenClaims.Sub
		email := IDTokenClaims.Email
		// if serverSession.Get("email") != email {
		// 	serverSession.Set("email", email)
		// }
		// if serverSession.Get("userId") != email {
		// 	serverSession.Set("userId", subject)
		// }
		// serverSession.Save()
		c.Set("email", email)
		c.Set("userid", subject)
		for _, b := range user_access_roles {
			if b == role {
				c.Next()
				return
			}
		}
		// serverSession.Clear()

		authorisationFailed("user not allowed to access this api", c)
	}
}
