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
import TunnelService from '../../modules/api/tunnel';
import AppPageContent from '../../modules/components/AppPageContent';
import AppPageActions from '../../modules/components/AppPageActions';
import { withKeycloak } from '@react-keycloak/web';
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

class TunnelAgent extends React.Component {
    constructor(props, context) {
        super(props, context);
        this.state = {
            tunnels: [],
            selected: [],
            order: 'asc',
            isLoading: false,
            showBackendError: false,
            backendErrorTitle: null,
            backendErrorMessage: null,
        };

        this.tunnelService = new TunnelService();
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
            const tunnels = await this.tunnelService.list(keycloak.token);
            this.setState({
                tunnels,
                isLoading: false,
            });
        } catch (error) {
            this.handleBackendError('Listing failed:', error.message);
        }
    };

    handleClick = (event, schedule) => {
        const { selected } = this.state;

        const selectedIndex = indexOf(selected, schedule.name);
        let newSelected = [];

        if (selectedIndex === -1) {
            newSelected = newSelected.concat(selected, schedule.name);
        } else if (selectedIndex === 0) {
            newSelected = newSelected.concat(selected.slice(1));
        } else if (selectedIndex === selected.length - 1) {
            newSelected = newSelected.concat(selected.slice(0, -1));
        } else if (selectedIndex > 0) {
            newSelected = newSelected.concat(
                selected.slice(0, selectedIndex),
                selected.slice(selectedIndex + 1)
            );
        }

        this.setState({ selected: newSelected });
    };

    handleSelectAllClick = (event, checked) => {
        if (checked) {
            this.setState({ selected: this.state.tunnels.map((p) => p.name) });
        } else {
            this.setState({ selected: [] });
        }
    };

    handleDeleteClick = async (event) => {
        const { keycloak } = this.props;
        try {
            const { selected } = this.state;
            if (selected.length > 0) {
                const promises = [];
                this.setState({ isLoading: true });
                selected.forEach((schedule) => {
                    promises.push(
                        this.scheduleService
                            .delete(schedule, keycloak.token)
                            .catch((error) => error)
                    );
                });
                const responses = await Promise.all(promises);
                const errorMessages = responses
                    .filter((response) => response instanceof Error)
                    .map((error) => error.message);
                if (errorMessages.length) {
                    throw Error(errorMessages.join('; '));
                }
                this.setState(
                    {
                        selected: [],
                        isLoading: false,
                    },
                    () => {
                        this.refreshList();
                    }
                );
            }
        } catch (error) {
            this.handleBackendError('Deletion failed:', error.message);
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
            tunnels,
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
                        onClick={this.handleClickRefresh}
                    >
                        <RefreshIcon className={classes.leftIcon} />
                        Refresh
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
                                            User
                                        </TableSortLabel>

                                    </Tooltip>
                                </TableCell>
                                <TableCell sortDirection={order}>
                                    Ip
                                </TableCell>
                                <TableCell sortDirection={order}>
                                    Region
                                </TableCell>
                                <TableCell sortDirection={order}>
                                    OS
                                </TableCell>
                                <TableCell sortDirection={order}>
                                    Started at
                                </TableCell>
                                <TableCell sortDirection={order}>
                                    Version
                                </TableCell>
                                <TableCell sortDirection={order}>
                                    Metadata
                                </TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {map(tunnels['data'], (tunnel) => {
                                const isSelected = indexOf(selected, tunnel.name) !== -1;
                                return (
                                    <TableRow
                                        hover
                                        role="checkbox"
                                        aria-checked={isSelected}
                                        tabIndex={-1}
                                        key={tunnel.name}
                                        selected={isSelected}
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
                                                {tunnel.ID}
                                            </span>
                                        </TableCell>
                                        <TableCell>
                                            <span
                                                className={classes.link}
                                            >
                                                {tunnel.userRemoteId}
                                            </span>
                                        </TableCell>
                                        <TableCell>
                                            <span
                                                className={classes.link}
                                            >
                                                {tunnel.ip}
                                            </span>
                                        </TableCell>
                                        <TableCell>
                                            <span
                                                className={classes.link}
                                            >
                                                {tunnel.region}
                                            </span>
                                        </TableCell>
                                        <TableCell>
                                            <span
                                                className={classes.link}
                                            >
                                                {tunnel.os}
                                            </span>
                                        </TableCell>
                                        <TableCell>
                                            <span
                                                className={classes.link}
                                            >
                                                {tunnel.startedAt
                                                }
                                            </span>
                                        </TableCell>

                                        <TableCell>
                                            <span
                                                className={classes.link}
                                            >
                                                {tunnel.version
                                                }
                                            </span>
                                        </TableCell>

                                        <TableCell>
                                            <span
                                                className={classes.link}
                                            >
                                                {tunnel.metadata
                                                }
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
)(TunnelAgent);