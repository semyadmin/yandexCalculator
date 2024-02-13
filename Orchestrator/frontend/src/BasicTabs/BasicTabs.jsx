import * as React from 'react';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import Box from '@mui/material/Box';
import CustomTabPanel from '../CustomTabPanel/CustomTabPanel';
  
  export default function BasicTabs({client}) {
    const [value, setValue] = React.useState(0);
  
    const handleChange = (event, newValue) => {
      setValue(newValue);
    };
  
    return (
      <Box sx={{ width: '100%' }}>
        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs value={value} onChange={handleChange} aria-label="basic tabs example">
            <Tab label="Рассчитать выражение" />
            <Tab label="Item Two" />
            <Tab label="Item Three"  />
          </Tabs>
        </Box>
        <CustomTabPanel value={value} index={0} client={client} />
      </Box>
    );
  }