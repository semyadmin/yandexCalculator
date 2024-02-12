import * as React from 'react';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import AutorenewIcon from '@mui/icons-material/Autorenew';



export default function BasicTable(rows) {
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
          {rows.rows.map((row) => (
            
            <TableRow
              key={row.id}
              sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
            >
              <TableCell><AutorenewIcon /></TableCell>
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