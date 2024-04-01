import * as React from 'react';

import Typography from '@mui/material/Typography';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid';



export default function ThreeTabPanel(props) {
  const { value, index, client } = props;
  const [agents, setAgents] = React.useState("0");
  const [workers, setWorkers] = React.useState("0");
  const [workersBusy, setWorkersBusy] = React.useState("0");
  const [expressions, setExpressions] = React.useState([]);
  

 /*  React.useEffect(() => {
    const getChargersData = () => {
      client.get('/workers')
        .then(response => {
          setAgents(response.data.agents);
          setWorkers(response.data.workers);
          setWorkers(response.data.workers);
          setWorkersBusy(response.data.workersBusy);
          setExpressions(response.data.expressions);
        })
    }
    getChargersData()
    const interval = setInterval(() => {
      getChargersData()
    },1000);
    return () => clearInterval(interval);    
  }, []) */
    
  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
    >
      {value === index && (
        <Box sx={{ flexGrow: 1 }}>
          <Grid container spacing={2}>
            <Grid sx={{m: 3}} item md={2}>
              <Typography variant="body1" gutterBottom>
                Данные по рабочим агентам (обновляются каждые 1 с)<br />
                Агентов: {agents}<br />
                Воркеров: {workers}<br/>
                Занятых воркеров: {workersBusy}<br/>
                Обрабатываемые операции:<br />
                {expressions.map(el => (
                  <p>{el}</p>
                ))}
              </Typography>
            </Grid>
          </Grid>
        </Box>      
      )}
    </div>
  );
}
