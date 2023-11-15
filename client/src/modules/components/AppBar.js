import * as React from 'react';
import classNames from 'classnames';

import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import Menu from '@mui/material/Menu';
import MenuIcon from '@mui/icons-material/Menu';
import Container from '@mui/material/Container';
import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import Tooltip from '@mui/material/Tooltip';
import MenuItem from '@mui/material/MenuItem';
import FilterDrama from '@mui/icons-material/FilterDrama';
import { useKeycloak } from '@react-keycloak/web';
import Divider from '@material-ui/core/Divider';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import Hidden from '@material-ui/core/Hidden';
import List from '@material-ui/core/List';
import map from 'lodash/map';
import startsWith from 'lodash/startsWith';
import Drawer from '@material-ui/core/Drawer';
import { ListItemButton } from '@mui/material';
import Brightness4Icon from '@mui/icons-material/Brightness4';
import Brightness7Icon from '@mui/icons-material/Brightness7';
const ColorModeContext = React.createContext({ toggleColorMode: () => {} });

function ResponsiveAppBar(props) {
  const [anchorElNav, setAnchorElNav] = React.useState(null);
  const [anchorElUser, setAnchorElUser] = React.useState(null);
  const [drawerOpen, setDrawerOpen] = React.useState(false);
  const colorMode = React.useContext(ColorModeContext);
  const { keycloak } = useKeycloak();
  
  const links = props.links;
  const handleOpenNavMenu = (event) => {
    setAnchorElNav(event.currentTarget);
  };
  const handleOpenUserMenu = (event) => {
    setAnchorElUser(event.currentTarget);
  };

  const handleCloseNavMenu = () => {
    setAnchorElNav(null);
  };

  const handleCloseUserMenu = () => {
    setAnchorElUser(null);
  };
  const drawer = (
    <div>
      <List dense disablePadding>
        {map(links, (page, _) => (
          <ListItem key={page.primary} button onClick={props.onClickLink(page)}>
            <ListItemButton>
              <ListItemIcon>{page.icon}</ListItemIcon>
              <ListItemText primary={page.primary} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
      <Divider />
    </div>
  );
  const toggleDrawer = (open) => (event) => {
    if (
      event.type === 'keydown' &&
      (event.key === 'Tab' || event.key === 'Shift')
    ) {
      return;
    }
    setDrawerOpen(open);
  };
  return (
    <div>
      <AppBar position="absolute" color="primary" elevation={2} square>
        <Container maxWidth="false">
          <Toolbar disableGutters>
            <IconButton color="inherit" aria-label="open drawer">
              <MenuIcon onClick={toggleDrawer(true)} />
              <Drawer
                anchor="left"
                open={drawerOpen}
                onClose={toggleDrawer(false)}
              >
                {drawer}
              </Drawer>
            </IconButton>
            <FilterDrama sx={{ display: { xs: 'none', md: 'flex' }, mr: 1 }} />
            <Typography
              variant="h6"
              noWrap
              component="a"
              href="/"
              sx={{
                mr: 2,
                display: { xs: 'none', md: 'flex' },
                fontFamily: 'monospace',
                fontWeight: 700,
                letterSpacing: '.3rem',
                color: 'inherit',
                textDecoration: 'none',
              }}
            >
              SNDP
            </Typography>

            <Box sx={{ flexGrow: 1, display: { xs: 'flex', md: 'none' } }}>
              <IconButton
                size="large"
                aria-label="account of current user"
                aria-controls="menu-appbar"
                aria-haspopup="true"
                onClick={handleOpenNavMenu}
                color="inherit"
              ></IconButton>
              <Menu
                id="menu-appbar"
                anchorEl={anchorElNav}
                anchorOrigin={{
                  vertical: 'bottom',
                  horizontal: 'left',
                }}
                keepMounted
                transformOrigin={{
                  vertical: 'top',
                  horizontal: 'left',
                }}
                open={Boolean(anchorElNav)}
                onClose={handleCloseNavMenu}
                sx={{
                  display: { xs: 'block', md: 'none' },
                }}
              >
                {links.map((page) => (
                  <MenuItem
                    key={page.primary}
                    onClick={props.onClickLink(page)}
                  >
                    <Typography textAlign="center">{page.primary}</Typography>
                  </MenuItem>
                ))}
              </Menu>
            </Box>
            <FilterDrama sx={{ display: { xs: 'flex', md: 'none' }, mr: 1 }} />
            <Typography
              variant="h5"
              noWrap
              component="a"
              href="/"
              sx={{
                mr: 2,
                display: { xs: 'flex', md: 'none' },
                flexGrow: 1,
                fontFamily: 'monospace',
                fontWeight: 700,
                letterSpacing: '.3rem',
                color: 'inherit',
                textDecoration: 'none',
              }}
            >
              LOGO
            </Typography>
            <Box sx={{ flexGrow: 1, display: { xs: 'none', md: 'flex' } }}>
              {links.map((page) => (
                <Button
                  key={page.primary}
                  onClick={props.onClickLink(page)}
                  sx={{ my: 2, color: 'white', display: 'block' }}
                >
                  {page.primary}
                </Button>
              ))}
            </Box>
            <Box sx={{ flexGrow: 0 }}>
              <Tooltip title="Open settings">
                <IconButton onClick={handleOpenUserMenu} sx={{ p: 0 }}>
                  {/* only when authenticated we have the keycloak.tokenParsed object */}
                  {keycloak.authenticated ? (
                    <Avatar
                      alt="user-avatar"
                      src={`https://avatars.githubusercontent.com/${keycloak.tokenParsed.preferred_username}`}
                    />
                  ) : null}
                </IconButton>
              </Tooltip>
              <Menu
                sx={{ mt: '45px' }}
                id="menu-appbar"
                anchorEl={anchorElUser}
                anchorOrigin={{
                  vertical: 'top',
                  horizontal: 'right',
                }}
                keepMounted
                transformOrigin={{
                  vertical: 'top',
                  horizontal: 'right',
                }}
                open={Boolean(anchorElUser)}
                onClose={handleCloseUserMenu}
              >
                <MenuItem key="Logout" onClick={handleCloseUserMenu}>
                  <Typography textAlign="center">
                    {!keycloak.authenticated && (
                      <div
                        type="button"
                        className="text-blue-800"
                        onClick={() => keycloak.login()}
                      >
                        Login
                      </div>
                    )}

                    {!!keycloak.authenticated && (
                      <div
                        type="button"
                        className="text-blue-800"
                        onClick={() => keycloak.logout()}
                      >
                        Logout
                      </div>
                    )}
                  </Typography>
                </MenuItem>
              </Menu>
            </Box>
          </Toolbar>
        </Container>
      </AppBar>
    </div>
  );
}
export default ResponsiveAppBar;
