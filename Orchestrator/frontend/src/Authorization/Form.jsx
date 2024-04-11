import * as React from 'react';
import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import TextField from '@mui/material/TextField';
import Box from '@mui/material/Box';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import { createTheme, ThemeProvider } from '@mui/material/styles';


const defaultTheme = createTheme();

export default function Form(props) {
    const { onClose, client, setUser } = props;
    const [errorLogin, setErrorLogin] = React.useState(true);
    const [login, setLogin] = React.useState("");
    const [password, setPassword] = React.useState("");
    const handleChangeLogin = (event) => {
        setLogin(event.target.value);
    }

    const handleChangePassword = (event) => {
        setPassword(event.target.value);
    }
    const handleSubmit = (event) => {
        event.preventDefault();
        setUser(login);
        onClose();
        return
        client.post('/api/v1/register', { login: login, password: password }, {
            headers: {
                'Content-Type': 'application/json'
            }
        }).then((response) => {
            if (response.data != null) {
                localStorage.setItem('token', response.data.token);
                setUser(login);
                setLogin("");
                setPassword("");
                onClose();
              }
        })
        .catch((error) => {
            setErrorLogin(true);
        });
    };
    const handleOnFocus = () => {
        setErrorLogin(false)
    }

  return (
    <ThemeProvider theme={defaultTheme}>
      <Container component="main" maxWidth="xs">
        <CssBaseline />
        <Box
          sx={{
            marginTop: 8,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
          }}
        >
          <Avatar sx={{ m: 1, bgcolor: 'secondary.main' }}>
            <LockOutlinedIcon />
          </Avatar>
          <Typography component="h1" variant="h5">
            Вход
          </Typography>
          <Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
            <TextField
                {...(errorLogin ? { error: true, helperText: 'Неверное имя пользователя или пароль' } : {})}
                onFocus={handleOnFocus}
              margin="normal"
              required
              fullWidth
              id="login"
              label="Имя пользователя"
              name="login"
              onChange={handleChangeLogin}
            />
            <TextField
                {...(errorLogin ? { error: true } : {})}
              margin="normal"
              required
              fullWidth
              name="password"
              label="Password"
              type="password"
              id="password"
              autoComplete="current-password"
              onChange={handleChangePassword}
            />
            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{ mt: 3, mb: 2 }}
            >
              Войти
            </Button>
          </Box>
        </Box>
      </Container>
    </ThemeProvider>
  );
}