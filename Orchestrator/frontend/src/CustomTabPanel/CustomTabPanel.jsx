import * as React from 'react';

import Typography from '@mui/material/Typography';
import Box from '@mui/material/Box';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Grid from '@mui/material/Grid';
import BasicTable from './BasicTable/BasicTable';
import Snackbar from '@mui/material/Snackbar';
import Alert from '@mui/material/Alert';


export default function CustomTabPanel(props) {
  const [openSuccess, setOpenSuccess] = React.useState(false);
  const [textSnackbarSuccess, setTextSnackbarSuccess] = React.useState('')
  const [textSnackbarError, setTextSnackbarError] = React.useState('')
  const [openError, setOpenError] = React.useState(false);
  const { value, index, client } = props;
  const [textValue, setTextValue] = React.useState('');
  const [idValue, setIdValue] = React.useState('');
  const [answer, setAnswer] = React.useState([]);
  const onChangeText = (event) => {
    setTextValue(event.target.value);
  }
  const onChangeId = (event) => {
    setIdValue(event.target.value);
  }
  const sendTextValue = () => {
    client
      .post('expression', textValue, {
        headers: { 
            'Content-Type' : 'text/plain' 
        }
      })
      .then((response) => {
        const res = {
          id: response.data.ID,
          expression: response.data.Expression,
          start: response.data.Start,
          end: response.data.End,
          status: response.data.Status
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
          setOpenSuccess(false)
          setTextSnackbarSuccess(`Выражение ${res.expression} успешно поставлено на обработку`)
          setOpenSuccess(true)
          setAnswer([...answer,res])
        } else {
          setOpenSuccess(false)
          setTextSnackbarSuccess(`Выражение ${res.expression} уже обрабатывается`)
          setOpenSuccess(true)
          setAnswer(answer)
        }
      })
      .catch(error => {
        console.log(error)
        setTextSnackbarError(`Введенное выражение некорректно!`)
        setOpenError(true)
      })
  }
  const sendIdValue = () => {
    client
      .get('/id/'+idValue)
      .then((response) => {
        const res = {
          id: response.data.ID,
          expression: response.data.Expression,
          start: response.data.Start,
          end: response.data.End,
          status: response.data.Status
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
          setOpenSuccess(false)
          setTextSnackbarSuccess(`Выражение ${res.expression} успешно поставлено на обработку`)
          setOpenSuccess(true)
          setAnswer([...answer,res])
        } else {
          setOpenSuccess(false)
          setTextSnackbarSuccess(`Выражение ${res.expression} уже обрабатывается`)
          setOpenSuccess(true)
          setAnswer(answer)
        }
      })
      .catch(error => {
        console.log(error)
        setTextSnackbarError(`Данного выражения не существует`)
        setOpenError(true)
      })
  }
  const setIdValueFromTable = (id) => {
    client
    .get('/id/'+id)
    .then((response) => {
      const res = {
        id: response.data.ID,
        expression: response.data.Expression,
        start: response.data.Start,
        end: response.data.End,
        status: response.data.Status
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
        setOpenSuccess(false)
        setTextSnackbarSuccess(`Выражение ${res.expression} успешно поставлено на обработку`)
        setOpenSuccess(true)
        setAnswer([...answer,res])
      } else {
        setOpenSuccess(false)
        setTextSnackbarSuccess(`Выражение ${res.expression} уже обрабатывается`)
        setOpenSuccess(true)
        setAnswer(answer)
      }
    })
    .catch(error => {
      console.log(error)
      setTextSnackbarError(`Данного выражения не существует`)
      setOpenError(true)
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
              <TextField 
                fullWidth
                id="find-id"
                label="Найти по ID"
                value={idValue}
                sx={{ mt: 3, mb: 2 }}
                onChange={onChangeId}
              />
              <Button
                fullWidth
                variant="contained"
                onClick={() => sendIdValue()}
                sx={{ mt: 3, mb: 2 }}
              >Найти</Button>
            </Grid>
            <Grid sx={{m: 3}}>
              <BasicTable rows={answer} sendIdValue={setIdValueFromTable} />
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
