import React from 'react';
import PropTypes from 'prop-types';

// Recompose
import { compose } from 'react-recompose';

// Router
import { Switch, Route, Redirect } from 'react-router-dom';

// Material-UI
import { withStyles } from '@material-ui/core/styles';

// Project
import withRoot from '../withRoot';
import withProps from '../withProps';
import AppFrame from '../modules/components/AppFrameFunction';

// Project Views
import NotFound from './NotFound/NotFound';

import ScheduleList from './Schedule/ScheduleList';
import ScheduleCreate from './Schedule/ScheduleCreate';
import ScheduleEdit from './Schedule/ScheduleEdit';

import Policy from './Policy/Policy';
import PolicyList from './Policy/PolicyList';
import { ReactKeycloakProvider } from '@react-keycloak/web';
import keycloak from '../keycloak';
const styles = (theme) => ({
  '@global': {
    'html, body, #root': {
      height: '100%',
    },
  },
  root: {},
});

class Index extends React.Component {
  constructor(props, context) {
    super(props, context);
    const refreshToken = localStorage.getItem('refreshToken');
    const setTokens = (token, idToken, refreshToken) => {
      localStorage.setItem('token', token);
      localStorage.setItem('refreshToken', refreshToken);
      localStorage.setItem('idToken', idToken);
    };
    this.state = {
      refreshToken: refreshToken,
      setTokens: setTokens,
    };
  }
  render() {
    const { classes } = this.props;
    const token = localStorage.getItem('token');
    const { refreshToken, setTokens } = this.state;

    return (
      <ReactKeycloakProvider
        authClient={keycloak}
        onTokens={(tokens) =>
          setTokens(
            tokens.token ?? '',
            tokens.idToken ?? '',
            tokens.refreshToken ?? ''
          )
        }
        initOptions={{ onLoad: 'login-required', token, refreshToken }}
      >
        <AppFrame className={classes.root}>
          <Switch>
            <Route
              exact
              path="/"
              render={() => <Redirect to="/schedules/browser" />}
            />
            <Route exact path="/schedules/create" component={ScheduleCreate} />
            <Route exact path="/schedules/browser" component={ScheduleList} />
            <Route
              exact
              path="/schedules/browser/:schedule"
              component={ScheduleEdit}
            />

            <Route
              exact
              path="/policies/create"
              component={withProps(Policy, { edit: false })}
            />
            <Route
              exact
              path="/policies/browser/:policy"
              component={withProps(Policy, { edit: true })}
            />
            <Route exact path="/policies/browser" component={PolicyList} />
            <Route component={NotFound} />
          </Switch>
        </AppFrame>
      </ReactKeycloakProvider>
    );
  }
}

Index.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default compose(withRoot, withStyles(styles))(Index);
