import * as React from 'react';
import PropTypes from 'prop-types';
import Typography from '@mui/material/Typography';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import Container from '@mui/material/Container';
import Toolbar from '@mui/material/Toolbar';
import Form from './Form';

function SimpleDialog(props) {
  const { onClose, open, client, setUser  } = props;

  return (
    <Dialog open={open}>
          <Form onClose={onClose} client={client} setUser={setUser} />
    </Dialog>
  );
}

SimpleDialog.propTypes = {
  onClose: PropTypes.func.isRequired,
  open: PropTypes.bool.isRequired,
};

export default function Authorization(props) {
  const {client} = props
  const [open, setOpen] = React.useState(false);
  const [user, setUser] = React.useState({});
  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClickExit = () => {
    setUser("");

  };

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <Container maxWidth="xl">
      <Toolbar>
      {
        user !== ""
        ? <Typography variant="h6" sx={{ margin: 1 }}>{user}</Typography>
        : null
      }
      { 
        user !== ""         
        ? (<Button variant="outlined" onClick={handleClickExit}>Выход</Button>)
        : (
          <Button variant="outlined" onClick={handleClickOpen}>
            Вход
          </Button>
          )
      
      }
        <SimpleDialog
          open={open}
          onClose={handleClose}
          client={client}
          setUser={setUser}
        />
      </Toolbar>   
    </Container>
  );

}