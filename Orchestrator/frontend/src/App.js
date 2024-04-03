import React from 'react';
import './App.css';
import BasicTabs from './BasicTabs/BasicTabs.jsx';
import axios from 'axios';

const client = axios.create({
  baseURL: window.location.href,
})

function App() {

  return (
    <div className="App">
      <BasicTabs client={client}/>  
    </div>
  );
}

export default App;
