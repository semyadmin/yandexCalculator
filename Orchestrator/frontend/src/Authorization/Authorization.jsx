import * as React from 'react';
import PropTypes from 'prop-types';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import Form from './Form';

function SimpleDialog(props) {
  const { onClose, open  } = props;

  return (
    <Dialog open={open}>
          <Form onClose={onClose} />
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

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <div>      
      <Button variant="outlined" onClick={handleClickOpen} client={client}>
        Вход
      </Button>
      <SimpleDialog
        open={open}
        onClose={handleClose}
      />
    </div>
  );

}