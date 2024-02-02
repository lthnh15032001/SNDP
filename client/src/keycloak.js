import Keycloak from 'keycloak-js';
const keycloak = new Keycloak({
  url: 'http://localhost:8080',
  realm: 'SNDP',
  clientId: 'sndp',
});

export default keycloak;
