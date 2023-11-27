import React from 'react';

// Recompose
import { compose } from 'react-recompose';

// Router
import { withRouter } from 'react-router-dom';

// Material UI
import { withStyles } from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TableSortLabel from '@material-ui/core/TableSortLabel';
import Tooltip from '@material-ui/core/Tooltip';
import Button from '@material-ui/core/Button';
import Checkbox from '@material-ui/core/Checkbox';
import AddIcon from '@material-ui/icons/Add';
import RefreshIcon from '@material-ui/icons/Refresh';
import EditIcon from '@material-ui/icons/Edit';
import DeleteIcon from '@material-ui/icons/Delete';

// Lodash
import map from 'lodash/map';
import indexOf from 'lodash/indexOf';

// Project
import AppPageContent from '../../modules/components/AppPageContent';
import AppPageActions from '../../modules/components/AppPageActions';
import { withKeycloak } from '@react-keycloak/web';
import UserService from '../../modules/api/user';
const styles = (theme) => ({
    root: {
        height: '100%',
    },
    button: {
        marginRight: theme.spacing(2),
    },
    leftIcon: {
        marginRight: theme.spacing(1),
    },
    link: {
        '&:hover': {
            textDecoration: 'underline',
            cursor: 'pointer',
        },
    },
    checkboxCell: {
        width: 48,
    },
});

class User extends React.Component {
    constructor(props, context) {
        super(props, context);
        this.state = {
            users: [],
            selected: [],
            order: 'asc',
            isLoading: false,
            showBackendError: false,
            backendErrorTitle: null,
            backendErrorMessage: null,
        };

        this.userService = new UserService();
    }

    componentDidMount() {
        this.refreshList();
    }

    handleRequestSort = (event) => {
        this.setState((prevState, props) => {
            let order = 'desc';
            if (prevState.order === 'desc') {
                order = 'asc';
            }

            const schedules =
                order === 'desc'
                    ? prevState.schedules.sort((a, b) =>
                        b.displayname || b.name < a.displayname || a.name ? -1 : 1
                    )
                    : prevState.schedules.sort((a, b) =>
                        a.displayname || a.name < b.displayname || b.name ? -1 : 1
                    );

            return {
                schedules,
                order,
            };
        });
    };

    handleClickNavigate = (path) => (event) => {
        const { history } = this.props;
        history.push(path);
    };

    handleClickRefresh = (event) => {
        this.refreshList();
    };

    refreshList = async () => {
        const { keycloak } = this.props;
        this.setState({ isLoading: true });
        try {
            const users = await this.userService.list(keycloak.token);
            this.setState({
                users,
                isLoading: false,
            });
        } catch (error) {
            this.handleBackendError('Listing failed:', error.message);
        }
    };


    handleBackendError = (title, message) => {
        this.setState({
            backendErrorTitle: title,
            backendErrorMessage: message,
            showBackendError: true,
            isLoading: false,
        });
    };

    handleErrorClose = () => {
        this.setState({
            showBackendError: false,
            isLoading: false,
        });
    };

    render() {
        const { classes } = this.props;
        const {
            users,
            selected,
            order,
            isLoading,
            backendErrorTitle,
            backendErrorMessage,
            showBackendError,
        } = this.state;


        return (
            <div className={classes.root}>
                <AppPageActions>
                    <Button
                        className={classes.button}
                        color="primary"
                        size="small"
                        onClick={this.handleClickNavigate(`/users/create`)}
                    >
                        <AddIcon className={classes.leftIcon} />
                        Add User
                    </Button>
                    <Button
                        className={classes.button}
                        color="primary"
                        size="small"
                        onClick={this.handleClickRefresh}
                    >
                        <RefreshIcon className={classes.leftIcon} />
                        Refresh
                    </Button>
                    <Button
                        className={classes.button}
                        color="primary"
                        size="small"
                        disabled={selected.length !== 1}
                        onClick={this.handleClickNavigate(
                            `/schedules/browser/${selected[0]}`
                        )}
                    >
                        <EditIcon className={classes.leftIcon} />
                        Edit
                    </Button>
                    <Button
                        className={classes.button}
                        color="primary"
                        size="small"
                        disabled={selected.length < 1}
                        onClick={this.handleDeleteClick}
                    >
                        <DeleteIcon className={classes.leftIcon} />
                        Delete
                    </Button>
                </AppPageActions>

                <AppPageContent
                    showBackendError={showBackendError}
                    backendErrorTitle={backendErrorTitle}
                    backendErrorMessage={backendErrorMessage}
                    onBackendErrorClose={this.handleErrorClose}
                    showLoadingSpinner={isLoading}
                >
                    <Table className={classes.table}>
                        <TableHead>
                            <TableRow>
                                {/* <TableCell padding="none" className={classes.checkboxCell}>
                                    <Checkbox
                                        indeterminate={numSelected > 0 && numSelected < rowCount}
                                        checked={rowCount > 0 && numSelected === rowCount}
                                        onChange={this.handleSelectAllClick}
                                    />
                                </TableCell> */}
                                <TableCell sortDirection={order}>
                                    <Tooltip
                                        title={order === 'desc' ? 'descending' : 'ascending'}
                                        placement="bottom-start"
                                        enterDelay={500}
                                    >
                                        <TableSortLabel
                                            active
                                            direction={order}
                                        // onClick={this.handleRequestSort}
                                        >
                                            ID
                                        </TableSortLabel>

                                    </Tooltip>
                                </TableCell>
                                <TableCell sortDirection={order}>
                                    Username
                                </TableCell>
                                <TableCell sortDirection={order}>
                                    Policy
                                </TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {map(users['data'], (user) => {
                                return (
                                    <TableRow
                                        hover
                                        role="checkbox"
                                        aria-checked={false}
                                        tabIndex={-1}
                                        key={user.ID}
                                        selected={false}
                                    >
                                        {/* <TableCell padding="none" className={classes.checkboxCell}>
                                            <Checkbox
                                                checked={isSelected}
                                                onClick={(event) => this.handleClick(event, schedule)}
                                            />
                                        </TableCell> */}

                                        <TableCell>
                                            <span
                                                // onClick={this.handleClickNavigate(
                                                //     `/schedules/browser/${tunnel.name}`
                                                // )}
                                                className={classes.link}
                                            >
                                                {user.ID}
                                            </span>
                                        </TableCell>
                                        <TableCell>
                                            <span
                                                className={classes.link}
                                            >
                                                {user.username}
                                            </span>
                                        </TableCell>
                                        <TableCell>
                                            <span
                                                className={classes.link}
                                            >
                                                {JSON.stringify( user.userRemotePolicy)}
                                            </span>
                                        </TableCell> 
                                    </TableRow>
                                );
                            })}
                        </TableBody>
                    </Table>
                </AppPageContent>
            </div>
        );
    }
}

export default compose(
    withRouter,
    withStyles(styles),
    withKeycloak
)(User);
