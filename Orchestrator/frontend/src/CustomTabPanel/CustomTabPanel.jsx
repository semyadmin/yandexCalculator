import * as React from 'react';
import useWebSocket, { ReadyState } from "react-use-websocket"

import Typography from '@mui/material/Typography';
import Box from '@mui/material/Box';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Grid from '@mui/material/Grid';
import Snackbar from '@mui/material/Snackbar';
import Alert from '@mui/material/Alert';
import BasicTable from './BasicTable/BasicTable';


export default function CustomTabPanel(props) {
  const [openSuccess, setOpenSuccess] = React.useState(false);
  const [textSnackbarSuccess, setTextSnackbarSuccess] = React.useState('')
  const [textSnackbarError, setTextSnackbarError] = React.useState('')
  const [openError, setOpenError] = React.useState(false);
  const { value, index, client} = props;
  const [textValue, setTextValue] = React.useState('');
  /* const [idValue, setIdValue] = React.useState(''); */
  const [answer, setAnswer] = React.useState([]);
  const address = "ws://" + window.location.host + "/ws"
  const { sendJsonMessage, lastJsonMessage, readyState } = useWebSocket(
    address,
    {
      share: false,
      shouldReconnect: () => true,
    },
  )
  React.useEffect(() => {
    client
      .get('getexpressions')
      .then((response) => {
        if (response.data != null) {
          response.data.forEach(el => {
            const res = {
              id: el.ID,
              expression: el.Expression,
              start: el.Start,
              end: el.End,
              status: el.Status
            }
            answer.push(res)
          })
          setAnswer([...answer])
        }
      })
  },[])

  // Run when the connection state (readyState) changes
  React.useEffect(() => {
    console.log("Connection state changed")
    if (readyState === ReadyState.OPEN) {
      sendJsonMessage({
        event: "subscribe",
        data: {
          channel: "general-chatroom",
        },
      })
    }
  }, [readyState])

  // Run when a new WebSocket message is received (lastJsonMessage)
  React.useEffect(() => {
    if (lastJsonMessage != null) {
    const res = {
      id: lastJsonMessage.ID,
      expression: lastJsonMessage.Expression,
      start: lastJsonMessage.Start,
      end: lastJsonMessage.End,
      status: lastJsonMessage.Status
    }
    let change = false
    answer.forEach(el => {
      if (el.id == res.id) {
        el.status = res.status
        el.expression = res.expression
        change = true
      }
    })
    if (change == false) {
      answer.push(res)
    }
    setAnswer([...answer])
    console.log("answer ", answer)
  }
  }, [lastJsonMessage])
  const onChangeText = (event) => {
    setTextValue(event.target.value);
  }
  /* 
  const onChangeId = (event) => {
    setIdValue(event.target.value);
  } */
  
  const sendTextValue = () => {
    client
      .post('expression', textValue, {
        headers: { 
            'Content-Type' : 'text/plain' 
        }
      })
  }
  const handleCloseSuccess = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }
    setOpenSuccess(false);
  };
  const handleCloseError = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }
    setOpenError(false);
  };
 
  
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
              <Typography variant="body1" gutterBottom sx={{ maxWidth: 600 }}>
                  Введите выражение, которое хотите посчитать.<br/>
                Поддерживаемые выражения: +, -, *, /. <br />
                Поддерживаются большие числа, но точность будет низкая.<br/>
                Инкремент не поддерживается(в этом случае надо вводить 0-ваше_число) <br />
                Дробные числа поддерживаются.
                Так же можно писать выражения в скобках. Например: (1+2)*3
              </Typography>
              <TextField 
                fullWidth
                id="outlined-basic"
                label="Выражение"
                value={textValue}
                sx={{ mt: 3, mb: 2 }}
                onChange={onChangeText}
              />
              <Button
                fullWidth
                variant="contained"
                onClick={() => sendTextValue()}
                sx={{ mt: 3, mb: 2 }}
              >Расчитать</Button>
            </Grid>
            <Grid sx={{m: 3}}>
              <BasicTable rows={answer} />
            </Grid>
          </Grid>
        </Box>      
      )}
      <Snackbar open={openSuccess} autoHideDuration={6000} onClose={handleCloseSuccess}>
        <Alert
          onClose={handleCloseSuccess}
          severity="success"
          variant="filled"
          sx={{ width: '100%' }}
        >
          {textSnackbarSuccess}
        </Alert>
      </Snackbar>
      <Snackbar open={openError} autoHideDuration={6000} onClose={handleCloseError}>
        <Alert
          onClose={handleCloseError}
          severity="error"
          variant="filled"
          sx={{ width: '100%' }}
        >
          {textSnackbarError}
        </Alert>
      </Snackbar>
    </div>
  );
}
