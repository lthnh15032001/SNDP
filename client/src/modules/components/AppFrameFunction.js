import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';

// Recompose
import { compose } from 'react-recompose';

// Router
import { withRouter } from 'react-router-dom';

import { withStyles } from '@material-ui/core/styles';

import PolicyIcon from '@material-ui/icons/LibraryBooks';
import ScheduleIcon from '@material-ui/icons/Schedule';

// Lodash
import find from 'lodash/find';
import startsWith from 'lodash/startsWith';

// Project
import { withKeycloak } from '@react-keycloak/web';
import ResponsiveAppBar from './AppBar';

const drawerWidth = 210;

const links = [
  {
    primary: 'Schedules',
    path: '/schedules/browser',
    icon: <ScheduleIcon />,
  },
  {
    primary: 'Policies',
    path: '/policies/browser',
    icon: <PolicyIcon />,
  },
];

const styles = (theme) => ({
  root: {
    width: '100%',
    height: '100%',
    zIndex: 1,
    overflow: 'scroll',
    display: 'flex',
  },
  appBar: {
    marginLeft: drawerWidth,
    [theme.breakpoints.up('md')]: {
      width: `calc(100% - ${drawerWidth}px)`,
    },
  },
  navIconHide: {
    [theme.breakpoints.up('md')]: {
      display: 'none',
    },
  },
  drawerHeader: {
    display: 'flex',
    alignItems: 'center',
    ...theme.mixins.toolbar,
  },
  drawerPaper: {
    width: drawerWidth,
    [theme.breakpoints.up('md')]: {
      position: 'relative',
      height: '100%',
    },
  },
  drawerDocked: {
    height: '100%',
  },
  content: {
    backgroundColor: theme.palette.background.default,
    width: '100%',
    height: 'calc(100% - 56px)',
    marginTop: 56,
    [theme.breakpoints.up('sm')]: {
      height: 'calc(100% - 64px)',
      marginTop: 64,
    },
  },
});

function AppFrame(props) {
  const [title, setTitle] = useState('');
  const keycloakInitialized = props.keycloakInitialized;
  useEffect(() => {
    const currentLink = find(links, (link) =>
      startsWith(props.history.location.pathname, link.path)
    );
    if (currentLink) {
      setTitle(currentLink.primary);
    }
  }, [props.history]);

  const handleClickLink = (link) => (event) => {
    props.history.push(link.path);
    setTitle(link.primary);
  };

  const classes = props.classes;
  const children = props.children;
  return (
    <div className={classes.root}>
      <ResponsiveAppBar
        title={title}
        onClickLink={handleClickLink}
        links={links}
      />
      <main className={classes.content}>{children}</main>
    </div>
  );
  if (keycloakInitialized) {
    return (
      <div className={classes.root}>
        <ResponsiveAppBar
          title={title}
          onClickLink={handleClickLink}
          links={links}
        />
        <main className={classes.content}>{children}</main>
      </div>
    );
  } else {
    return <div>Loading...</div>;
  }
}

AppFrame.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default compose(withRouter, withStyles(styles), withKeycloak)(AppFrame);
