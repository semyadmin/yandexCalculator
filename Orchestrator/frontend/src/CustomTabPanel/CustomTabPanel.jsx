import * as React from 'react';

import Typography from '@mui/material/Typography';
import Box from '@mui/material/Box';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Grid from '@mui/material/Grid';


export default function CustomTabPanel(props) {
  const { children, value, index, client } = props;
  const [textValue, setTextValue] = React.useState('');
  const onChangeText = (event) => {
    setTextValue(event.target.value);
  }
  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
    >
      {value === index && (
          <Box sx={{ flexGrow: 1 }}>
            <Grid container spacing={2}>
              <Grid sx={{m: 3}}>
                <Typography>{children}</Typography>
                <TextField 
                  id="outlined-basic"
                  label="Outlined"
                  variant="outlined"
                  value={textValue}
                  onChange={onChangeText}
                />
                <Button
                  fullWidth
                  variant="contained"
                  onClick={() => onChangeText()}
                  sx={{ mt: 3, mb: 2 }}
                >Contained</Button>
              </Grid>
              <Grid sx={{m: 3}}>
                <Typography>{textValue}</Typography>
              </Grid>
            </Grid>
          </Box>
          
      )}
    </div>
  );
}
