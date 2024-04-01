import * as React from 'react';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import AutorenewIcon from '@mui/icons-material/Autorenew';
import CheckIcon from '@mui/icons-material/Check';
import IconButton from '@mui/material/IconButton';
import Tooltip from '@mui/material/Tooltip';
import CloseIcon from '@mui/icons-material/Close';



export default function BasicTable(props) {
  const onClickIdValue = (id) => {
    props.sendIdValue(id)
  }
  const address = "ws://" + window.location.host + "/ws"
  console.log(address)
  const socket = new WebSocket(address)

  // Connection opened
  socket.addEventListener("open", event => {
    socket.send("Connection established")
  });

  // Listen for messages
  socket.addEventListener("message", event => {
    console.log("Message from server ", event.data)
  });
  return (
    <TableContainer component={Paper}>
      <Table sx={{ minWidth: 650 }} aria-label="simple table">
        <TableHead>
          <TableRow>
            <TableCell>Статус</TableCell>
            <TableCell>Идентификатор</TableCell>
            <TableCell align="right">Выражение</TableCell>
            <TableCell align="right">Дата начала обработки</TableCell>
            <TableCell align="right">Возможная дата окончания обработки</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {props.rows.map((row) => (
            
            <TableRow
              key={row.id}
              sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
            >
              <TableCell >
                
                  {row.status == 'progress' 
                  ? <Tooltip title="Обновить информацию"><IconButton><AutorenewIcon
                      /* onClick={() => {
                        onClickIdValue(row.id)
                    }} */
                      sx={{ "&:hover": { color: "green" } }}
                  /></IconButton></Tooltip>
                  : row.status == 'completed' ? <Tooltip title="Выражение посчитано"><CheckIcon/></Tooltip> 
                  : <Tooltip title="Выражение посчитано с ошибкой"><CloseIcon/></Tooltip> 
                  }
                
               </TableCell>
              <TableCell component="th" scope="row">
                {row.id}
              </TableCell>
              <TableCell align="right">{row.expression}</TableCell>
              <TableCell align="right">{row.start}</TableCell>
              <TableCell align="right">{row.end}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}