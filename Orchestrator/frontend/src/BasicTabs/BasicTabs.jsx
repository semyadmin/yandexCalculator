import * as React from 'react';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import Box from '@mui/material/Box';
import CustomTabPanel from '../CustomTabPanel/CustomTabPanel';
import TwoTabPanel from '../TwoTabPanel/TwoTabPanel';
  
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
            <Tab label="Установить продолжительность" />
            <Tab label="Item Three"  />
          </Tabs>
        </Box>
        <CustomTabPanel value={value} index={0} client={client} />
        <TwoTabPanel value={value} index={1} client={client} />
      </Box>
    );
  }