import React from 'react';
import { NavLink } from 'react-router-dom';
import Button from '@material-ui/core/Button';
import { makeStyles } from '@material-ui/core/styles';


const useStyles = makeStyles((theme) => ({
  root: {
    '& > *': {
      margin: theme.spacing(1),
    },
  },
}));
 
const Navigation = () => {
    const classes = useStyles();
    return (
        <nav>
            <ul className={classes.root}>
                <Button variant="contained" color="primary">
                    <li><NavLink exact activeClassName="current" to="/">Home</NavLink></li>
                </Button>
                <Button variant="contained" color="primary">
                    <li><NavLink exact activeClassName="current" to="/deploy">Deploy</NavLink></li>
                </Button>

            </ul>
        </nav>
    )
}
 
export default Navigation;