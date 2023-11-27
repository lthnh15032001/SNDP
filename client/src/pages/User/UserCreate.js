import React from 'react';

// Material UI
import { withStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';
import IconButton from '@material-ui/core/IconButton';
import TextField from '@material-ui/core/TextField';

// Project
import ScheduleTimeTable from '../../modules/components/ScheduleTimeTable';
import ScheduleTimeZone from '../../modules/components/ScheduleTimeZone';
import AppPageContent from '../../modules/components/AppPageContent';
import AppPageActions from '../../modules/components/AppPageActions';
import UserService from '../../modules/api/user';
import { getDefaultSchedule } from '../../modules/utils/schedule';
import { compose } from 'react-recompose';
import { withKeycloak } from '@react-keycloak/web';
import map from 'lodash/map';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import DeleteForeverIcon from '@mui/icons-material/DeleteForever';

const styles = (theme) => ({
    root: {
        height: '100%',
    },
    button: {
        marginRight: theme.spacing(2),
    },
    textField: {
        minWidth: 250,
        marginBottom: theme.spacing(3),
        marginRight: theme.spacing(2),
    },
});

class UserCreate extends React.Component {
    constructor(props, context) {
        super(props, context);
        this.state = {
            user: {
                username: "",
                password: ""
            },
            policy: [
                {
                    txt: ''
                }
            ],
            nameError: false,
            timezones: [],
            isLoading: false,
            showBackendError: false,
            backendErrorTitle: null,
            backendErrorMessage: null,
            exitPage: null,
            keycloak: this.props.keycloak,
        };
        this.userService = new UserService();
    }

    async componentDidMount() {
        const { keycloak } = this.props;
        try {
            // this.setState({ isLoading: true });
            // const response = await this.scheduleService.timezones(keycloak.token);
            // this.setState({
            //   timezones: response.Timezones,
            //   isLoading: false,
            // });
        } catch (error) {
            this.handleBackendError(
                'Loading timezones failed:',
                error.message,
                '/schedules/browser'
            );
        }
    }

    handleChange = (name) => (event) => {
        const { user } = this.state;
        user[name] = event.target.value;
        this.setState({ user });
    };
    handleCreate = async (event) => {
        try {
            const { history, keycloak } = this.props;
            const { user, policy } = this.state;
            const nameRe = /^[a-zA-Z][\w-]*[a-zA-Z0-9]$/;
            if (!nameRe.test(user.username)) {
                this.setState({
                    nameError: true,
                });
                return;
            }
            this.setState({ isLoading: true });
            const userCreateDTO = {
                username: user.username,
                password: user.password,
                userRemotePolicy: policy.map(x => x['txt'])
            }
            await this.userService.create(userCreateDTO, keycloak.token);
            this.setState({ isLoading: false });
            history.push('/users');
        } catch (error) {
            this.handleBackendError('Saving failed:', error.message);
        }
    };

    handleRequestCancel = (event) => {
        const { history } = this.props;
        history.goBack();
    };

    handleBackendError = (title, message, exitPage) => {
        this.setState({
            backendErrorTitle: title,
            backendErrorMessage: message,
            showBackendError: true,
            isLoading: false,
            exitPage,
        });
    };

    handleErrorClose = () => {
        const { history } = this.props;
        const { exitPage } = this.state;
        this.setState({
            showBackendError: false,
            isLoading: false,
        });
        if (exitPage) {
            history.push(exitPage);
        }
    };
    handleAddMultiplePolicy = () => {
        const { policy } = this.state
        this.setState({
            policy: [
                {
                    txt: ''
                },
                ...policy
            ]
        })
    }

    handleDeletePolicy = (index) => {
        const { policy } = this.state
        policy.splice(index, 1)
        this.setState({
            policy: policy
        })
    }


    handlePolicyChange = (index, value) => {
        const { policy } = this.state
        policy[index].txt = value
        this.setState({
            policy: policy
        })
    }
    render() {
        const { classes } = this.props;
        const {
            user,
            policy,
            nameError,
            isLoading,
            showBackendError,
            backendErrorTitle,
            backendErrorMessage,
        } = this.state;
        return (
            <div className={classes.root}>
                <AppPageActions>
                    <IconButton
                        color="primary"
                        aria-label="Back"
                        onClick={this.handleRequestCancel}
                    >
                        <ArrowBackIcon />
                    </IconButton>
                    <Typography variant="subtitle1" color="primary">
                        Create User
                    </Typography>
                </AppPageActions>

                <AppPageContent
                    showBackendError={showBackendError}
                    backendErrorTitle={backendErrorTitle}
                    backendErrorMessage={backendErrorMessage}
                    onBackendErrorClose={this.handleErrorClose}
                    showLoadingSpinner={isLoading}
                >
                    <TextField
                        id="username"
                        error={nameError}
                        helperText="Required. server will authenticate incoming connections by validating this username along with password."
                        placeholder="username"
                        className={classes.textField}
                        value={user.username}
                        onChange={this.handleChange('username')}
                        margin="none"
                        autoFocus
                    />

                    <TextField
                        id="password"
                        type="password"
                        helperText="Required. Password for authenticating user"
                        placeholder="Password"
                        className={classes.textField}
                        value={user.password}
                        onChange={this.handleChange('password')}
                        margin="none"
                    />
                    <div className={classes.row}>
                        {
                            map(policy, (p, i) => {
                                return <div key={i.toString()}>
                                    <TextField
                                        id={`policy#${policy.length - i}`}
                                        type="text"
                                        helperText="Required. Policy of authenticated user's client. Ex: R:0.0.0.0:7000"
                                        placeholder={`policy#${policy.length - i - 1}`}
                                        className={classes.textField}
                                        value={p.txt}
                                        onChange={(e) => this.handlePolicyChange(i, e.target.value)}
                                        margin="none"
                                    />
                                    {i !== policy.length - 1 && <DeleteForeverIcon
                                        style={{ cursor: 'pointer' }}
                                        onClick={() => this.handleDeletePolicy(i)}
                                    />}
                                    {
                                        i === 0 && <AddCircleOutlineIcon
                                            style={{ cursor: 'pointer' }}
                                            onClick={() => this.handleAddMultiplePolicy()}
                                        />
                                    }
                                </div>
                            })
                        }

                    </div>



                    <Button
                        className={classes.button}
                        variant="contained"
                        color="primary"
                        size="small"
                        onClick={this.handleCreate}
                    >
                        Create
                    </Button>
                    <Button
                        className={classes.button}
                        variant="outlined"
                        color="primary"
                        size="small"
                        onClick={this.handleRequestCancel}
                    >
                        Cancel
                    </Button>
                </AppPageContent>
            </div>
        );
    }
}

export default compose(withStyles(styles), withKeycloak)(UserCreate);
