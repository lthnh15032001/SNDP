import React from 'react';
import PropTypes from 'prop-types';

// Material UI
import { withStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import FormGroup from '@material-ui/core/FormGroup';
import Button from '@material-ui/core/Button';
import IconButton from '@material-ui/core/IconButton';
import ClearIcon from '@material-ui/icons/Clear';

// Lodash
import map from 'lodash/map';

const TEXT_FIELD_WIDTH = 250;

const styles = (theme) => ({
  root: {
    marginBottom: theme.spacing(3),
  },
  textField: {
    width: TEXT_FIELD_WIDTH,
    marginRight: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },
  iconButton: {
    width: 32,
    height: 32,
  },
  addButton: {
    width: TEXT_FIELD_WIDTH * 2 + theme.spacing(1),
  },
  sizeSmallButton: {
    padding: 0,
    minHeight: 24,
  },
});

class PolicyProjects extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.state = {
      projects: [
        {
          name: '',
          credentialRef: '',
        },
      ],
    };
  }

  componentDidMount() {
    if (this.props.projects && this.props.projects.length > 0) {
      let projects = [];
      this.props.projects.forEach((project) => {
        projects.push({
          name: project.name,
          credentialRef: project.credentialRef,
        });
      });
      this.setState({
        projects,
      });
    }
  }

  publishChanges = (shouldUpdateErrors) => {
    const projects = map(this.state.projects, (project) => ({
      name: project.name,
      credentialRef: project.credentialRef,
    }));
    this.props.onChange(projects, shouldUpdateErrors);
  };

  handleChange = (index, name) => (event) => {
    const projects = this.state.projects.slice();
    projects[index][name] = event.target.value;
    this.setState({ projects }, () => this.publishChanges(false));
  };

  handleClearProject = (index) => (event) => {
    const projects = this.state.projects.slice();
    if (projects.length > 1) {
      projects.splice(index, 1);
      this.setState({ projects }, () => this.publishChanges(true));
    }
  };

  handleAddProject = (event) => {
    const projects = this.state.projects.slice();
    projects.push({
      name: '',
      credentialRef: '',
    });
    this.setState({ projects }, () => this.publishChanges(false));
  };

  render() {
    const { classes, error } = this.props;
    const { projects } = this.state;
    return (
      <div className={classes.root}>
        {map(projects, (project, index) => (
          <FormGroup row key={index}>
            <TextField
              id="project-name-value"
              error={error[index] && error[index][0]}
              helperText="Name of the project"
              placeholder="Name"
              className={classes.textField}
              value={projects[index].name}
              onChange={this.handleChange(index, 'name')}
              margin="none"
            />
            <TextField
              id="project-credential-value"
              error={error[index] && error[index][1]}
              helperText="Credential reference to access this project"
              placeholder="Credential"
              className={classes.textField}
              value={projects[index].credentialRef}
              onChange={this.handleChange(index, 'credentialRef')}
              margin="none"
            />

            {projects.length > 1 && (
              <IconButton
                className={classes.iconButton}
                aria-label="Clear"
                onClick={this.handleClearProject(index)}
                classes={{
                  root: classes.iconButton,
                }}
              >
                <ClearIcon />
              </IconButton>
            )}
          </FormGroup>
        ))}

        {projects.length < 7 && (
          <Button
            variant="outlined"
            color="primary"
            size="small"
            className={classes.addButton}
            onClick={this.handleAddProject}
            classes={{
              sizeSmall: classes.sizeSmallButton,
            }}
          >
            Add project
          </Button>
        )}
      </div>
    );
  }
}
PolicyProjects.propTypes = {
  classes: PropTypes.object.isRequired,
  onChange: PropTypes.func.isRequired,
  error: PropTypes.array.isRequired,
};

export default withStyles(styles)(PolicyProjects);
