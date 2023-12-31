import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';

// Recompose
import { compose } from 'react-recompose';

// Router
import { withRouter } from 'react-router-dom';

import { withStyles } from '@material-ui/core/styles';

import PolicyIcon from '@material-ui/icons/LibraryBooks';
import SupervisedUserCircle from '@material-ui/icons/SupervisedUserCircle';
import ScheduleIcon from '@material-ui/icons/Schedule';

// Lodash
import find from 'lodash/find';
import startsWith from 'lodash/startsWith';

// Project
import { withKeycloak } from '@react-keycloak/web';
import ResponsiveAppBar from './AppBar';

import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import Drawer from '@material-ui/core/Drawer';
import map from 'lodash/map';
import { ListItemButton } from '@mui/material';
import Divider from '@material-ui/core/Divider';

const drawerWidth = 280;

const links = [
  // {
  //   primary: 'Schedules',
  //   path: '/schedules/browser',
  //   icon: <ScheduleIcon />,
  // },
  // {
  //   primary: 'Policies',
  //   path: '/policies/browser',
  //   icon: <PolicyIcon />,
  // },

  {
    primary: 'Auth Token',
    path: '/auth/',
    icon: <SupervisedUserCircle />,
  },
  {
    primary: 'Users',
    path: '/users/',
    icon: <SupervisedUserCircle />,
  },
  {
    primary: 'Tunnels',
    path: '/tunnels/agent',
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
  sidebar: {
    width: drawerWidth,
    height: '100%',
    backgroundColor: '#21214c'
  },
  content: {
    backgroundColor: theme.palette.grey,
    width: '100%',
    [theme.breakpoints.up('sm')]: {
      height: 'calc(100% - 64px)',
    },
    marginLeft: drawerWidth,
    [theme.breakpoints.up('md')]: {
      width: `calc(100% - ${drawerWidth}px)`,
    },
  },
});

function AppFrame(props) {
  const [title, setTitle] = useState('');
  const keycloakInitialized = props.keycloakInitialized;
  const [drawerOpen, setDrawerOpen] = React.useState(false);
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
  const toggleDrawer = (open) => (event) => {
    if (
      event.type === 'keydown' &&
      (event.key === 'Tab' || event.key === 'Shift')
    ) {
      return;
    }
    setDrawerOpen(open);
  };
  const drawer = (
    <div className={classes.sidebar}>
      <List dense disablePadding>
        {map(links, (page, _) => (
          <ListItem
            disableGutters
            key={page.primary}
            style={{
              paddingLeft: 0,
              paddingRight: 0,
            }}
            button onClick={handleClickLink(page)}
          >
            <ListItemButton disablePadding style={{
              color: 'white',
              backgroundColor: title === page.primary && '#1853db',
            }}>
              <ListItemText
             
                primary={
                  <div style={{
                    fontWeight: title === page.primary ? 500 : 400
                  }}>
                    {page.primary}
                  </div>
                }
              />
            </ListItemButton>

          </ListItem>
        ))}
      </List>
      {/* <Divider style={{ backgroundColor: 'white' }} /> */}
    </div>
  );

  if (keycloakInitialized) {
    return (
      <div className={classes.root}>
        <Drawer
          anchor="left"
          open={true}
          onClose={toggleDrawer(false)}
          variant='permanent'
        >
          {drawer}
        </Drawer>

        <main className={classes.content}>
          <ResponsiveAppBar
            title={title}
            onClickLink={handleClickLink}
            links={links}
          />
          <div style={{ maxWidth: 1200, margin: '0 auto' }}>
            {children}
          </div>
        </main>
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
