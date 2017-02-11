import React, { Component } from 'react';
import SearchForm from './searchform.js';
import logo from './logo.svg';
import './App.css';

class App extends Component {
  render() {
    return (
      <div className="App">
        <div className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h2>Flight Search</h2>
        </div>
        <SearchForm />
      </div>
    );
  }
}

export default App
