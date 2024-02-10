import * as React from 'react';

import Typography from '@mui/material/Typography';
import Box from '@mui/material/Box';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Grid from '@mui/material/Grid';


export default function CustomTabPanel(props) {
  const { onChange, children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      {...other}
    >
      {value === index && (
            <Grid item md={12} container spacing={2}>
              <Grid xs={4} sx={{m: 3}}>
                <Box sx={{ p: 3}} align="center" display="flex" flexDirection="column" spacing={2}>
                  <Typography>{children}</Typography>
                  <TextField 
                    id="outlined-basic"
                    label="Outlined"
                    onChange={onChange}
                  />
                  <Button variant="contained">Contained</Button>
                </Box>
              </Grid>
              <Grid item xs={4}>
            </Grid>
 
            </Grid>
          
      )}
    </div>
  );
}
